package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/QosmuratSamat0/catalog/internal/app/grpc"
	"github.com/QosmuratSamat0/catalog/internal/cache"
	"github.com/QosmuratSamat0/catalog/internal/repository/postgresql"
	"github.com/QosmuratSamat0/catalog/internal/services/product"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, redisAddr string, tokenTTL time.Duration) *App {
	storage, err := postgresql.New(storagePath)
	if err != nil {
		panic(err)
	}

	redisCache := cache.New(redisAddr)

	productService := product.New(storage, redisCache)

	grpcApp := grpcapp.New(log, productService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
