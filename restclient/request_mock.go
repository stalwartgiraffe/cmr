package restclient

import (
	"context"

	"github.com/go-resty/resty/v2"
)

var _ Request = &RequestMock{}

type RequestMock struct {
	Context   context.Context
	AuthToken string
	Query     string
	Headers   []string
	Values    []string
	Body      any

	Responses map[string]*resty.Response
	Err       error
}

func (a *RequestMock) SetContext(ctx context.Context) Request {
	a.Context = ctx
	return a
}
func (a *RequestMock) SetAuthToken(token string) Request {
	a.AuthToken = token
	return a
}
func (a *RequestMock) SetQueryString(query string) Request {
	a.Query = query
	return a
}

func (a *RequestMock) SetHeader(header, value string) Request {
	a.Headers = append(a.Headers, header)
	a.Values = append(a.Values, value)
	return a
}

func (a *RequestMock) SetBody(body any) Request {
	a.Body = body
	return a
}

func (a *RequestMock) fetch(url string) (*resty.Response, error) {
	r, ok := a.Responses[url]
	if !ok {
		return nil, a.Err
	}
	return r, nil
}

func (a *RequestMock) Get(url string) (*resty.Response, error) {
	return a.fetch(url)
}
func (a *RequestMock) Head(url string) (*resty.Response, error) {
	panic("not implemented")
}
func (a *RequestMock) Post(url string) (*resty.Response, error) {
	panic("not implemented")
}
func (a *RequestMock) Put(url string) (*resty.Response, error) {
	panic("not implemented")
}
func (a *RequestMock) Delete(url string) (*resty.Response, error) {
	panic("not implemented")
}
func (a *RequestMock) Options(url string) (*resty.Response, error) {
	panic("not implemented")
}
func (a *RequestMock) Patch(url string) (*resty.Response, error) {
	panic("not implemented")
}
