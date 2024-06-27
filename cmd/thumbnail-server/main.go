package main

import (
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/app"
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/logger"
	thumbs_service "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services/thumbs"
	runtimecache "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/storage/runtime_cache"
)


func main() {
	log := logger.SetupLogger("prod")
	storage := runtimecache.New(log)
	service := thumbs_service.New(log, storage)
	app := app.New(log, service, 8080)

	if err := app.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

