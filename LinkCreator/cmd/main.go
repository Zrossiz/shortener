package main

import (
	"fmt"
	"github.com/Zrossiz/LinkCreator/creator/internal/delivery/rest"
	"github.com/Zrossiz/LinkCreator/creator/internal/repository/postgresql"
	redisdb "github.com/Zrossiz/LinkCreator/creator/internal/repository/redis"
	"github.com/Zrossiz/LinkCreator/creator/internal/service"
	"github.com/Zrossiz/LinkCreator/creator/pkg/config"
	"github.com/Zrossiz/LinkCreator/creator/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	log, err := logger.New(cfg.Server.LogLevel)
	if err != nil {
		fmt.Println("error init logger: ", err)
		os.Exit(1)
	}

	postgresConn, err := postgresql.Connect(cfg.Postgres.DBURI)
	if err != nil {
		log.Error("postgres connect error", zap.Error(err))
		os.Exit(1)
	}
	defer postgresConn.Close()

	redisConn, err := redisdb.Connect(cfg.Redis.Address, cfg.Redis.Password)
	if err != nil {
		log.Error("redis connect error", zap.Error(err))
		os.Exit(1)
	}
	defer redisConn.Close()

	postgresRepo := postgresql.NewPostgresRepo(postgresConn)
	redisRepo := redisdb.New(redisConn)

	serv := service.NewService(log, postgresRepo, redisRepo)

	hand := rest.NewUrlHandler(serv)

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: http.HandlerFunc(hand.CreateUrl),
	}

	log.Sugar().Infof("Starting server on %s", cfg.Server.Address)
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", zap.Error(err))
	}
}
