package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type domain struct {
	Domain     string `json:"domain"`
	SubDomain  string `json:"sub"`
	RecordType string `json:"type"`
}

type Config struct {
	Type      string   `json:"type"`
	AccessKey string   `json:"key"`
	SecretKey string   `json:"secret"`
	Domains   []domain `json:"domains"`
}

// 解析配置文件
func GetConfig(filePath string) ([]Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("config not exist")
	}

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result []Config
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
