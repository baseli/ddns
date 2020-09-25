package request

import (
	"errors"
	"strings"
)

type DDNsContext struct {
	Request DDnsRequest
}

func NewContext(strategy string, accessKey string, secretKey string, req Request) (*DDNsContext, error) {
	context := new(DDNsContext)

	switch strings.ToUpper(strategy) {
	case "TENCENT":
		context.Request = newTencent(accessKey, secretKey, req)
	default:
		return nil, errors.New("unknown type")
	}

	return context, nil
}
