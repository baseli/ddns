package request

import (
	"net/http"
	"strings"
	"time"
)

type Request struct {
	client http.Client
}

func NewRequest() Request {
	return Request{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// 正常的post请求
func (req Request) GeneralPost(params string, url string) (*http.Response, error) {
	return req.client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params))
}

// get请求
func (req Request) Get(url string) (*http.Response, error) {
	return req.client.Get(url)
}
