package model

type Route struct {
	Base
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Method      []string         `json:"method" gorm:"serializer:json"`
	Path        string           `json:"path"`
	Matcher     []*MatcherConfig `json:"matcher" gorm:"serializer:json"`
}

type MatcherConfig struct {
	Source string `json:"source"` // 匹配源
	Name   string `json:"name"`   // 匹配值
	Type   string `json:"type"`   // 匹配方式
	Value  string `json:"value"`  // 匹配内容
}
