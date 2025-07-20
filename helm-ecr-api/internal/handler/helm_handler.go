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
// - tag 쿼리 파라미터가 있으면: 특정 태그의 상세 정보를 반환
// - tag 쿼리 파라미터가 없으면: 모든 태그 목록을 반환
// - view=values-schema 파라미터가 있으면: values.schema.json 파일 내용을 반환
// 예: GET /v1/helm-charts/my-repo?tag=1.2.3&view=values-schema
func (h *HelmHandler) GetHelmChart(w http.ResponseWriter, r *http.Request) {
	// URL 경로에서 리포지토리 이름을 추출합니다.
	repoName := strings.TrimPrefix(r.URL.Path, "/v1/helm-charts/")
	tag := r.URL.Query().Get("tag")
	view := r.URL.Query().Get("view")

	if repoName == "" {
		h.logger.Warn("missing required query parameters", "repo", repoName)
		http.Error(w, "missing repository name in URL path", http.StatusBadRequest)
		return
	}

	// view 파라미터에 따라 분기 처리
	if view == "values-schema" {
		if tag == "" {
			http.Error(w, "tag is required for values-schema view", http.StatusBadRequest)
			return
		}
		h.getValuesSchema(w, r, repoName, tag)
		return
	}

	// 저장소에 있는 이미지 정보 조회 (tag 유무에 따라 서비스에서 다르게 처리)
	h.logger.Info("request to get helm chart info", "repo", repoName, "tag", tag)

	chart, err := h.chartService.DescribeHelmChart(r.Context(), repoName, tag)
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

func (h *HelmHandler) getValuesSchema(w http.ResponseWriter, r *http.Request, repoName, tag string) {
	h.logger.Info("request to get values.schema.json", "repo", repoName, "tag", tag)

	schemaBytes, err := h.chartService.GetValuesSchema(r.Context(), repoName, tag)
	if err != nil {
		h.logger.Error("failed to get values.schema.json", "error", err)
		// TODO: 서비스 계층에서 반환된 에러 유형에 따라 더 구체적인 상태 코드를 반환하도록 개선할 수 있습니다.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(schemaBytes)
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
