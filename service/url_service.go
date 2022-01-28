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
	SaveUrl(longUrl string, expireAt int64) error
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

func (s *UrlService) SaveUrl(longUrl string, exipreAt int64) error {

	hash := s.generateHash(longUrl)

	u := m.Url{
		Hash:      hash,
		LongUrl:   longUrl,
		CreatedAt: time.Now().Unix(),
		ExipreAt:  exipreAt,
	}

	// Save to cache, db...
	s.cache.SetLongUrl(&u)

	return nil

}

func (s *UrlService) GetUrlByHash(hash string) (u *m.Url, err error) {
	u, err = s.cache.GetLongUrl(hash)
	if u.ExipreAt != 0 && u.ExipreAt > time.Now().Unix() {
		return nil, err
	}
	return u, err
}
