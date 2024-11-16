package agent

import (
	"net/url"
	"reflect"
	"testing"
)

func TestMatchPath(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		match   bool
		value   url.Values
		err     error
	}{
		{"/", "/", true, url.Values{}, nil},
		{"/{name}/foo", "/foo/foo", true, url.Values{"name": []string{"foo"}}, nil},
		{"/{name}/foo/", "/foo/foo", false, nil, nil},
		{"/{name}/foo/", "/foo/foo/", true, url.Values{"name": []string{"foo"}}, nil},
		{"/{name}/foo/", "/foo/foo/abc", true, url.Values{"name": []string{"foo"}, "$": []string{"abc"}}, nil},
		{"/{name}.txt", "/foo.txt", true, url.Values{"name": []string{"foo"}}, nil},
		{"/{name}.{ext}", "/foo.txt", true, url.Values{"name": []string{"foo"}, "ext": []string{"txt"}}, nil},
		{"/{name}/{value}/{test}/test", "/foo1/foo2/foo3/test", true, url.Values{"name": []string{"foo1"}, "value": []string{"foo2"}, "test": []string{"foo3"}}, nil},
		{"/{name}/{value}/{test}", "/foo1/foo2/foo3", true, url.Values{"name": []string{"foo1"}, "value": []string{"foo2"}, "test": []string{"foo3"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			ok, vars, err := TestPath(tt.pattern, tt.path)
			if !reflect.DeepEqual(vars, tt.value) {
				t.Errorf("MatchPath() got = %v, want %v", ok, tt.value)
			}
			if ok != tt.match {
				t.Errorf("MatchPath() got1 = %v, want %v", vars, tt.match)
			}
			if err != tt.err {
				t.Errorf("MatchPath() err = %v, want %v", vars, tt.match)
			}
		})
	}
}
