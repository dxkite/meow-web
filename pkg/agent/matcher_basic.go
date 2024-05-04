package agent

import (
	"fmt"
	"net/http"
	"strings"
)

var _ RequestPathMatcher = (*BasicMatcher)(nil)

type BasicMatcher struct {
	Host   []string
	Method []string
	Path   RequestPathMatcher
	Extra  []*ExtraMatchOption
}

func NewBasicMatcher() *BasicMatcher {
	return &BasicMatcher{}
}

type ExtraMatchOption struct {
	Type   string
	Source string
	Name   string
	Value  string
}

func (b *BasicMatcher) MatchPathType() PathType {
	return b.Path.MatchPathType()
}

func (b *BasicMatcher) MatchPathPriority() int {
	return b.Path.MatchPathPriority()
}

func (b *BasicMatcher) MatchRequest(req *http.Request) bool {
	if len(b.Host) > 0 && !InStringSlice(req.Host, b.Host) {
		return false
	}

	if len(b.Method) > 0 && !InStringSlice(req.Method, b.Method) && !InStringSlice("Any", b.Method) {
		return false
	}

	if !b.Path.MatchRequest(req) {
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

func (b *BasicMatcher) String() string {
	return fmt.Sprintf("%v %v", b.Method, b.Path)
}

type PathType int

const (
	PathTypeNone PathType = iota
	PathTypeFull
	PathTypePrefix
	PathTypeParam
)

type pathMatcher struct {
	path string
	typ  PathType
}

func (m *pathMatcher) MatchPathType() PathType {
	return m.typ
}

func (m *pathMatcher) MatchPathPriority() int {
	switch m.typ {
	case PathTypeFull, PathTypePrefix:
		return len(m.path)
	}
	return 0
}

func (m *pathMatcher) MatchRequest(req *http.Request) bool {
	path := req.URL.Path
	switch m.typ {
	case PathTypeFull:
		return path == m.path
	case PathTypeParam:
		ok, _, _ := TestPath(m.path, path)
		return ok
	case PathTypePrefix:
		fallthrough
	default:
		return strings.HasPrefix(path, m.path)
	}
}

func NewRequestPathMatcher(path string) RequestPathMatcher {
	typ := PathTypePrefix
	if strings.IndexByte(path, '{') >= 0 {
		typ = PathTypeParam
	}
	return NewRequestPathMatcherWithType(typ, path)
}

func NewRequestPathMatcherWithType(typ PathType, path string) RequestPathMatcher {
	return &pathMatcher{path: path, typ: typ}
}
