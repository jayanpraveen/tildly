package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jayanpraveen/tildly/datastore"
	m "github.com/jayanpraveen/tildly/entity"
)

type UrlService struct {
	cache UrlCache
	AC    *atomicCounter
}

type UrlRepository interface {
	SaveUrl(longUrl string) error
	GetUrlByHash(hash string) (*m.Url, error)
}

func NewUrlService(uc UrlCache, es *datastore.EtcdStore) *UrlService {
	return &UrlService{
		cache: uc,
		AC:    NewAtomicCounter(es),
	}
}

func (s *UrlService) generateHash(longUrl string) (hash string) {
	md5hash := md5.New()
	md5hash.Write([]byte(fmt.Sprintf("%d_%s", s.AC.next(), longUrl)))
	return hex.EncodeToString(md5hash.Sum(nil))[:7]
}

func (s *UrlService) SaveUrl(longUrl string) error {

	hash := s.generateHash(longUrl)

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
