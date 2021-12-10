package ssclient

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
)

type methodType int

const (
	methodGet methodType = iota
	methodPost
	methodPut
	methodDelete
	methodPatch
	methodHead
	methodOptions
)

const userAgent = "terraform"

func (m methodType) String() string {
	switch m {
	case methodGet:
		return "GET"
	case methodPost:
		return "POST"
	case methodPut:
		return "PUT"
	case methodDelete:
		return "DELETE"
	case methodPatch:
		return "PATCH"
	case methodHead:
		return "HEAD"
	case methodOptions:
		return "OPTIONS"
	default:
		panic("WRONG METHOD ID")
	}
}

func makeRequest(
	client *resty.Client,
	url string,
	method methodType,
	payload interface{},
	res interface{},
) (interface{}, error) {
	request := client.R().
		SetError(&ErrorBodyResponse{}).
		SetHeaders(map[string]string{
			"User-Agent": userAgent,
		})

	if res != nil {
		request = request.SetResult(res)
	}
	if payload != nil {
		request = request.SetBody(payload)
	}
	var (
		resp *resty.Response
		err  error
	)
	debugReqest, err := json.MarshalIndent(struct {
		Method string      `json:"method,omitempty"`
		URL    string      `json:"url,omitempty"`
		Body   interface{} `json:"body,omitempty"`
	}{
		Method: method.String(),
		URL:    url,
		Body:   payload,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Default().Println("[DEBUG]  Request: ", string(debugReqest))

	switch method {
	case methodGet:
		resp, err = request.Get(url)
	case methodPost:
		resp, err = request.Post(url)
	case methodPut:
		resp, err = request.Put(url)
	case methodDelete:
		resp, err = request.Delete(url)
	case methodPatch:
		resp, err = request.Patch(url)
	case methodHead:
		resp, err = request.Head(url)
	case methodOptions:
		resp, err = request.Options(url)
	default:
		return nil, errors.New("wrong method type")
	}

	if err != nil {
		return nil, NewRequestError(resp, err)
	}

	respBody := resp.Result()
	debugReqest, err = json.MarshalIndent(struct {
		Method     string      `json:"method"`
		URL        string      `json:"url"`
		Body       interface{} `json:"body"`
		Statuscode int         `json:"status_code"`
		Status     string      `json:"status"`
		Response   interface{} `json:"response"`
	}{
		Method:     method.String(),
		URL:        resp.Request.URL,
		Body:       payload,
		Statuscode: resp.StatusCode(),
		Status:     resp.Status(),
		Response:   respBody,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Default().Println("[DEBUG]  Performed request ", string(debugReqest))

	if err != nil {
		return nil, NewRequestError(resp, err)
	}

	if resp.IsError() {
		return nil, NewRequestError(resp, nil)
	}

	return respBody, nil
}
