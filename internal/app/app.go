package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zrossiz/shortener/internal/config"
	"github.com/Zrossiz/shortener/internal/service"
	"github.com/Zrossiz/shortener/internal/storage"
	"github.com/Zrossiz/shortener/internal/transport/http/handler"
	"github.com/Zrossiz/shortener/internal/transport/http/router"
	logger "github.com/Zrossiz/shortener/pkg/log"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("error init config: ", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println("error init logger: ", err)
		os.Exit(1)
	}

	conn, err := storage.Connect(cfg.DBURI)
	if err != nil {
		log.Fatal("error connect", zap.Error(err))
	}
	defer conn.Close()

	store := storage.New(conn)

	serv := service.NewService(store, log)

	h := handler.NewHandler(&serv)

	r := router.NewRouter(h)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go startServer(cfg.ServerPort, r, log, stop)

	<-stop
	log.Info("Shutting down gracefully...")
}

func startServer(addr string, h http.Handler, log *zap.Logger, stop chan os.Signal) {
	server := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	log.Sugar().Infof("Starting server on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Failed to start server", zap.Error(err))
	}

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Error("Server Shutdown failed", zap.Error(err))
	} else {
		log.Info("Server stopped gracefully")
	}
}
