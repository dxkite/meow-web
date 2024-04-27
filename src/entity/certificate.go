package entity

import "time"

type Certificate struct {
	Base

	Name        string    `json:"name"`
	Domain      []string  `json:"domain" gorm:"serializer:json"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Key         string    `json:"key"`
	Certificate string    `json:"certificate"`
}
