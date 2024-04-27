package model

type Link struct {
	Base
	Direct   string `gorm:"index"`
	SourceId uint64
	LinkedId uint64
}
