package domain

type UrlRedisRepo interface {
	Get(hash string) (string, error)
}

type UrlPostresRepo interface {
	Get(hash string) (string, error)
}

type UrlService interface {
	Get(hash string) (string, error)
}
