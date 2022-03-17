package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type params map[string]interface{}

type Request struct {
	Method     string
	Endpoint   string
	Query      url.Values
	Form       url.Values
	recvWindow int64
	Header     http.Header
	Body       io.Reader
	FullURL    string
}

// setParam set param with key/value to query string
func (r *Request) SetParam(key string, value interface{}) *Request {
	if r.Query == nil {
		r.Query = url.Values{}
	}
	r.Query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setParams set params with key/values to query string
func (r *Request) SetParams(m params) *Request {
	for k, v := range m {
		r.SetParam(k, v)
	}
	return r
}
