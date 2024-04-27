package valueobject

type MatcherConfig struct {
	Source string `json:"source" binding:"required"` // 匹配源
	Name   string `json:"name" binding:"required"`   // 匹配值
	Type   string `json:"type" binding:"required"`   // 匹配方式
	Value  string `json:"value" binding:"required"`  // 匹配内容
}
