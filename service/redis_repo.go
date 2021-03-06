package service

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	m "github.com/jayanpraveen/tildly/entity"
)

type UrlCache interface {
	SetUrl(url *m.Url) error
	GetUrl(hash string) (*m.Url, error)
}

// CacheRepo implements the interface UrlCache
type CacheRepo struct {
	cache *cache.Cache
}

func NewCacheRepo(c *cache.Cache) *CacheRepo {
	return &CacheRepo{
		cache: c,
	}
}

func (c *CacheRepo) SetUrl(u *m.Url) error {
	ctx := context.Background()

	if err := c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   u.Hash,
		Value: u,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}
	return nil
}

func (c *CacheRepo) GetUrl(hash string) (*m.Url, error) {
	var u m.Url
	ctx := context.Background()
	if err := c.cache.Get(ctx, hash, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
