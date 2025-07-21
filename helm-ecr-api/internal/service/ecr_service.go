// filepath: helm-ecr-api/internal/service/ecr_service.go
package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// ChartService는 Helm 차트 관련 비즈니스 로직에 대한 인터페이스입니다.
// 이를 통해 핸들러는 실제 구현으로부터 분리되어 테스트 용이성이 높아집니다.
type ChartService interface {
	DescribeHelmChart(ctx context.Context, repoName, tag, digest string) ([]types.ImageDetail, error)
	ListHelmCharts(ctx context.Context) ([]types.Repository, error)
	GetChartFile(ctx context.Context, repoName, tag, digest, fileName string) ([]byte, error)
}

// ECRService는 ChartService 인터페이스의 구현체입니다.
// ECRService는 ECR과 상호작용하는 비즈니스 로직을 담당합니다.
type ECRService struct {
	client *ecr.Client
}

// NewECRService는 ECRService의 새 인스턴스를 생성합니다.
func NewECRService(cfg aws.Config) *ECRService {
	return &ECRService{client: ecr.NewFromConfig(cfg)}
}

// DescribeHelmChart는 ECR에서 특정 Helm 차트(OCI 이미지)의 상세 정보를 조회합니다.
func (s *ECRService) DescribeHelmChart(ctx context.Context, repoName, tag, digest string) ([]types.ImageDetail, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	}

	// tag 또는 digest 파라미터가 있을 때만 특정 이미지로 필터링
	if tag != "" {
		input.ImageIds = []types.ImageIdentifier{{ImageTag: aws.String(tag)}}
	} else if digest != "" {
		input.ImageIds = []types.ImageIdentifier{{ImageDigest: aws.String(digest)}}
	}

	result, err := s.client.DescribeImages(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.ImageDetails) == 0 {
		if tag != "" {
			return nil, &types.ImageNotFoundException{Message: aws.String(fmt.Sprintf("chart not found with tag: %s", tag))}
		}
		if digest != "" {
			return nil, &types.ImageNotFoundException{Message: aws.String(fmt.Sprintf("chart not found with digest: %s", digest))}
		}
		return nil, &types.RepositoryNotFoundException{Message: aws.String(fmt.Sprintf("chart not found in repository: %s", repoName))}
	}

	return result.ImageDetails, nil
}

// ListHelmCharts는 ECR에 있는 모든 리포지토리를 조회합니다.
// ECR API는 페이지네이션을 사용하므로, 모든 결과를 가져오기 위해 반복 호출합니다.
func (s *ECRService) ListHelmCharts(ctx context.Context) ([]types.Repository, error) {

	result, err := s.client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return nil, err
	}

	if len(result.Repositories) == 0 {
		return nil, &types.RepositoryNotFoundException{Message: aws.String("chart not found in any repository")}
	}
	return result.Repositories, nil
}

// GetChartFile은 ECR에서 차트(.tar.gz)를 다운로드하고 압축을 해제하여
// 특정 파일(예: 'values.yaml', 'Chart.yaml')의 내용을 반환합니다.
func (s *ECRService) GetChartFile(ctx context.Context, repoName, tag, digest, fileName string) ([]byte, error) {
	// 1. 이미지 매니페스트를 가져와서 차트 레이어의 다이제스트(digest)를 찾습니다.
	// Helm 차트의 .tar.gz 파일은 OCI 이미지의 레이어로 저장됩니다.
	batchGetImageInput := &ecr.BatchGetImageInput{
		RepositoryName: aws.String(repoName),
		AcceptedMediaTypes: []string{
			"application/vnd.oci.image.manifest.v1+json",
		},
	}
	if tag != "" {
		batchGetImageInput.ImageId = &types.ImageIdentifier{ImageTag: aws.String(tag)}
	} else if digest != "" {
		batchGetImageInput.ImageId = &types.ImageIdentifier{ImageDigest: aws.String(digest)}
	} else {
		return nil, fmt.Errorf("either tag or digest must be provided")
	}

	batchGetImageOutput, err := s.client.BatchGetImage(ctx, batchGetImageInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get image manifest: %w", err)
	}
	if len(batchGetImageOutput.Images) == 0 || batchGetImageOutput.Images[0].ImageManifest == nil {
		return nil, &types.ImageNotFoundException{Message: aws.String("image manifest not found")}
	}

	var manifest struct {
		Layers []struct {
			MediaType string `json:"mediaType"`
			Digest    string `json:"digest"`
		} `json:"layers"`
	}
	if err := json.Unmarshal([]byte(*batchGetImageOutput.Images[0].ImageManifest), &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal image manifest: %w", err)
	}

	var chartLayerDigest string
	for _, layer := range manifest.Layers {
		// Helm 차트 콘텐츠의 mediaType은 'application/vnd.cncf.helm.chart.content.v1.tar+gzip' 입니다.
		if layer.MediaType == "application/vnd.cncf.helm.chart.content.v1.tar+gzip" {
			chartLayerDigest = layer.Digest
			break
		}
	}
	if chartLayerDigest == "" {
		return nil, fmt.Errorf("could not find helm chart content layer in manifest")
	}

	// 2. 찾은 다이제스트를 이용해 레이어(.tar.gz 파일)를 다운로드합니다.
	// GetLayers API는 Deprecated 되었으므로 GetDownloadUrlForLayer를 사용합니다.
	getDownloadURLOutput, err := s.client.GetDownloadUrlForLayer(ctx, &ecr.GetDownloadUrlForLayerInput{
		RepositoryName: aws.String(repoName),
		LayerDigest:    aws.String(chartLayerDigest),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get download url for layer: %w", err)
	}

	// 3. 다운로드 URL을 이용해 실제 레이어 데이터를 가져옵니다.
	resp, err := http.Get(*getDownloadURLOutput.DownloadUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download layer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download layer: status code %d", resp.StatusCode)
	}

	layerBlob, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read layer blob: %w", err)
	}

	// 4. 다운로드한 blob의 압축을 해제하고 요청된 파일을 찾습니다.
	gzipReader, err := gzip.NewReader(bytes.NewReader(layerBlob))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 파일 끝
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar archive: %w", err)
		}

		// 파일 경로는 'chart-name/values.yaml' 형태일 수 있으므로 HasSuffix로 확인합니다.
		if strings.HasSuffix(header.Name, fileName) {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("'%s' not found in chart archive", fileName)
}
		if strings.HasSuffix(header.Name, fileName) {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("'%s' not found in chart archive", fileName)
}
