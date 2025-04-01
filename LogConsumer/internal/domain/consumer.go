package domain

type ClickhouseDB interface {
	Create(data RegisterRedirectEvent) error
	Get() ([]RedirectEventDAO, error)
}
