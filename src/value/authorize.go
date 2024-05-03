package value

type AuthorizeAttribute struct {
	Binary *AuthorizeAttributeBinary `json:"binary,omitempty"`
}

type AuthorizeAttributeBinary struct {
	Key     string             `json:"key"`
	Header  string             `json:"header"`
	Sources []*AuthorizeSource `json:"sources" binding:"required,dive,required"`
}

type AuthorizeSource struct {
	Source string `json:"source" binding:"required"` // 匹配源
	Name   string `json:"name" binding:"required"`   // 匹配值
}
