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

const (
	ytThumbURLPattern = `(?:youtube\.com/watch\?v=|youtu\.be/)([^&]+)`
	ytThumbURLFormat  = "https://i.ytimg.com/vi/%s/maxresdefault.jpg"
)

type Storage interface {
	GetThumb(context.Context, string) ([]byte, error)
	PutThumb(context.Context, string, []byte) error
}

type ThumbsService struct {
	storage Storage
	log     *slog.Logger
}

func New(log *slog.Logger, storage Storage) *ThumbsService {
	return &ThumbsService{storage, log}
}

func (s *ThumbsService) ParseUrls(urls []string) ([]string, error) {
	const op = "thumbs_service.ParseUrls"
	log := s.log.With(slog.String("op", op))

	re := regexp.MustCompile(ytThumbURLPattern)
	ids := make([]string, len(urls))

	for i, url := range urls {
		matches := re.FindStringSubmatch(url)
		if len(matches) < 2 {
			log.Error("Bad url", slog.String("url", url))
			return []string{}, services.ErrBadURL
		}
		ids[i] = matches[1]
	}

	return ids, nil
}

// GetImage - get image from storage or receive it from api and write to cache
func (s *ThumbsService) GetImage(ctx context.Context, videoID string) ([]byte, error) {
	const op = "thumbs_service.GetImage"
	log := s.log.With(slog.String("op", op))

	thumb, err := s.storage.GetThumb(ctx, videoID)
	if err != nil {
		return []byte{}, err
	}

	if len(thumb) == 0 {
		log.Debug("Thumbnail not found in storage, requesting from API", slog.String("videoID", videoID))
		thumb, err = s.requestImageThumb(ctx, videoID)
		if err != nil {
			return []byte{}, err
		}
		log.Debug("Got thumbnail from API", slog.String("videoID", videoID))
		err := s.storage.PutThumb(ctx, videoID, thumb)
		if err != nil {
			return []byte{}, err
		}
	} else {
		log.Debug("Thumbnail found in storage", slog.String("videoID", videoID))
	}
	return thumb, nil

}

func (s *ThumbsService) requestImageThumb(ctx context.Context, videoID string) ([]byte, error) {
	const op = "thumbs_service.RequestImageThumb"
	log := s.log.With(slog.String("op", op))

	log.Debug("Requesing thumbnail from api", slog.String("videoID", videoID))

	requestURL := fmt.Sprintf(ytThumbURLFormat, videoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		log.Error("error creating request", slog.Any("err", err))
		return []byte{}, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Error requesting thumbnail", slog.String("url", requestURL), slog.Any("err", err))
		return []byte{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected response code: %d", response.StatusCode)
		log.Error(err.Error())
		return []byte{}, services.ErrVideoNotFound
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Error reading response body", slog.String("url", requestURL), slog.Any("error", err))
		return []byte{}, err
	}

	log.Debug("Got thumbnail from api", slog.String("videoID", videoID))

	return data, nil
}
