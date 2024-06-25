package main

import (
	"log/slog"
	"os"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/app"
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

// https://img.youtube.com/vi/l_o0izWJeCA/0.jpg

func main() {
	log := setupLogger("local")

	service := services.New()

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