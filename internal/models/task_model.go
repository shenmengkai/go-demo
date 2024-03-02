package models

type Task struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Text   string `json:"text"`
	Status int    `json:"status"`
}
