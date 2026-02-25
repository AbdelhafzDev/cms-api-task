package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Request struct {
	client  *Client
	method  string
	path    string
	headers http.Header
	query   map[string]string
	body    io.Reader
}


func (r *Request) Header(key, value string) *Request {
	r.headers.Set(key, value)
	return r
}

func (r *Request) Query(key, value string) *Request {
	r.query[key] = value
	return r
}


func (r *Request) JSON(body interface{}) *Request {
	if data, err := json.Marshal(body); err == nil {
		r.body = bytes.NewReader(data)
		r.headers.Set("Content-Type", "application/json")
	}
	return r
}


func (r *Request) Text(body string) *Request {
	r.body = bytes.NewReader([]byte(body))
	r.headers.Set("Content-Type", "text/plain")
	return r
}

func (r *Request) Raw(body interface{}) *Request {
	switch v := body.(type) {
	case []byte:
		r.body = bytes.NewReader(v)
	case string:
		r.body = bytes.NewReader([]byte(v))
	case io.Reader:
		r.body = v
	}
	return r
}


func (r *Request) Do(ctx context.Context) (*Response, error) {
	url := r.buildURL()

	req, err := http.NewRequestWithContext(ctx, r.method, url, r.body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	r.applyHeaders(req)
	r.applyQuery(req)

	resp, err := r.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}

func (r *Request) Into(ctx context.Context, target interface{}) error {
	resp, err := r.Do(ctx)
	if err != nil {
		return err
	}
	if !resp.OK() {
		return fmt.Errorf("status %d: %s", resp.StatusCode, resp.String())
	}
	return resp.Decode(target)
}

func (r *Request) DoStream(ctx context.Context) (*http.Response, error) {
	url := r.buildURL()

	req, err := http.NewRequestWithContext(ctx, r.method, url, r.body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	r.applyHeaders(req)
	r.applyQuery(req)

	resp, err := r.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (r *Request) buildURL() string {
	if r.client.baseURL != "" && !isAbsoluteURL(r.path) {
		return r.client.baseURL + r.path
	}
	return r.path
}

func (r *Request) applyHeaders(req *http.Request) {
	for k, v := range r.client.headers {
		req.Header[k] = v
	}
	for k, v := range r.headers {
		req.Header[k] = v
	}
}

func (r *Request) applyQuery(req *http.Request) {
	if len(r.query) == 0 {
		return
	}
	q := req.URL.Query()
	for k, v := range r.query {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func isAbsoluteURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}
