// filepath: helm-ecr-api/internal/handler/helm_handler.go
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"helm-ecr-api/internal/service"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// contextKey는 컨텍스트 값의 키로 사용되어 충돌을 방지합니다.
type contextKey string

const (
	chartNameKey contextKey = "chart-name"
	fileNameKey  contextKey = "file-name"
)

// HelmHandler는 HTTP 요청을 처리하고 서비스 계층을 호출합니다.
type HelmHandler struct {
	chartService service.ChartService
	logger       *slog.Logger
	filePattern  *regexp.Regexp
}

// NewHelmHandler는 HelmHandler의 새 인스턴스를 생성합니다.
func NewHelmHandler(chartService service.ChartService, logger *slog.Logger) *HelmHandler {
	return &HelmHandler{
		chartService: chartService,
		logger:       logger,
		// 정규식을 한 번만 컴파일하여 재사용합니다.
		filePattern: regexp.MustCompile(`^(.+)/files/([^/]+)$`),
	}
}

// GetHelmChart는 특정 차트의 정보를 조회하는 HTTP 핸들러입니다.
// - tag 또는 digest 쿼리 파라미터가 있으면: 특정 버전의 상세 정보를 반환
// - 쿼리 파라미터가 없으면: 모든 태그 목록을 반환
func (h *HelmHandler) GetHelmChart(w http.ResponseWriter, r *http.Request) {
	// URL 경로에서 리포지토리 이름을 추출합니다.
	repoName, ok := r.Context().Value(chartNameKey).(string)
	if !ok || repoName == "" {
		h.respondError(w, http.StatusBadRequest, "missing repository name in URL path")
		return
	}
	tag := r.URL.Query().Get("tag")
	digest := r.URL.Query().Get("digest")

	if tag != "" && digest != "" {
		h.respondError(w, http.StatusBadRequest, "tag and digest cannot be specified simultaneously")
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

		if errors.Is(err, service.ErrRepositoryNotAllowed) {
			h.respondError(w, http.StatusForbidden, err.Error())
		} else if errors.As(err, &notFoundErr) || errors.As(err, &repoNotFoundErr) {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.respondJSON(w, http.StatusOK, chart)
}

// GetChartFile은 차트 아카이브 내의 특정 파일을 조회하는 핸들러입니다.
// 예: GET /v1/helm-charts/my-repo/my-app/files/values.yaml?tag=1.2.3
func (h *HelmHandler) GetChartFile(w http.ResponseWriter, r *http.Request) {
	repoName, ok := r.Context().Value(chartNameKey).(string)
	if !ok || repoName == "" {
		h.respondError(w, http.StatusBadRequest, "missing repository name in URL path")
		return
	}
	fileName, ok := r.Context().Value(fileNameKey).(string)
	if !ok || fileName == "" {
		h.respondError(w, http.StatusBadRequest, "missing file name in URL path")
		return
	}
	tag := r.URL.Query().Get("tag")
	digest := r.URL.Query().Get("digest")

	if tag == "" && digest == "" {
		h.respondError(w, http.StatusBadRequest, "tag or digest is required")
		return
	}

	if tag != "" && digest != "" {
		h.respondError(w, http.StatusBadRequest, "tag and digest cannot be specified simultaneously")
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

		// 서비스 계층에서 정의한 커스텀 에러를 확인하여 적절한 상태 코드를 반환합니다.
		if errors.Is(err, service.ErrChartNotFound) {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrRepositoryNotAllowed) {
			h.respondError(w, http.StatusForbidden, err.Error())
		} else {
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(fileBytes)
}

// RouteChartDetails는 /v1/helm-charts/ 경로에 대한 요청을 분석하여
// 적절한 핸들러로 분기하는 서브-라우터 역할을 합니다.
func (h *HelmHandler) RouteChartDetails(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/helm-charts/")

	// 미리 컴파일된 정규식을 사용하여 .../files/... 패턴을 찾습니다.
	matches := h.filePattern.FindStringSubmatch(path)

	if len(matches) == 3 {
		// 파일 조회 요청: /v1/helm-charts/{chart-name}/files/{file-name}
		chartName := matches[1]
		fileName := matches[2]

		// 컨텍스트에 파라미터를 추가하여 핸들러에 전달합니다.
		ctx := context.WithValue(r.Context(), chartNameKey, chartName)
		ctx = context.WithValue(ctx, fileNameKey, fileName)
		h.GetChartFile(w, r.WithContext(ctx))
	} else {
		// 차트 정보 조회 요청: /v1/helm-charts/{chart-name}
		chartName := path
		ctx := context.WithValue(r.Context(), chartNameKey, chartName)
		h.GetHelmChart(w, r.WithContext(ctx))
	}
}

// ListHelmCharts는 ECR의 모든 Helm 차트 리포지토리를 조회하는 핸들러입니다.
// 예: GET /v1/helm-charts
func (h *HelmHandler) ListHelmCharts(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("request to list all helm chart repositories")

	charts, err := h.chartService.ListHelmCharts(r.Context())
	if err != nil {
		h.logger.Error("failed to list helm charts", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, charts)
}

// HealthCheck는 서비스의 상태를 확인하는 간단한 핸들러입니다.
// 200 OK 응답을 반환하여 서비스가 살아있음을 알립니다.
func (h *HelmHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// respondJSON은 JSON 응답을 작성하는 헬퍼 함수입니다.
func (h *HelmHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// respondError는 JSON 형식의 에러 응답을 작성하는 헬퍼 함수입니다.
func (h *HelmHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
