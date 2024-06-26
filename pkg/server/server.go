package server

import (
	"context"
	"log/slog"

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
	log *slog.Logger
}

var _ *Server = (*Server)(nil)

func Register(gRPCServer *grpc.Server, log *slog.Logger, service ThumbsService){
	ytthumbs.RegisterYouTubeThumbnailServiceServer(gRPCServer, &Server{service: service, log: log})
}

func (s *Server) GetThumbnails(ctx context.Context, r *ytthumbs.ThumbnailsRequest) (*ytthumbs.ThumbnailsResponse, error) {
	const op = "server.GetThumbnails"
	log := s.log.With(slog.String("op", op))

	log.Info("New request received")

	urls, err := s.service.ParseUrls(r.VideoUrls)
	if err != nil {
		return &ytthumbs.ThumbnailsResponse{}, status.Error(codes.InvalidArgument, services.ErrBadURL.Error())
	}
	res := make([]*ytthumbs.Thumbnail, len(urls))
	for i, videoID := range urls {
		image, err := s.service.GetImage(ctx, videoID)
		if err != nil {
			return &ytthumbs.ThumbnailsResponse{}, status.Error(codes.Internal, "internal server error")
		}
		res[i] = &ytthumbs.Thumbnail{VideoUrl: r.VideoUrls[i], Thumbnail: image}
	}

	return &ytthumbs.ThumbnailsResponse{Thumbnails: res}, nil
}

