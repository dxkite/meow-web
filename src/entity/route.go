package entity

import "dxkite.cn/meownest/src/value"

type Route struct {
	Base
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Method       []string             `json:"method" gorm:"serializer:json"`
	Path         string               `json:"path"`
	MatchOptions []*value.MatchOption `json:"match_options" gorm:"serializer:json"`
}
