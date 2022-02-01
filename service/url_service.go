package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jayanpraveen/tildly/datastore"
	m "github.com/jayanpraveen/tildly/entity"
)

type UrlService struct {
	csdra UrlCsdra
	cache UrlCache
	AC    AtmoicCounterRepo
}

type UrlRepository interface {
	SaveUrl(longUrl string, expireAt int64) (string, error)
	GetUrlByHash(hash string) (*m.Url, error)
}

func NewUrlService(csdra UrlCsdra, uc UrlCache, es *datastore.EtcdStore) *UrlService {
	return &UrlService{
		csdra: csdra,
		cache: uc,
		AC:    NewAtomicCounter(es),
	}
}

func (s *UrlService) generateHash(longUrl string) (msg string) {
	md5hash := md5.New()
	md5hash.Write([]byte(fmt.Sprintf("%d_%s", s.AC.next(), longUrl)))
	return hex.EncodeToString(md5hash.Sum(nil))[:7]
}

func (s *UrlService) SaveUrl(longUrl string, exipreAt int64) (hash string, err error) {
	hash = s.generateHash(longUrl)
	u := m.Url{
		Hash:      hash,
		LongUrl:   longUrl,
		CreatedAt: time.Now().Unix(),
		ExipreAt:  exipreAt,
	}

	err = s.csdra.SetUrl(&u)
	if err != nil {
		return "csdra error", err
	}
	err = s.cache.SetUrl(&u)
	if err != nil {
		return "redis error", err
	}
	return hash, nil
}

func (s *UrlService) GetUrlByHash(hash string) (u *m.Url, err error) {
	u, err = s.cache.GetUrl(hash)
	if err == nil && !isUrlExpired(u) {
		return u, err
	}
	u, err = s.csdra.GetUrl(hash)
	if err == nil && !isUrlExpired(u) {
		return u, err
	}
	return nil, errors.New("csdra & cache error")
}

func isUrlExpired(u *m.Url) bool {
	if u.ExipreAt != 0 && u.ExipreAt > time.Now().Unix() {
		return true
	}
	return false
}
