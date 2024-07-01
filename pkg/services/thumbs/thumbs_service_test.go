package thumbs_service

import (
	"context"
	"testing"

	"github.com/moxicom/grpc-youtube-thumbnail-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	getThumbFunc func(ctx context.Context, videoID string) ([]byte, error)
	putThumbFunc func(ctx context.Context, videoID string, data []byte) error
}

func (m *mockStorage) GetThumb(ctx context.Context, videoID string) ([]byte, error) {
	return m.getThumbFunc(ctx, videoID)
}

func (m *mockStorage) PutThumb(ctx context.Context, videoID string, data []byte) error {
	return m.putThumbFunc(ctx, videoID, data)
}

func TestParseUrls(t *testing.T) {
	log := logger.SetupLogger(logger.EnvLocal)
	service := New(log, nil)

	testCases := []struct {
		name        string
		urls        []string
		expectedIDs []string
		expectError bool
	}{
		{
			name:        "Valid URLs",
			urls:        []string{"https://youtu.be/dQw4w9WgXcQ", "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
			expectedIDs: []string{"dQw4w9WgXcQ", "dQw4w9WgXcQ"},
			expectError: false,
		},
		{
			name:        "Invalid URL",
			urls:        []string{"https://invalid.url"},
			expectedIDs: nil,
			expectError: true,
		},
		{
			name:        "Mixed URLs",
			urls:        []string{"https://youtu.be/dQw4w9WgXcQ", "https://invalid.url"},
			expectedIDs: nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ids, err := service.ParseUrls(tc.urls)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedIDs, ids)
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()
	log := logger.SetupLogger(logger.EnvLocal)

	mockStorage := &mockStorage{
		getThumbFunc: func(ctx context.Context, videoID string) ([]byte, error) {
			if videoID == "11" {
				return []byte("image data"), nil
			} else {
				return []byte{}, nil
			}
		},
		putThumbFunc: func(ctx context.Context, videoID string, data []byte) error {
			return nil
		},
	}

	service := New(log, mockStorage)

	testCases := []struct {
		name         string
		url          string
		expectedData []byte
		expectError  error
	}{
		{
			name:         "Thumbnail in storage",
			url:          "11",
			expectedData: []byte("image data"),
			expectError:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run("Thumbnail in storage", func(t *testing.T) {
			thumb, err := service.GetImage(ctx, "11")

			if testCase.expectError == nil {
				assert.Equal(t, testCase.expectedData, thumb)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, testCase.expectError)
			}
		})
	}
}
