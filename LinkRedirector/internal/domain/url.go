package domain

type UrlRedisRepo interface {
	Get(hash string) (string, error)
}

type UrlPostresRepo interface {
	Get(hash string) (string, error)
}

type UrlService interface {
	Get(hash GetUrlDTO) (string, error)
}

type UrlKafka interface {
	Send(message UrlKafkaDTO) error
}
