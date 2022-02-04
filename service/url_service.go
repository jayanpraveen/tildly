package service

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
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

func genMD5(str string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(str))
	return hex.EncodeToString(md5hash.Sum(nil))
}

func genBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (s *UrlService) genHash(longUrl string) (msg string) {
	c := strconv.Itoa(s.AC.next())
	p1 := genBase64(genMD5(fmt.Sprintf("%s_%s", longUrl, c)))
	p2 := genBase64(genMD5(c))
	return p1[:3] + p2[:4]
}

func (s *UrlService) SaveUrl(longUrl string, exipreAt int64) (hash string, err error) {
	hash = s.genHash(longUrl)

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
