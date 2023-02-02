package ssclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

var HOST_MAP = map[string]string{
	"02": "https://api.serverspace.by",
	"04": "https://api.serverspace.io",
	"06": "https://api.serverspace.ru",
	"07": "https://api.lincore.kz",
	"08": "https://api.serverspace.us",
	"09": "https://api.serverspace.com.tr",
	"0a": "https://api.serverspace.in",
	"0f": "https://api.itglobal.com",
	"14": "https://api.serverspace.kz",
}

type SSClient struct {
	client *resty.Client
	Key    string
	Host   string
}

func NewClient(key string, host string) (*SSClient, error) {
	if host == "" {
		if len(key) < 2 {
			return nil, NewWrongKeyFormatError(nil)
		}

		var ok bool
		hostIndex := key[:2]
		if host, ok = HOST_MAP[hostIndex]; !ok {
			return nil, NewWrongKeyFormatError(nil)
		}
	}

	client := resty.New()
	client.SetHeader("X-API-KEY", key)

	baseURL := fmt.Sprintf("%s/%s", host, "api/v1/")
	client.SetHostURL(baseURL)

	return &SSClient{client, key, host}, nil
}
