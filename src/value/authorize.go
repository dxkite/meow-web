package value

type AuthorizeAttribute struct {
	Binary  *AuthorizeAttributeBinary `json:"binary,omitempty"`
	Sources []*AuthorizeSource        `json:"sources" binding:"required,dive,required"`
}

type AuthorizeAttributeBinary struct {
	Key string `json:"key"`
}

type AuthorizeSource struct {
	Source string `json:"source" binding:"required"` // 匹配源
	Name   string `json:"name" binding:"required"`   // 匹配值
}
