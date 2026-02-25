package httpclient

import (
	"encoding/json"
	"net/http"
)


type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}


func (r *Response) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}


func (r *Response) Decode(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func (r *Response) String() string {
	return string(r.Body)
}


func (r *Response) Status() int {
	return r.StatusCode
}
