package model

import "dxkite.cn/meownest/src/valueobject"

type Route struct {
	Base
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Method      []string                     `json:"method" gorm:"serializer:json"`
	Path        string                       `json:"path"`
	Matcher     []*valueobject.MatcherConfig `json:"matcher" gorm:"serializer:json"`
}
