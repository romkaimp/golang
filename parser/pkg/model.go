package model

import (
    _ "gorm.io/driver/postgres"
    _ "gorm.io/gorm"
    "time"
)

type Product struct {
    ID          uint           `gorm:"primaryKey"`
    Name        string         `gorm:"type:varchar(255);not null"`
    Prices      []ProductPrice `gorm:"foreignKey:ProductID"`
}

type ProductPrice struct {
    ID        uint      `gorm:"primaryKey"`
    ProductID uint      `gorm:"not null"`
    Price     float64   `gorm:"type:decimal(10,2);not null"`
    ValidFrom time.Time `gorm:"type:timestamptz;not null;default:now()"`
}