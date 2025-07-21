// filepath: helm-ecr-api/internal/handler/helm_handler.go
package handler

import (
	"encoding/json"
	"errors"
	"helm-ecr-api/internal/service"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// HelmHandler는 HTTP 요청을 처리하고 서비스 계층을 호출합니다.
type HelmHandler struct {
	chartService service.ChartService
	logger       *slog.Logger
}

// NewHelmHandler는 HelmHandler의 새 인스턴스를 생성합니다.
func NewHelmHandler(chartService service.ChartService, logger *slog.Logger) *HelmHandler {
	return &HelmHandler{
		chartService: chartService,
		logger:       logger,
	}
}

// GetHelmChart는 특정 차트의 정보를 조회하는 HTTP 핸들러입니다.
// - tag 또는 digest 쿼리 파라미터가 있으면: 특정 버전의 상세 정보를 반환
// - 쿼리 파라미터가 없으면: 모든 태그 목록을 반환
func (h *HelmHandler) GetHelmChart(w http.ResponseWriter, r *http.Request) {
	// URL 경로에서 리포지토리 이름을 추출합니다.
	repoName := r.PathValue("chart-name")
	tag := r.URL.Query().Get("tag")
	digest := r.URL.Query().Get("digest")

	if repoName == "" {
		h.logger.Warn("missing repository name in URL path")
		http.Error(w, "missing repository name in URL path", http.StatusBadRequest)
		return
	}

	if tag != "" && digest != "" {
		http.Error(w, "tag and digest cannot be specified simultaneously", http.StatusBadRequest)
		return
	}

	// 저장소에 있는 이미지 정보 조회 (tag 또는 digest 유무에 따라 서비스에서 다르게 처리)
	h.logger.Info("request to get helm chart info", "repo", repoName, "tag", tag, "digest", digest)

	chart, err := h.chartService.DescribeHelmChart(r.Context(), repoName, tag, digest)
	if err != nil {
		h.logger.Error("failed to describe helm chart", "error", err)

		// 에러 유형에 따라 적절한 HTTP 상태 코드 반환
		var notFoundErr *types.ImageNotFoundException
		var repoNotFoundErr *types.RepositoryNotFoundException

		if errors.As(err, &notFoundErr) || errors.As(err, &repoNotFoundErr) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chart); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// GetChartFile은 차트 아카이브 내의 특정 파일을 조회하는 핸들러입니다.
// 예: GET /v1/files/my-repo/my-app/values.yaml?tag=1.2.3
func (h *HelmHandler) GetChartFile(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("chart-name")
	fileName := r.PathValue("file-name")
	tag := r.URL.Query().Get("tag")
	digest := r.URL.Query().Get("digest")

	if repoName == "" || fileName == "" {
		http.Error(w, "repository name and file name are required in URL path", http.StatusBadRequest)
		return
	}

	if tag == "" && digest == "" {
		http.Error(w, "tag or digest is required", http.StatusBadRequest)
		return
	}

	if tag != "" && digest != "" {
		http.Error(w, "tag and digest cannot be specified simultaneously", http.StatusBadRequest)
		return
	}

	h.logger.Info("request to get chart file", "repo", repoName, "tag", tag, "digest", digest, "file", fileName)

	// 파일 확장자에 따라 적절한 Content-Type을 설정합니다.
	// YAML의 공식 IANA MIME 타입은 없지만, 'application/x-yaml'이 널리 사용되는 관례입니다.
	var contentType string
	if strings.HasSuffix(fileName, ".json") {
		contentType = "application/json"
	} else if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		contentType = "application/x-yaml; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	fileBytes, err := h.chartService.GetChartFile(r.Context(), repoName, tag, digest, fileName)
	if err != nil {
		h.logger.Error("failed to get chart file", "error", err, "file", fileName)
		// TODO: 서비스 계층에서 반환된 에러 유형에 따라 더 구체적인 상태 코드를 반환하도록 개선할 수 있습니다.
		// 예를 들어, 파일이 없는 경우 404 Not Found를 반환할 수 있습니다.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(fileBytes)
}

// ListHelmCharts는 ECR의 모든 Helm 차트 리포지토리를 조회하는 핸들러입니다.
// 예: GET /v1/helm-charts
func (h *HelmHandler) ListHelmCharts(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("request to list all helm chart repositories")

	charts, err := h.chartService.ListHelmCharts(r.Context())
	if err != nil {
		h.logger.Error("failed to list helm charts", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(charts); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// HealthCheck는 서비스의 상태를 확인하는 간단한 핸들러입니다.
// 200 OK 응답을 반환하여 서비스가 살아있음을 알립니다.
func (h *HelmHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// 간단한 JSON 응답을 보냅니다.
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		h.logger.Error("failed to write health check response", "error", err)
	}
}
