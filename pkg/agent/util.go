package agent

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrInvalidPattern = errors.New("invalid value pattern")

// 路径参数模式匹配
// 支持匹配带参数数据格式 /{name}/xx
func TestPath(pattern, path string) (bool, url.Values, error) {
	p := make(url.Values)
	var i, j int
	pnN := len(pattern)
	phN := len(path)
	for i < phN {
		switch {
		case j >= pnN: // 匹配到末尾
			if pattern != "/" && pnN > 0 && pattern[pnN-1] == '/' {
				extPath := path[i:]
				value, err := url.QueryUnescape(extPath)
				if err != nil {
					p.Add("$", extPath)
					return true, p, nil
				}
				p.Add("$", value)
				return true, p, nil
			}
			return false, nil, nil
		case pattern[j] == '{': // {name}
			nameEnd := j + 1
			for ; nameEnd < pnN && pattern[nameEnd] != '}'; nameEnd++ {
			}

			name := pattern[j+1 : nameEnd]
			if pattern[nameEnd] != '}' {
				return false, nil, ErrInvalidPattern
			}

			matchEnd := i
			for ; matchEnd < phN && path[matchEnd] != '/'; matchEnd++ {
				if nameEnd+1 < pnN && pattern[nameEnd+1] == path[matchEnd] {
					break
				}
			}

			matchValue := path[i:matchEnd]
			value, err := url.QueryUnescape(matchValue)
			if err != nil {
				return false, nil, err
			}

			p.Add(name, value)
			i = matchEnd
			j = nameEnd + 1
		case path[i] == pattern[j]:
			i++
			j++
		default:
			return false, nil, nil
		}
	}

	if j != pnN {
		return false, nil, nil
	}

	return true, p, nil
}

func VarFrom(req *http.Request, source, name string) string {
	switch source {
	case "cookie":
		if c, err := req.Cookie(name); err != http.ErrNoCookie {
			return ""
		} else if c != nil {
			return c.Value
		}
	case "header":
		if v := req.Header.Get(name); v != "" {
			return v
		}
	case "query":
		if v := req.URL.Query().Get(name); v != "" {
			return v
		}
	}
	return ""
}

func InStringSlice(v string, slice []string) bool {
	for _, m := range slice {
		if v == m {
			return true
		}
	}
	return false
}

func printLog(format string, values ...interface{}) {
	fmt.Printf(format, values...)
}
