package repository

import "voucher-go/model"

type Repositorier interface {
	GetVoucher() ([]model.Voucher, error)
	ShowVoucher(code string) (model.Voucher, error)
	CreateVoucher(req []model.VoucherReq) ([]model.Voucher, error)
	DeleteVoucher(id int) (int, error)
}
