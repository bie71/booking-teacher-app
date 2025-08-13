package model

import "time"

type PaymentMethod struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;type:varchar(100);not null" json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `gorm:"default:false" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponsePaginate struct {
	Data       []PaymentMethod `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	Total       int `json:"total"`
	Limit       int `json:"limit"`
}
