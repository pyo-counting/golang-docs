package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Item은 API에서 사용될 데이터 구조체입니다.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items []Item
var logger *slog.Logger

// getItems는 모든 아이템을 반환합니다.
// 참조: https://pkg.go.dev/net/http#HandleFunc
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		logger.Error("Failed to encode items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// createItem은 새 아이템을 생성합니다.
// 참조: https://pkg.go.dev/encoding/json#Decoder.Decode
func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	// 요청 본문을 디코딩하고 에러를 처리합니다.
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	items = append(items, item)
	logger.Info("New item created", "id", item.ID, "name", item.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created 상태 코드를 반환합니다.
	if err := json.NewEncoder(w).Encode(item); err != nil {
		logger.Error("Failed to encode created item", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// 1. 구조화된 로거(slog) 설정
	// DevOps/SRE 환경에서는 JSON 형식의 구조화된 로그가 관찰가능성에 매우 중요합니다.
	// 참조: https://pkg.go.dev/log/slog
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 초기 데이터
	items = append(items, Item{ID: "1", Name: "Item 1"})
	items = append(items, Item{ID: "2", Name: "Item 2"})

	// 2. 라우터 및 서버 설정
	// ServeMux를 사용해 라우팅을 명시적으로 관리합니다.
	mux := http.NewServeMux()
	mux.HandleFunc("/items", getItems)
	mux.HandleFunc("/items/new", createItem)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 3. Graceful Shutdown 구현
	// 서버가 SIGINT 또는 SIGTERM 신호를 받으면, 진행 중인 요청을 처리한 후 안전하게 종료됩니다.
	// 이는 무중단 배포(Zero-downtime deployment)의 핵심 요소입니다.
	// 참조: https://pkg.go.dev/net/http#Server.Shutdown
	go func() {
		logger.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 종료 신호를 기다립니다.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 30초의 타임아웃으로 컨텍스트를 생성합니다.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 서버를 안전하게 종료합니다.
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited properly")
}
