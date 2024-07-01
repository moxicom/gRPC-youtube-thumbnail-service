package main

import (
	"flag"
	"fmt"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/app"
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/logger"
	thumbs_service "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services/thumbs"
	runtimecache "github.com/moxicom/grpc-youtube-thumbnail-service/pkg/storage/runtime_cache"
)

var envLog string

func init() {
	flag.StringVar(
		&envLog,
		"envLog",
		logger.EnvProd,
		fmt.Sprintf("'%s' or '%s' to setup logger", logger.EnvProd, logger.EnvLocal),
	)
}

func main() {
	log := logger.SetupLogger(envLog)
	storage := runtimecache.New(log)
	service := thumbs_service.New(log, storage)
	app := app.New(log, service, 8080)

	if err := app.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}
