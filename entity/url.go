package entity

type Url struct {
	Hash      string `json:"hash"`
	LongUrl   string `json:"longUrl"`
	CreatedAt string `json:"createdAt"`
}
type UrlStore interface {
	SaveUrl(u Url) error
	GetUrlByHash(hash string) (Url, error)
}
