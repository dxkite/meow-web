package mux

import (
	"errors"
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
		case j >= pnN:
			if pattern != "/" && pnN > 0 && pattern[pnN-1] == '/' {
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
