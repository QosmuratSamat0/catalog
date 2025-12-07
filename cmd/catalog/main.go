package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/QosmuratSamat0/product-catalog/internal/app"
	"github.com/QosmuratSamat0/product-catalog/internal/config"
	"github.com/QosmuratSamat0/product-catalog/internal/pkg/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting product-catalog Service",
		slog.String("env", cfg.Env),
		slog.Int("grpc_port", cfg.GRPC.Port),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.RedisAddr, cfg.TokenTTL)

	go application.GRPCSrv.MustLoad()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Shutting down...", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
