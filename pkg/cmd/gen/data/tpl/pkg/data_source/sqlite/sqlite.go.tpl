package sqlite

import (
	"{{ .Pkg }}/pkg/data_source"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DataSource struct {
	db *gorm.DB
}

func NewSQLiteDataSource(db *gorm.DB) *DataSource {
	return &DataSource{db: db}
}

func Open(dsn string) (*DataSource, error) {
	db, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, err
	}
	db = db.Debug()
	return NewSQLiteDataSource(db), nil
}

func (s *DataSource) RawSource() interface{} {
	return s.db
}

func (s *DataSource) Transaction(fn func(s data_source.DataSource) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewSQLiteDataSource(tx))
	})
}
