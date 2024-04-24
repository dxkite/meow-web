package datasource

import "gorm.io/gorm"

type DataSource struct {
	Raw *gorm.DB
}

func New() *DataSource {
	return &DataSource{}
}
