package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
    return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }

    return json.Unmarshal(b, &a)
} 



type Products struct {
	UserID             int      `json:"user_id" gorm:"foreignKey:UserID;references:UserID"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      StringArray   `json:"product_images" gorm:"type:json"`
	ProductPrice       float64  `json:"product_price"`
	CompressedImages	 StringArray   `json:"compressed_images" gorm:"type:json"`
}

func MigrateProducts(db *gorm.DB) error {
	err := db.AutoMigrate(&Products{},&Users{})
	return err
}