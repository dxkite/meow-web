package suda

type RouteInfo struct {
	Name string
	RouteConfig
}

func (r *RouteInfo) RouteName() string {
	return r.Name
}

type Token struct {
	ExpireAt int64  `json:"exp"`
	Value    string `json:"val"`
}
