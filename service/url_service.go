package service

import (
	"fmt"

	m "github.com/jayanpraveen/tildly/entity"
)

type UrlCache interface {
	SetLongUrl(url *m.Url) error
	GetLongUrl(hash string) (*m.Url, error)
}

type UrlService struct {
	cache UrlCache
}

func NewUrlService(uc UrlCache) *UrlService {
	return &UrlService{
		cache: uc,
	}
}

// !Change this
func (s *UrlService) SaveUrl(u *m.Url) error {
	hash := "dQw4w9WgXcQ"
	fmt.Println("hash: ", hash)

	// Save to cache, db...
	s.cache.SetLongUrl(u)

	return nil

}

func (s *UrlService) GetUrlByHash(hash string) (*m.Url, error) {

	var u *m.Url

	u, err := s.cache.GetLongUrl(hash)

	if err != nil {
		return nil, err
	}

	return u, nil
}
