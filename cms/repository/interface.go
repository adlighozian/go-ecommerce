package repositoy

type Repositorier interface {
	Get(..., err error)
	GetDetail(..., err error)
	Create(..., err error)
	Update(..., err error)
	Delete(..., err error)
}