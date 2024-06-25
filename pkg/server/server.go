package server

import (
	"context"
	"log/slog"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/grpc/ytthumbs"
	"google.golang.org/grpc"
)

type ThumbsService interface {
	ParseUrls([]string) ([]string, error)
	GetImage(string) ([]byte, error)
}

type Server struct {
	ytthumbs.UnimplementedYouTubeThumbnailServiceServer
	service ThumbsService
	log *slog.Logger
}

var _ *Server = (*Server)(nil)

func Register(gRPCServer *grpc.Server, log *slog.Logger, service ThumbsService){
	ytthumbs.RegisterYouTubeThumbnailServiceServer(gRPCServer, &Server{service: service, log: log})
}

func (s *Server) GetThumbnails(ctx context.Context, r *ytthumbs.ThumbnailsRequest) (*ytthumbs.ThumbnailsResponse, error) {
	const op = "server.GetThumbnails"
	log := s.log.With(slog.String("op", op))

	urls := r.VideoUrls
	res := make([]*ytthumbs.Thumbnail, len(urls))
	for _, videoID := range urls {
		image, err := s.service.GetImage(videoID)
		if err != nil {
			log.Error(err.Error())
			return &ytthumbs.ThumbnailsResponse{}, err
		}
		res = append(res, &ytthumbs.Thumbnail{VideoUrl: videoID, Thumbnail: image, Error: ""})
	}

	return &ytthumbs.ThumbnailsResponse{Thumbnails: res}, nil
}

