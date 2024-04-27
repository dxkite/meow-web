package gateway

import (
	"fmt"
	"net/http"
	"net/url"
)

type Metadata map[string]string

func (md *Metadata) FromHeaders(header http.Header) *Metadata {
	for k, v := range header {
		(*md)[k] = v[0]
	}
	return md
}

func (md *Metadata) FromQuery(query url.Values) *Metadata {
	for k, v := range query {
		(*md)[k] = v[0]
	}
	return md
}

func (md *Metadata) FromString(str string) *Metadata {
	val, _ := url.ParseQuery(str)
	md.FromQuery(val)
	return md
}

func (md *Metadata) FromCookies(query []*http.Cookie) *Metadata {
	for _, v := range query {
		(*md)[v.Name] = v.Value
	}
	return md
}

func (md *Metadata) Equal(t *Metadata) bool {
	for k, v := range *t {
		if (*t)[k] != v {
			fmt.Println(k, v, "not equal")
			return false
		}
	}
	return true
}
