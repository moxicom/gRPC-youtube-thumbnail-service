package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/client"
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/logger"
)

var (
	// serverAddr string
	async      bool
	envLog     string
	serverAddr string
)

func init() {
	flag.StringVar(&serverAddr, "server_addr", "localhost:8080", "The server address in the format of host:port")
	flag.BoolVar(&async, "async", false, "Download thumbnails asynchronously")
	flag.StringVar(
		&envLog,
		"envLog",
		logger.EnvProd,
		fmt.Sprintf("'%s' or '%s' to setup logger", logger.EnvProd, logger.EnvLocal),
	)
}

func main() {
	flag.Parse()

	// Remaining command-line arguments are treated as video URLs
	videoUrls := flag.Args()
	if len(videoUrls) == 0 {
		flag.Usage()
		panic("No video URLs provided")
	}
	log := logger.SetupLogger(envLog)
	client := client.New(log, serverAddr, videoUrls)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	start := time.Now()

	if async {
		client.FetchThumbnailsAsync(ctx)
	} else {
		client.FetchThumbnails(ctx)
	}
	fmt.Println()
	fmt.Println(time.Since(start))
}
