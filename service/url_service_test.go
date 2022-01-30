package service

// import (
// 	"testing"

// 	m "github.com/jayanpraveen/tildly/entity"
// )

// var hash = "MsQlS"
// var longUrl = "https://pkg.go.dev"

// var u = m.Url{
// 	Hash:      hash,
// 	LongUrl:   longUrl,
// 	CreatedAt: 1257894000,
// 	ExipreAt:  1357894000,
// }

// type MockCacheRepo struct {
// 	SetUrlFunc func(u *m.Url) error
// 	GetUrlFunc func(hash string) (*m.Url, error)
// }

// func (c *MockCacheRepo) SetUrl(u *m.Url) error {
// 	return nil
// }

// func (c *MockCacheRepo) GetUrl(hash string) (*m.Url, error) {
// 	return c.GetUrlFunc(hash)
// }

// func TestGetUrl(t *testing.T) {

// 	mcr := MockCacheRepo{
// 		GetUrlFunc: func(hash string) (url *m.Url, err error) {
// 			return &u, nil
// 		},
// 	}

// 	us := NewUrlService(&mcr, nil)
// 	u, err := us.GetUrlByHash(hash)

// 	t.Log(u)
// 	t.Log(err)

// }
