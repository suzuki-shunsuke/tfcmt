package domain

type ComplementEntry interface {
	Entry() (string, error)
	Type() string
}
