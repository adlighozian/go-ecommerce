package service

import "voucher-go/model"

type Servicer interface {
	GetVoucher() (model.Respon, error)
	ShowVoucher(code string) (model.Respon, error)
	CreateVoucher(req []model.VoucherReq) (model.Respon, error)
	DeleteVoucher(idVoucher int) (model.Respon, error)
}
