package ssclient

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type BaseClientError struct {
	Msg string
	Err error
}

func (e *BaseClientError) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Err)
}

func (e *BaseClientError) Unwrap() error {
	return e.Err
}

type WrongKeyFormatError struct {
	BaseClientError
}

func NewWrongKeyFormatError(err error) *WrongKeyFormatError {
	return &WrongKeyFormatError{
		BaseClientError{
			Msg: "Wrong key format",
			Err: err,
		},
	}
}

type ErrorBodyResponse struct {
	Errors []*struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

type RequestError struct {
	BaseClientError
	Response *resty.Response
	Status   int
	Body     string
}

func NewRequestError(response *resty.Response, err error) *RequestError {
	return &RequestError{
		BaseClientError: BaseClientError{
			Msg: "Request isn't ok",
			Err: err,
		},
		Response: response,
		Status:   response.StatusCode(),
		Body:     response.String(),
	}
}
func (e *RequestError) Error() string {
	if e.Err == nil {

		errData, err := json.MarshalIndent(struct {
			Method     string      `json:"method,omitempty"`
			URL        string      `json:"url,omitempty"`
			Body       interface{} `json:"body,omitempty"`
			Statuscode int         `json:"status_code,omitempty"`
			Status     string      `json:"status,omitempty"`
			Response   interface{} `json:"response,omitempty"`
		}{
			Method:     e.Response.Request.Method,
			URL:        e.Response.Request.URL,
			Body:       e.Response.Request.Body,
			Statuscode: e.Response.StatusCode(),
			Status:     e.Response.Status(),
			Response:   e.Response.Error(),
		}, "", "  ")
		if err != nil {
			return fmt.Errorf("Error on marshaling json body on error (%s): %w", e.Msg, err).Error()
		}
		return fmt.Sprintf("%s:\n %s", e.Msg, errData)

	}
	errResp := e.Err.Error()
	return fmt.Sprintf("%s: %s", e.Msg, errResp)
}
