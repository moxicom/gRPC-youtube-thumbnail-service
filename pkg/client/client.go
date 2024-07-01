package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/grpc/ytthumbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	client ytthumbs.YouTubeThumbnailServiceClient
	urls   []string
}

func New(log *slog.Logger, connectionStr string, urls []string) Client {
	cc, err := grpc.NewClient(connectionStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := ytthumbs.NewYouTubeThumbnailServiceClient(cc)

	return Client{log, client, urls}
}

func (c *Client) FetchThumbnailsAsync(ctx context.Context) {
	log := c.log.With(slog.String("op", "client.FetchThumbnailsAsync"))

	log.Info("received urls. fetching asynchronously")

	var wg sync.WaitGroup
	var mu sync.Mutex
	amount := len(c.urls)

	for i, videoUrl := range c.urls {
		wg.Add(1)
		go func(i int, videoUrl string) {
			defer wg.Done()
			req := &ytthumbs.ThumbnailsRequest{VideoUrls: []string{videoUrl}}
			resp, err := c.client.GetThumbnails(ctx, req)
			if err != nil {
				log.Error("failed to get thumbnail", slog.String("url", videoUrl))
				return
			}

			log.Debug(
				"received thumbnail",
				slog.String("url", videoUrl),
				slog.String("progress", string(i+1)+string(amount)))

			mu.Lock()
			c.SaveThumbnails(resp.Thumbnails)
			mu.Unlock()
		}(i, videoUrl)
	}

	wg.Wait()
}

func (c *Client) FetchThumbnails(ctx context.Context) {
	log := c.log.With(slog.String("op", "client.FetchThumbnails"))

	log.Info("received urls. fetching")

	req := &ytthumbs.ThumbnailsRequest{VideoUrls: c.urls}
	resp, err := c.client.GetThumbnails(ctx, req)
	if err != nil {
		log.Error("Failed not get thumbnails", slog.Any("error", err))
		return
	}
	c.SaveThumbnails(resp.Thumbnails)

	log.Info("Thumbnails successfully saved")
}

func (c *Client) SaveThumbnails(thumbs []*ytthumbs.Thumbnail) {
	log := c.log.With(slog.String("op", "client.FetchThumbnails"))

	const path = "./output/"
	for _, thumb := range thumbs {
		log.Debug("saving thumbnail", slog.String("url", thumb.VideoUrl))
		fileName := fmt.Sprintf("%v.jpg", thumb.VideoUrl)
		filePath := filepath.Join(path, fileName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModeDir|0755)
			if err != nil {
				log.Error(
					"Failed to make a dir to save",
					slog.String("path", path),
					slog.Any("err", err),
				)
				return
			}
		}
		err := os.WriteFile(filePath, thumb.Thumbnail, 0644)
		if err != nil {
			log.Error(
				"Failed to save thumbnail",
				slog.String("url", thumb.VideoUrl),
				slog.Any("err", err),
			)
			return
		}
		log.Debug("Thumbnail saved successfully", slog.String("url", thumb.VideoUrl))
	}
}
