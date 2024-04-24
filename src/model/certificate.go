package model

import "gorm.io/gorm"

type Certificate struct {
	gorm.Model

	Domain      []string `json:"domain"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	Certificate string   `json:"certificate"`

	ExpireAt uint64 `json:"expire_at"`
}
