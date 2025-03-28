package domain

type UrlRedisRepo interface {
	Create(url string, hash string) error
}

type UrlPostresRepo interface {
	Create(url string, hash string) error
}

type UrlService interface {
	Create(url string) (string, error)
}
