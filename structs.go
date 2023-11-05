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
