package repository

import (
	"context"
	"math"
	"payment/internal/model"
	"payment/internal/pkg"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreatePayment(payment *model.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *Repository) GetPaymentByTrxID(trxId string) (*model.Payment, error) {
	var payment model.Payment
	err := r.DB.Where("midtrans_transaction_id = ?", trxId).First(&payment).Error
	return &payment, err
}

func (r *Repository) GetPayments(pg int, limit int) (pkg.ResponsePaginate, error) {
	var payments []model.Payment
	var total int64
	if pg == 0 {
		pg = 1
	}

	if limit == 0 {
		limit = 10
	}

	err := r.DB.Model(&model.Payment{}).Count(&total).Error
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}
	err = r.DB.Limit(limit).Offset((pg - 1) * limit).Find(&payments).Error
	if err != nil {
		return pkg.ResponsePaginate{}, err
	}

	return pkg.ResponsePaginate{
		Data: payments,
		Pagination: pkg.PaginationPage{
			CurrentPage: pg,
			TotalPage:   int(math.Ceil(float64(total) / float64(limit))),
			TotalData:   int(total),
			Limit:       limit,
		},
	}, nil
}

func (r *Repository) CreatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error {
	return r.DB.WithContext(ctx).Create(method).Error
}

func (r *Repository) GetPaymentMethodByID(ctx context.Context, id uint) (*model.PaymentMethod, error) {
	var method model.PaymentMethod
	err := r.DB.WithContext(ctx).First(&method, id).Error
	return &method, err
}

func (r *Repository) GetPaymentMethods(ctx context.Context, page int, limit int) ([]model.PaymentMethod, int64, error) {
	var methods []model.PaymentMethod
	var total int64

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	err := r.DB.Model(&model.PaymentMethod{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.WithContext(ctx).Limit(limit).Offset((page - 1) * limit).Find(&methods).Error
	if err != nil {
		return nil, 0, err
	}

	return methods, total, nil
}

func (r *Repository) UpdatePaymentMethod(ctx context.Context, method *model.PaymentMethod) error {
	return r.DB.WithContext(ctx).Save(method).Error
}

func (r *Repository) DeletePaymentMethod(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&model.PaymentMethod{}, id).Error
}

func (r *Repository) UpdatePayment(payment *model.Payment) error {
	return r.DB.Save(payment).Error
}

func (r *Repository) GetAll(ctx context.Context) ([]model.PaymentMethod, error) {
	var methods []model.PaymentMethod
	err := r.DB.WithContext(ctx).Find(&methods).Error
	return methods, err
}

func (r *Repository) UpdateStatus(ctx context.Context, code string, active bool) error {
	return r.DB.WithContext(ctx).Model(&model.PaymentMethod{}).
		Where("code = ?", code).
		Update("is_active", active).Error
}

func (r *Repository) GetPaymentMethodByCode(ctx context.Context, code string) (*model.PaymentMethod, error) {
	var method model.PaymentMethod
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&method).Error
	return &method, err
}

func (r *Repository) GetPaymentMethodByName(ctx context.Context, code string) (*model.PaymentMethod, error) {
	var method model.PaymentMethod
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&method).Error
	return &method, err
}

func (r *Repository) GetPaymentById(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.DB.Where("id = ?", id).First(&payment).Error
	return &payment, err
}
