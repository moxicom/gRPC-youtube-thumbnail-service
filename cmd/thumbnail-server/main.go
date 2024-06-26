package main

import (
	"log/slog"
	"os"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/app"
	thumbs_service "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services/thumbs"
	runtimecache "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/storage/runtime_cache"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	log := setupLogger("local")
	storage := runtimecache.New(log)
	service := thumbs_service.New(log, storage)
	app := app.New(log, service, 8080)

	if err := app.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return log
}