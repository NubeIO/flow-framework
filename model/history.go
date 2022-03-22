package model

import "time"

//History for storing the all history
type History struct {
	ID        int       `json:"id" gorm:"primary_key"`
	UUID      string    `json:"uuid" gorm:"primary_key"`
	Value     float64   `json:"value" gorm:"primary_key"`
	Timestamp time.Time `json:"timestamp" gorm:"primary_key"`
}
