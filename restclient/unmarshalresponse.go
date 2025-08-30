package restclient

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"

	"github.com/stalwartgiraffe/cmr/withstack"
)

// Unmarshal If response is success, will unmarshal json from a response body else return a failure error.
func Unmarshal[T any](resp *resty.Response) (*T, error) {
	if resp == nil {
		return nil, NewFailureResponse("Response object was nil", resp)
	}
	if resp.IsError() {
		const includeBody = true
		return nil, NewFailureResponse(SprintResponse(resp), resp)
	}
	if !resp.IsSuccess() {
		return nil, NewFailureResponse("Response object had failure status", resp)
	}

	var body T
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, withstack.Errorf("Unmarshal error:%w", err)
	}
	return &body, nil
}

// Define a custom error type
type FailureResponse struct {
	Msg     string         `json:"msg"`
	Status  string         `json:"status"`
	Request map[string]any `json:"request"`
}

// NewFailureResponse returns a custom response error.
func NewFailureResponse(msg string, resp *resty.Response) *FailureResponse {
	if resp == nil {
		// The client did not error but gave a nil response.
		// This can happen in test but we hope is rare in production.
		// Wrap state with a meaningful error indicating upstream service was not available.
		return &FailureResponse{
			Msg:    withstack.StackTrace() + msg,
			Status: "502 Bad Gateway",
		}
	}
	if resp.Request == nil {
		return &FailureResponse{
			Msg:    withstack.StackTrace() + msg,
			Status: "502 Bad Gateway, missing request",
		}
	}
	return &FailureResponse{
		Msg:    withstack.StackTrace() + msg,
		Status: resp.Status(),
		Request: map[string]any{
			"method": resp.Request.Method,
			"url":    resp.Request.URL,
			"query":  resp.Request.QueryParam,
			"form":   resp.Request.FormData,
			//			"body":   resp.Request.Body,
		},
	}
}

// Implement the Error() method for the custom error type
func (e *FailureResponse) Error() string {
	if b, err := yaml.Marshal(*e); err != nil {
		return "yaml.Marshal error: " + err.Error()
	} else {
		return string(b)
	}
}
