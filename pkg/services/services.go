package services

import (
	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services/thumbs_service"
)

type GeneralService struct {
	thumbs_service.ThumbsService
}

func New() *GeneralService {
	return &GeneralService{}
}
