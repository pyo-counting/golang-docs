// filepath: helm-ecr-api/cmd/api/main.go
package main

import (
	"context"
	"errors"
	"helm-ecr-api/internal/handler"
	"helm-ecr-api/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// 1. 로거 설정
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 환경 변수에서 허용할 Helm 리포지토리 목록을 읽어옵니다.
	allowedReposStr := os.Getenv("HELM_REPOSITORIES")
	if allowedReposStr == "" {
		logger.Error("HELM_REPOSITORIES environment variable must be set")
		os.Exit(1)
	}
	allowedRepos := strings.Split(allowedReposStr, ",")

	// 2. AWS 설정 로드
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Error("failed to load AWS configuration", "error", err)
		os.Exit(1)
	}

	// 3. 서비스 및 핸들러 계층 초기화 (의존성 주입)
	ecrSvc := service.NewECRService(cfg, allowedRepos)
	helmHandler := handler.NewHelmHandler(ecrSvc, logger)

	// 4. 라우터 설정
	mux := http.NewServeMux()

	// RESTful API 경로 설계 (개선된 버전)
	// /v1/helm-charts/ 로 시작하는 모든 요청을 RouteChartDetails 핸들러로 위임합니다.
	mux.HandleFunc("GET /v1/helm-charts", helmHandler.ListHelmCharts)
	mux.HandleFunc("GET /v1/helm-charts/", helmHandler.RouteChartDetails) // Prefix-based routing
	mux.HandleFunc("GET /health", helmHandler.HealthCheck)

	// 환경 변수에서 포트를 읽어오고, 설정되지 않은 경우 기본값 8080을 사용합니다.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// 5. Graceful Shutdown과 함께 서버 시작
	go func() {
		logger.Info("starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited properly")
}
