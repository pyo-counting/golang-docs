// filepath: helm-ecr-api/internal/service/ecr_service.go
package service

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

var (
	// ErrChartNotFound는 차트를 찾을 수 없을 때 반환되는 에러입니다.
	ErrChartNotFound = errors.New("chart not found")
	// ErrRepositoryNotAllowed는 허용되지 않은 리포지토리에 접근 시 반환되는 에러입니다.
	ErrRepositoryNotAllowed = errors.New("repository not allowed")
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
	sts    *sts.Client
	awsCfg aws.Config

	allowedRepos map[string]struct{} // 빠른 조회를 위해 map 사용

	accountIDOnce sync.Once
	accountID     string
	accountIDErr  error
}

// NewECRService는 ECRService의 새 인스턴스를 생성합니다.
func NewECRService(cfg aws.Config, allowedRepos []string) *ECRService {
	allowedReposMap := make(map[string]struct{}, len(allowedRepos))
	for _, repo := range allowedRepos {
		allowedReposMap[repo] = struct{}{}
	}

	return &ECRService{
		client:       ecr.NewFromConfig(cfg),
		sts:          sts.NewFromConfig(cfg),
		awsCfg:       cfg,
		allowedRepos: allowedReposMap,
	}
}

// isRepoAllowed는 요청된 리포지토리가 허용 목록에 있는지 확인합니다.
func (s *ECRService) isRepoAllowed(repoName string) bool {
	_, ok := s.allowedRepos[repoName]
	return ok
}

// getAccountID는 sync.Once를 사용하여 AWS 계정 ID를 한 번만 조회하고 캐싱합니다.
func (s *ECRService) getAccountID(ctx context.Context) (string, error) {
	s.accountIDOnce.Do(func() {
		identity, err := s.sts.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			s.accountIDErr = fmt.Errorf("failed to get caller identity: %w", err)
			return
		}
		s.accountID = aws.ToString(identity.Account)
	})
	return s.accountID, s.accountIDErr
}

// DescribeHelmChart는 ECR에서 특정 Helm 차트(OCI 이미지)의 상세 정보를 조회합니다.
func (s *ECRService) DescribeHelmChart(ctx context.Context, repoName, tag, digest string) ([]types.ImageDetail, error) {
	if !s.isRepoAllowed(repoName) {
		return nil, fmt.Errorf("%w: %s", ErrRepositoryNotAllowed, repoName)
	}

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
	repoNames := make([]string, 0, len(s.allowedRepos))
	for repo := range s.allowedRepos {
		repoNames = append(repoNames, repo)
	}

	if len(repoNames) == 0 {
		return []types.Repository{}, nil // 허용된 리포지토리가 없으면 빈 목록 반환
	}

	result, err := s.client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: repoNames,
	})
	if err != nil {
		return nil, err
	}

	return result.Repositories, nil
}

// GetChartFile은 ECR에서 차트(.tar.gz)를 다운로드하고 압축을 해제하여
// 특정 파일(예: 'values.yaml', 'Chart.yaml')의 내용을 반환합니다.
// 이 함수는 go-containerregistry 라이브러리를 사용하여 OCI 표준 방식으로 차트를 가져옵니다.
func (s *ECRService) GetChartFile(ctx context.Context, repoName, tag, digest, fileName string) ([]byte, error) {
	if !s.isRepoAllowed(repoName) {
		return nil, fmt.Errorf("%w: %s", ErrRepositoryNotAllowed, repoName)
	}

	// 1. 캐시된 AWS 계정 ID를 가져옵니다.
	accountID, err := s.getAccountID(ctx)
	if err != nil {
		return nil, err
	}

	// 2. ECR 리포지토리 URI를 구성합니다.
	// 예: 123456789012.dkr.ecr.ap-northeast-2.amazonaws.com/my-helm-charts/my-app
	repoURI := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s", accountID, s.awsCfg.Region, repoName)

	var ref name.Reference
	if tag != "" {
		ref, err = name.NewTag(fmt.Sprintf("%s:%s", repoURI, tag))
	} else if digest != "" {
		ref, err = name.NewDigest(fmt.Sprintf("%s@%s", repoURI, digest))
	} else {
		return nil, fmt.Errorf("either tag or digest must be provided")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse image reference: %w", err)
	}

	// 3. ECR 인증 토큰을 가져와서 직접 인증을 설정합니다.
	token, err := s.getECRAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ECR auth token: %w", err)
	}

	// 4. go-containerregistry를 사용하여 OCI 이미지를 가져옵니다.
	img, err := remote.Image(ref, remote.WithContext(ctx), remote.WithAuth(token))
	if err != nil {
		// 404 Not Found와 같은 특정 오류를 확인하여 커스텀 에러를 반환합니다.
		var transportErr *transport.Error
		if errors.As(err, &transportErr) && transportErr.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s", ErrChartNotFound, ref.Name())
		}
		return nil, fmt.Errorf("failed to get remote image: %w", err)
	}

	// 5. 이미지 레이어를 순회하여 Helm 차트 콘텐츠 레이어를 찾습니다.
	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get image layers: %w", err)
	}

	for _, layer := range layers {
		mediaType, err := layer.MediaType()
		if err != nil {
			continue
		}

		// Helm 차트 콘텐츠의 mediaType은 'application/vnd.cncf.helm.chart.content.v1.tar+gzip' 입니다.
		if string(mediaType) == "application/vnd.cncf.helm.chart.content.v1.tar+gzip" {
			// 6. 레이어의 압축을 풀고 tar 아카이브에서 원하는 파일을 찾습니다.
			// layer.Uncompressed()는 라이브러리가 gzip 압축을 자동으로 처리하도록 합니다.
			rc, err := layer.Uncompressed()
			if err != nil {
				return nil, fmt.Errorf("failed to uncompress layer: %w", err)
			}
			defer rc.Close()

			tarReader := tar.NewReader(rc)
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
		}
	}

	return nil, fmt.Errorf("'%s' not found in chart archive", fileName)
}

// getECRAuthToken은 AWS ECR로부터 인증 토큰을 가져와서 Basic 인증 형태로 반환합니다.
func (s *ECRService) getECRAuthToken(ctx context.Context) (authn.Authenticator, error) {
	// ECR GetAuthorizationToken API를 호출하여 인증 토큰을 가져옵니다.
	output, err := s.client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ECR authorization token: %w", err)
	}

	if len(output.AuthorizationData) == 0 || output.AuthorizationData[0].AuthorizationToken == nil {
		return nil, fmt.Errorf("no authorization data received from ECR")
	}

	// Base64로 인코딩된 토큰을 디코딩합니다.
	token := *output.AuthorizationData[0].AuthorizationToken
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ECR authorization token: %w", err)
	}

	// 토큰은 "username:password" 형태입니다.
	parts := strings.SplitN(string(decodedToken), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ECR authorization token format")
	}

	// Basic 인증 형태로 반환합니다.
	return &authn.Basic{
		Username: parts[0],
		Password: parts[1],
	}, nil
}