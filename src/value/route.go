package value

type MatcherOption struct {
	Source string `json:"source" binding:"required"` // 匹配源
	Name   string `json:"name" binding:"required"`   // 匹配值
	Type   string `json:"type" binding:"required"`   // 匹配方式
	Value  string `json:"value" binding:"required"`  // 匹配内容
}

type ForwardHeaderOption struct {
	Type  string `json:"type" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Value string `json:"value"`
}

type ForwardRewriteOption struct {
	// 转发正则
	Regex string `json:"regex" binding:"required"`
	// 转发配置
	Replace string `json:"replace" binding:"required"`
}

type ForwardEndpoint struct {
	Static *ForwardEndpointStatic `json:"static"`
}

type ForwardEndpointStatic struct {
	Address []string `json:"address" binding:"required"`
}
