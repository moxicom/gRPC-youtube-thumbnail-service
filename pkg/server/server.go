package server

import (
	"context"
	"log/slog"
	"sync"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/grpc/ytthumbs"
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThumbsService interface {
	ParseUrls([]string) ([]string, error)
	GetImage(context.Context, string) ([]byte, error)
}

type Server struct {
	ytthumbs.UnimplementedYouTubeThumbnailServiceServer
	service ThumbsService
	log     *slog.Logger
}

var _ ytthumbs.YouTubeThumbnailServiceServer = (*Server)(nil)

func Register(gRPCServer *grpc.Server, log *slog.Logger, service ThumbsService) {
	ytthumbs.RegisterYouTubeThumbnailServiceServer(gRPCServer, &Server{service: service, log: log})
}

func (s *Server) GetThumbnails(ctx context.Context, r *ytthumbs.ThumbnailsRequest) (*ytthumbs.ThumbnailsResponse, error) {
	const op = "server.GetThumbnails"
	log := s.log.With(slog.String("op", op))

	log.Info("New request received")

	// Parse the URLs to extract video IDs.
	videoIDs, err := s.service.ParseUrls(r.VideoUrls)
	if err != nil {
		log.Error("failed to parse video urls", slog.Any("err", err))
		return &ytthumbs.ThumbnailsResponse{}, status.Error(codes.InvalidArgument, services.ErrBadURL.Error())
	}
	log.Debug("video url parsed successfully")

	res := make([]*ytthumbs.Thumbnail, len(videoIDs))
	errChan := make(chan error, len(videoIDs))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, videoID := range videoIDs {
		wg.Add(1)
		go func(i int, videoID string) {
			defer wg.Done()
			image, err := s.service.GetImage(ctx, videoID)
			if err != nil {
				if err == services.ErrVideoNotFound {
					errChan <- status.Error(codes.NotFound, "video not found with id="+videoID)
					return
				}
				errChan <- status.Error(codes.Internal, "internal server error")
				return
			}
			mu.Lock()
			res[i] = &ytthumbs.Thumbnail{VideoUrl: videoID, Thumbnail: image}
			mu.Unlock()
		}(i, videoID)
	}

	// Wait for all goroutines to complete.
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check if any errors occurred.
	for err := range errChan {
		if err != nil {
			log.Error("Error fetching thumbnails", slog.Any("err", err))
			return nil, err
		}
	}

	return &ytthumbs.ThumbnailsResponse{Thumbnails: res}, nil
}
