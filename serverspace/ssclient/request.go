package ssclient

import (
	"errors"

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

func makeRequest(
	client *resty.Client,
	url string,
	method methodType,
	payload interface{},
	res interface{},
) (interface{}, error) {
	request := client.R().
		SetResult(res). //SetError(map[string]interface{}{})
		SetError(&ErrorBodyResponse{})
	if payload != nil {
		request = request.SetBody(payload)
	}
	var (
		resp *resty.Response
		err  error
	)
	switch method {
	case methodGet:
		resp, err = request.Get(url)
		// if strings.Contains(url, "servers") {
		// 	tmpReq := client.R().
		// 		SetResult(map[string]interface{}{}). //SetError(map[string]interface{}{})
		// 		SetError(&ErrorBodyResponse{})
		// 	if payload != nil {
		// 		request = request.SetBody(payload)
		// 	}
		// 	tmpResp, _ := tmpReq.Get(url)
		// 	log.Default().Println("Result %ы 00000000000: ::", tmpReq.URL, tmpResp.Result())
		// }
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
		return nil, errors.New("Wrong method type")
	}

	// marshaledBody, err := json.MarshalIndent(resp.Request.Body, "", "\t")
	// log.Default().Printf("%s, %s, %w", marshaledBody, err)

	if err != nil {
		return nil, NewRequestError(resp, err)
	}

	if resp.IsError() {
		return nil, NewRequestError(resp, nil)
	}

	return resp.Result(), nil
}
