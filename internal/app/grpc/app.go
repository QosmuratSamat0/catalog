package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/grpc"

	pService "github.com/QosmuratSamat0/product-catalog/internal/grpc/product"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, productService pService.Product, port int) *App {
	gRPCServer := grpc.NewServer()
	pService.Register(gRPCServer, productService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) MustLoad() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.String("port", strconv.Itoa(app.port)),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("grpc server listening on " + l.Addr().String())
	log.With(
		slog.String("op", op),
		slog.String("port", strconv.Itoa(app.port)),
	)
	if err := app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"
	app.log.With("op", op, "port", app.port).Info("Grpc server stopping")
	app.gRPCServer.GracefulStop()
}
