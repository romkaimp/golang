package model

import (
	"time"

	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

/*type Product struct {
	ID         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"uniqueIndex;type:varchar(255);not null"`
	ImageBig   string         `gorm:"size:256"`
	ImageSmall string         `gorm:"size:256"`
	Ref        string         `gorm:"size:256"`
	Prices     []ProductPrice `gorm:"foreignKey:ProductID"`
}

type ProductPrice struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	ProductID uint      `gorm:"not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	ValidFrom time.Time `gorm:"type:timestamptz;not null;default:now()"`
}*/

type Product struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"index:name,priority:1;type:varchar(255);not null"`
	ImageBig   string    `gorm:"size:256"`
	ImageSmall string    `gorm:"size:256"`
	Ref        string    `gorm:"size:256"`
	Price      float64   //`gorm:"foreignKey:ProductID"`
	CreatedAt  time.Time `gorm:"index:time,priority:2"`
	UpdatedAt  time.Time
}
