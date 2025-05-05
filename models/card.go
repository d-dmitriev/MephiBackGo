package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `json:"-"`
	AccountID  uint   `json:"-"`
	Number     string `json:"number"`      // зашифрованный номер карты
	ExpireDate string `json:"expire_date"` // формат: MM/YY
	CVV        string `json:"cvv"`         // хешированный или зашифрованный
	IssuedAt   string `json:"issued_at"`   // дата выпуска
	Status     string `json:"status"`      // active, blocked, expired
}
