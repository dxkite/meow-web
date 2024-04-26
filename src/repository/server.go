package repository

import "gorm.io/gorm"

type Server interface {
}

func NewServer(db *gorm.DB) Server {
	return nil
}

type server struct {
	db *gorm.DB
}
