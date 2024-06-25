package thumbs_service

import (
	"fmt"
	"io"
	"net/http"
)

type ThumbsService struct {}

func New() *ThumbsService{
	return &ThumbsService{}
}

func (s *ThumbsService) ParseUrls([]string) ([]string, error) {
	return []string{}, nil
}

func (s *ThumbsService) GetImage(url string) ([]byte, error) {
	const op = "thumbs_service.GetImage"

	requestUrl := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", url)
	fmt.Println(requestUrl)

	response, err := http.Get(requestUrl)
	if err != nil {
		return []byte{}, fmt.Errorf("%s error on image request: %s", op, err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
        return []byte{}, fmt.Errorf("%s response code is not ok:  %d", op,  response.StatusCode)
    }

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%s error while reading response body: %s", op, err.Error())
	}

	return data, nil
}