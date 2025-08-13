package model

import "time"

type Payment struct {
	ID                    uint   `gorm:"primaryKey"`
	MidtransTransactionID string `gorm:"size:100"`
	Amount                float64
	Status                string `gorm:"type:enum('pending','settlement','failed','cancel');default:'pending'"`
	PaymentMethod         string
	BookingID             uint `gorm:"index"`
	PaidAt                *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type BookingResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentResponse struct {
	Id                    uint       `json:"id"`
	MidtransTransactionID string     `json:"midtrans_transaction_id"`
	Amount                float64    `json:"amount"`
	Status                string     `json:"status"`
	PaymentMethod         string     `json:"payment_method"`
	BookingID             uint       `json:"booking_id"`
	PaidAt                *time.Time `json:"paid_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
