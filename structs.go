package suda

type RouteInfo struct {
	Name      string        `yaml:"name"`
	Auth      bool          `yaml:"auth"`
	Rewrite   RewriteConfig `yaml:"rewrite"`
	EndPoints []string
}

func (r *RouteInfo) RouteName() string {
	return r.Name
}

type Token struct {
	ExpireAt int64  `json:"exp"`
	Value    string `json:"val"`
	Scope    string `json:"sco"`
}
