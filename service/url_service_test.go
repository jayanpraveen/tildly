package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	m "github.com/jayanpraveen/tildly/entity"
)

var hash = "MsQlS"
var longUrl = "https://pkg.go.dev"

var u = m.Url{
	Hash:      hash,
	LongUrl:   longUrl,
	CreatedAt: 1257894000,
	ExipreAt:  0,
}

type MockCsdraRepo struct {
	SetUrlFunc func(u *m.Url) error
	GetUrlFunc func(hash string) (*m.Url, error)
}

func (c *MockCsdraRepo) SetUrl(u *m.Url) error {
	return c.SetUrlFunc(u)
}

func (c *MockCsdraRepo) GetUrl(hash string) (*m.Url, error) {
	return c.GetUrlFunc(hash)
}

type MockCacheRepo struct {
	SetUrlFunc func(u *m.Url) error
	GetUrlFunc func(hash string) (*m.Url, error)
}

func (c *MockCacheRepo) SetUrl(u *m.Url) error {
	return c.SetUrlFunc(u)
}

func (c *MockCacheRepo) GetUrl(hash string) (*m.Url, error) {
	return c.GetUrlFunc(hash)
}

type MockAtomicCounter struct {
	nextFunc                func() int
	DisplayCurrentRangeFunc func() string
}

func (ac *MockAtomicCounter) next() int {
	return ac.nextFunc()
}

func (ac *MockAtomicCounter) DisplayCurrentRange() string {
	return ac.DisplayCurrentRangeFunc()
}

func TestGetUrl(t *testing.T) {

	mac := MockAtomicCounter{
		nextFunc: func() int {
			return 100
		},
	}

	t.Run("get url from cache", func(t *testing.T) {

		cacheRepo := MockCacheRepo{
			GetUrlFunc: func(hash string) (url *m.Url, err error) {
				return &u, nil
			},
		}

		csdraRepo := MockCsdraRepo{
			GetUrlFunc: func(hash string) (*m.Url, error) {
				return &u, nil
			},
		}

		us := UrlService{
			csdra: &csdraRepo,
			cache: &cacheRepo,
			AC:    &mac,
		}

		u, err := us.GetUrlByHash(hash)

		t.Log(u)
		t.Log(err)
	})

	t.Run("get url from cassandra", func(t *testing.T) {

		cacheRepo := MockCacheRepo{
			GetUrlFunc: func(hash string) (url *m.Url, err error) {
				return nil, errors.New("cache error")
			},
		}

		csdraRepo := MockCsdraRepo{
			GetUrlFunc: func(hash string) (*m.Url, error) {
				return &u, nil
			},
		}

		us := UrlService{
			csdra: &csdraRepo,
			cache: &cacheRepo,
			AC:    &mac,
		}

		res, _ := us.GetUrlByHash(hash)
		if cmp.Equal(res, u) {
			t.Error("hash don't match")
		}
	})

	t.Run("csdra cache error", func(t *testing.T) {

		cacheRepo := MockCacheRepo{
			GetUrlFunc: func(hash string) (url *m.Url, err error) {
				return nil, errors.New("cache error")
			},
		}

		csdraRepo := MockCsdraRepo{
			GetUrlFunc: func(hash string) (*m.Url, error) {
				return nil, errors.New("cache error")
			},
		}

		us := UrlService{
			csdra: &csdraRepo,
			cache: &cacheRepo,
			AC:    &mac,
		}

		_, err := us.GetUrlByHash(hash)
		if err == errors.New("csdra & cache error") {
			t.Error("failing: error thown")
		}

	})

}

func TestSetUrl(t *testing.T) {

	mac := MockAtomicCounter{
		nextFunc: func() int {
			return 100
		},
	}

	t.Run("set url pass", func(t *testing.T) {

		cacheRepo := MockCacheRepo{
			SetUrlFunc: func(u *m.Url) error {
				return nil
			},
		}

		csdraRepo := MockCsdraRepo{
			SetUrlFunc: func(u *m.Url) error {
				return nil
			},
		}

		us := UrlService{
			csdra: &cacheRepo,
			cache: &csdraRepo,
			AC:    &mac,
		}

		hash, err := us.SaveUrl(u.Hash, u.ExipreAt)

		if hash != "a5fd263" {
			t.Errorf("exp: a5fd263 act:%s", hash)
		}

		if err != nil {
			t.Error(err)
		}

	})

	t.Run("set url fails", func(t *testing.T) {

		cacheRepo := MockCacheRepo{
			SetUrlFunc: func(u *m.Url) error {
				return fmt.Errorf("cache url fails")
			},
		}

		csdraRepo := MockCsdraRepo{
			SetUrlFunc: func(u *m.Url) error {
				return fmt.Errorf("csdra url fails")
			},
		}

		us := UrlService{
			csdra: &cacheRepo,
			cache: &csdraRepo,
			AC:    &mac,
		}

		_, err := us.SaveUrl(u.Hash, u.ExipreAt)

		if err == nil {
			t.Error("failing save fails")
		}

	})

}
