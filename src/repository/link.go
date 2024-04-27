package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Link interface {
	Link(ctx context.Context, direct string, sourceId, linkedId uint64) error
	LinkOnce(ctx context.Context, direct string, sourceId, linkedId uint64) error
	LinkOf(ctx context.Context, direct string, sourceId uint64) ([]*entity.Link, error)
}

func NewLink(db *gorm.DB) Link {
	return &link{db: db}
}

type link struct {
	db *gorm.DB
}

func (r *link) Link(ctx context.Context, direct string, sourceId, linkedId uint64) error {
	link := entity.Link{}
	link.Direct = direct
	link.SourceId = sourceId
	link.LinkedId = linkedId
	return r.db.Create(&direct).Error
}

func (r *link) LinkOnce(ctx context.Context, direct string, sourceId, linkedId uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(entity.Link{Direct: direct, SourceId: sourceId}).Error; err != nil {
			return err
		}
		link := entity.Link{}
		link.Direct = direct
		link.SourceId = sourceId
		link.LinkedId = linkedId
		return r.db.Create(&direct).Error
	})
}

func (r *link) LinkOf(ctx context.Context, direct string, sourceId uint64) ([]*entity.Link, error) {
	links := []*entity.Link{}
	if err := r.db.Where(entity.Link{SourceId: sourceId}).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}
