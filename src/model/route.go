package model

import "gorm.io/gorm"

type Route struct {
	gorm.Model
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Method      []string         `json:"method"`
	Path        string           `json:"path"`
	Matcher     []*MatcherConfig `json:"matcher"`
}

type MatcherConfig struct {
	Source string `json:"source"` // 匹配源
	Name   string `json:"name"`   // 匹配值
	Type   string `json:"type"`   // 匹配方式
	Value  string `json:"value"`  // 匹配内容
}
