package service

import (
	"github.com/gocql/gocql"
	m "github.com/jayanpraveen/tildly/entity"
)

type UrlCsdra interface {
	SetUrl(url *m.Url) error
	GetUrl(hash string) (*m.Url, error)
}

type Cassandra struct {
	cs *gocql.Session
}

func NewCassandra(cs *gocql.Session) *Cassandra {
	return &Cassandra{
		cs: cs,
	}
}

func (c *Cassandra) SetUrl(u *m.Url) error {
	return c.cs.Query(
		"INSERT INTO tildly.url (hash, createdat, exipireat, longurl) VALUES (?, ?, ?, ?)",
		u.Hash, u.CreatedAt, u.ExipreAt, u.LongUrl).Exec()

}

func (c *Cassandra) GetUrl(hash string) (*m.Url, error) {
	var u m.Url
	err := c.cs.Query("SELECT * FROM tildly.url WHERE hash = ?", hash).
		Scan(&u.Hash, &u.CreatedAt, &u.ExipreAt, &u.LongUrl)
	return &u, err
}
