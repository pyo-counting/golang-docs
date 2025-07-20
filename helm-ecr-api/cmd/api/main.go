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
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// 1. 로거 설정
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. AWS 설정 로드
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Error("failed to load AWS configuration", "error", err)
		os.Exit(1)
	}

	// 3. 서비스 및 핸들러 계층 초기화 (의존성 주입)
	ecrSvc := service.NewECRService(cfg)
	helmHandler := handler.NewHelmHandler(ecrSvc, logger)

	// 4. 라우터 설정
	mux := http.NewServeMux()

	// RESTful API 경로 설계 (개선된 버전)
	// - `GET /v1/helm-charts`: 모든 차트 저장소 목록 조회
	// - `GET /v1/helm-charts/{chart-name}`: 특정 차트의 모든 태그 목록 조회
	// - `GET /v1/helm-charts/{chart-name}?tag=...`: 특정 차트의 특정 태그 정보 조회
	mux.HandleFunc("GET /v1/helm-charts", helmHandler.ListHelmCharts)
	mux.HandleFunc("GET /v1/helm-charts/", helmHandler.GetHelmChart) // 후행 슬래시로 하위 경로 처리
	mux.HandleFunc("GET /health", helmHandler.HealthCheck)

	server := &http.Server{
		Addr:    ":8080",
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
