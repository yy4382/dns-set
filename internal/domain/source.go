package domain

type DomainSource interface {
	GetDomains() ([]string, error)
	Name() string
}
