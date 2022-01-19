package service

import (
	"fmt"
	"time"

	m "github.com/jayanpraveen/tildly/entity"
)

type UrlService struct {
	cache UrlCache
}

type UrlRepository interface {
	SaveUrl(longUrl string) error
	GetUrlByHash(hash string) (*m.Url, error)
}

func NewUrlService(uc UrlCache) *UrlService {
	return &UrlService{
		cache: uc,
	}
}

// !Change this
func (s *UrlService) SaveUrl(longUrl string) error {
	hash := "dQw4w9WgXcQ"
	fmt.Println("hash: ", hash)
	fmt.Println("longUrl: ", longUrl)

	u := m.Url{
		Hash:      hash,
		LongUrl:   longUrl,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.000000"),
	}

	// Save to cache, db...
	s.cache.SetLongUrl(&u)

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
