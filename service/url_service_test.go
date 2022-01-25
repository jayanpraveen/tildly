package service

import (
	"testing"
	"time"

	m "github.com/jayanpraveen/tildly/entity"
)

var hash = "MsQlS"
var longUrl = "https://pkg.go.dev"

var u = m.Url{
	Hash:      hash,
	LongUrl:   longUrl,
	CreatedAt: time.Now().Format("2006-01-02 15:04:05.000000"),
}

type MockCacheRepo struct {
	SetLongUrlFunc func(u *m.Url) error
	GetLongUrlFunc func(hash string) (*m.Url, error)
}

func (c *MockCacheRepo) SetLongUrl(u *m.Url) error {
	return nil
}

func (c *MockCacheRepo) GetLongUrl(hash string) (*m.Url, error) {
	return c.GetLongUrlFunc(hash)
}

func TestGetLongUrl(t *testing.T) {

	mcr := MockCacheRepo{
		GetLongUrlFunc: func(hash string) (url *m.Url, err error) {
			return &u, nil
		},
	}

	us := NewUrlService(&mcr, nil)

	u, err := us.GetUrlByHash(hash)

	t.Log(u)
	t.Log(err)

}
