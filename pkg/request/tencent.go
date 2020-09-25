package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/baseli/ddns/internal/crypto"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type tencent struct {
	AccessKey string
	SecretKey string

	req Request
}

type recordItem struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	RecordType string `json:"type"`
	Value      string `json:"value"`
}

type recordListData struct {
	Domain []recordItem `json:"records"`
}

type recordListResult struct {
	Code    uint64         `json:"code"`
	Message string         `json:"message,omitempty"`
	Data    recordListData `json:"data,omitempty"`
}

const TENCENT_URI = "https://cns.api.qcloud.com/v2/index.php"

func newTencent(accessKey string, secretKey string, req Request) tencent {
	return tencent{
		AccessKey: accessKey,
		SecretKey: secretKey,
		req:       req,
	}
}

// 腾讯云接口鉴权处理
func (tencent tencent) auth(params map[string]string) string {
	// 加入时间戳和随机数
	params["Timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["Nonce"] = fmt.Sprint(rand.Uint32())
	params["SecretId"] = tencent.AccessKey
	params["SignatureMethod"] = "HmacSHA256"

	// 排序
	keys := make([]string, len(params))
	index := 0
	for k := range params {
		keys[index] = k
		index++
	}
	sort.Strings(keys)

	// 拼接字符串
	signatureStr := ""
	for _, key := range keys {
		signatureStr += key + "=" + params[key] + "&"
	}
	signatureStr = signatureStr[:len(signatureStr)-1]
	signature := "POSTcns.api.qcloud.com/v2/index.php?" + signatureStr

	// 生成签名
	return signatureStr + "&Signature=" + crypto.Base64Encode(crypto.Sha256(tencent.SecretKey, signature))
}

// 获取历史记录，从而获取到recordId
func (tencent tencent) getRecordId(domain string, recordType string, subDomain string) (uint64, error) {
	params := tencent.auth(map[string]string{
		"Action":     "RecordList",
		"domain":     domain,
		"subDomain":  subDomain,
		"recordType": recordType,
	})

	res, err := tencent.req.GeneralPost(params, TENCENT_URI)
	if err != nil {
		return 0, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	recordList := &recordListResult{}
	err = json.Unmarshal(body, recordList)
	if err != nil {
		return 0, err
	}

	for _, item := range recordList.Data.Domain {
		if item.Name == subDomain && item.RecordType == recordType {
			return item.Id, nil
		}
	}
	return 0, nil
}

func (tencent tencent) Update(domain string, recordType string, subDomain string, ipAddress string) error {
	recordId, err := tencent.getRecordId(domain, recordType, subDomain)
	if err != nil {
		return err
	}

	params := map[string]string{
		"domain":     domain,
		"subDomain":  subDomain,
		"recordId":   strconv.FormatUint(recordId, 10),
		"recordType": recordType,
		"recordLine": "默认",
		"value":      ipAddress,
	}

	params["Action"] = "RecordModify"
	if recordId == 0 {
		// 新增
		params["Action"] = "RecordCreate"
	}

	time.Sleep(time.Second)
	res, err := tencent.req.GeneralPost(tencent.auth(params), TENCENT_URI)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(res.Body)
	recordList := &recordListResult{}
	err = json.Unmarshal(body, recordList)
	if err != nil {
		return err
	}

	if recordList.Code == 0 {
		return nil
	}

	return errors.New(recordList.Message)
}
