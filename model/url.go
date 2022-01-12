package model

type Url struct {
	Hash      string `json:"hash"`
	LongUrl   string `json:"longUrl"`
	CreatedAt string `json:"createdAt"`
}

type UrlRepository interface {
	CreateHashForUrl(string) (string, error)
	GetLongUrlByHash(string) (string, error)
}
