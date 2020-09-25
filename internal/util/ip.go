package util

import (
	"github.com/baseli/ddns/pkg/request"
	"io/ioutil"
	"strings"
)

// 获取ipv6地址
func GetIpv6(req request.Request) (string, error) {
	return extractorIp("http://ip6only.me/api/", req)
}

// 获取ipv4地址
func GetIpv4(req request.Request) (string, error) {
	return extractorIp("http://ip4only.me/api/", req)
}

func extractorIp(url string, req request.Request) (string, error) {
	res, err := req.Get(url)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(res.Body)
	return strings.Split(string(body), ",")[1], nil
}
