package agent

import "net/http"

type BasicMatcher struct {
	Host   []string
	Method []string
	Path   string
	Extra  []*ExtraMatchOption
}

type ExtraMatchOption struct {
	Type   string
	Source string
	Name   string
	Value  string
}

func (b *BasicMatcher) MatchRequest(req *http.Request) bool {
	if !InStringSlice(req.Host, b.Host) {
		return false
	}

	if !InStringSlice(req.Method, b.Method) && !InStringSlice("Any", b.Method) {
		return false
	}

	if ok, _, _ := TestPath(b.Path, req.URL.Path); !ok {
		return false
	}

	for _, e := range b.Extra {
		value := VarFrom(req, e.Source, e.Name)
		switch e.Type {
		case "equal", "=", "":
			if value != e.Value {
				return false
			}
		}
	}

	return true
}
