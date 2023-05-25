package  service

type Servicer interface {
	Get(..., err error)
	GetDetail(..., err error)
	Create(..., err error)
	Update(..., err error)
	Delete(..., err error)
}