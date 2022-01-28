package entity

type Url struct {
	Hash      string `json:"hash"`
	LongUrl   string `json:"longUrl"`
	CreatedAt int64  `json:"createdAt"`
	ExipreAt  int64  `json:"exipreAt"`
}
type UrlStore interface {
	SaveUrl(u *Url, expireAt int64) error
	GetUrlByHash(hash string) (*Url, error)
}
