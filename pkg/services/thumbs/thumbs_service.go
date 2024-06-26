package thumbs_service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/services"
)

type Storage interface {
	GetThumb(context.Context, string) ([]byte, error)
	PutThumb(context.Context, string, []byte) error
}

type ThumbsService struct {
	storage Storage
	log *slog.Logger
}

func New(log *slog.Logger, storage Storage) *ThumbsService{
	return &ThumbsService{storage, log}
}

func (s *ThumbsService) ParseUrls(urls []string) ([]string, error) {
	ids := make([]string, len(urls))
	for i, url := range urls {
		re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([^&]+)`)
    	matches := re.FindStringSubmatch(url)
		if len(matches) < 2 {
			return []string{}, services.ErrBadURL
		}
		ids[i] = matches[1]
	}
	return ids, nil
}

// GetImage - get image from storage or receive it from api and write to cache
func (s *ThumbsService) GetImage(ctx context.Context, url string) ([]byte, error) {
	const op = "thumbs_service.GetImage"
	log := s.log.With(slog.String("op", op))

	thumb, err := s.storage.GetThumb(ctx, url)
	if err != nil {
		return []byte{}, err
	}

	if len(thumb) == 0 {
		thumb, err = s.requestImageThumb(ctx, url)
		if err != nil {
			return []byte{}, err
		}
		log.Debug("Got from API")	
		err := s.storage.PutThumb(ctx, url, thumb)
		if err != nil {
			return []byte{}, err
		}
	} else {
		log.Debug("Found in storage")	
	}
	return thumb, nil

}

func (s *ThumbsService) requestImageThumb(ctx context.Context, url string) ([]byte, error) {
	const op = "thumbs_service.RequestImageThumb"
	log := s.log.With(slog.String("op", op))

	requestUrl := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
    if err != nil {
		log.Error("%s error creating request: %s", op, err.Error())
        return []byte{}, err
    }

    response, err := http.DefaultClient.Do(req)
    if err != nil {
		log.Error("%s error when requesting an image: %s", op, err.Error())
        return []byte{}, err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s response code is not ok: %d", op, response.StatusCode)
		log.Error(err.Error())
        return []byte{}, err 
    }

    data, err := io.ReadAll(response.Body)
    if err != nil {
		log.Error("%s error while reading response body: %s", op, err.Error())
        return []byte{}, err
    }

    return data, nil
}