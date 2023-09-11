package model

import "time"

type DataModel struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	Timestamp time.Time
	Data      string
}
