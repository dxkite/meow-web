package route

import (
	"dxkite.cn/gateway/config"
	"fmt"
	"testing"
)

func TestRoute_Match(t *testing.T) {
	cfg := config.Config{}
	if err := cfg.LoadFromFile("./tests/config.yml"); err != nil {
		t.Error("load config error", err)
	}
	tests := []struct {
		path string
		want string
	}{
		{"/user/signin", "/user/signin"},
		{"/user/signout", "/user/signout"},
		{"/user/info", "/user"},
		{"/swiper/selectFront", "/"},
		{"/demo/get", "/demo"},
	}
	r := NewRoute()
	r.Load(cfg.Routes)
	r.Load([]config.Route{
		{
			Pattern: "/demo",
			Backend: []string{"http://127.0.0.1:8888"},
		},
	})
	fmt.Println(r.re)
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, _ := r.Match(tt.path)
			if got != tt.want {
				t.Errorf("Match() got = %v, want %v", got, tt.want)
			}
		})
	}
}
