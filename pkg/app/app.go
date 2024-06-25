package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	thumbsService server.ThumbsService,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	server.Register(gRPCServer, log, thumbsService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	log := a.log.With(slog.String("op", op))
	log.Info("Starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reflection.Register(a.gRPCServer)

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}