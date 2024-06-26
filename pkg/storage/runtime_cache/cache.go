package runtimecache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	log *slog.Logger
	cache *cache.Cache
}

func New(log *slog.Logger) *Cache {
	cache := cache.New(5*time.Minute, 10*time.Minute)
	return &Cache{
		log: log,
		cache: cache,
	}
}

func (c *Cache) GetThumb(ctx context.Context, url string) ([]byte, error) {
	const op = "runtimecache.GetThumb"
	log := c.log.With(slog.String("op", op))

	raw, found := c.cache.Get(url)
	if !found {
		return []byte{}, nil
	}
	
	thumb, ok := raw.([]byte)
	if !ok {
		err := fmt.Errorf("type assertion to []byte failed")
		log.Error(err.Error())
		return []byte{}, err 
	}
	return thumb, nil
}

func (c *Cache) PutThumb(ctx context.Context, url string, thumb []byte) error {
	const op = "runtimecache.GetThumb"
	log := c.log.With(slog.String("op", op))

	c.cache.Set(url, thumb, cache.DefaultExpiration)
	log.Debug("New value cached")
	return nil
}

