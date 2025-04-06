package domain

type ClickhouseDB interface {
	Create(data RegisterRedirectEventDTO) error
	Get() ([]RedirectEventDAO, error)
}

type Kafka interface {
	Create(data RegisterRedirectEventDTO) error
}
