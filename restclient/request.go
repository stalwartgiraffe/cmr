package restclient

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type Request interface {
	SetContext(ctx context.Context) Request
	SetAuthToken(token string) Request
	SetQueryString(query string) Request

	SetHeader(header, value string) Request
	SetBody(body any) Request

	Get(url string) (*resty.Response, error)

	Head(url string) (*resty.Response, error)
	Post(url string) (*resty.Response, error)
	Put(url string) (*resty.Response, error)
	Delete(url string) (*resty.Response, error)
	Options(url string) (*resty.Response, error)
	Patch(url string) (*resty.Response, error)
}

type requestAdapter struct {
	request *resty.Request
}

func (a *requestAdapter) SetContext(ctx context.Context) Request {
	a.request.SetContext(ctx)
	return a
}

func (a *requestAdapter) SetAuthToken(token string) Request {
	a.request.SetAuthToken(token)
	return a
}
func (a *requestAdapter) SetQueryString(query string) Request {
	a.request.SetQueryString(query)
	return a
}

func (a *requestAdapter) SetHeader(header, value string) Request {
	a.request.SetHeader(header, value)
	return a
}

func (a *requestAdapter) SetBody(body any) Request {
	a.request.SetBody(body)
	return a
}

func (a *requestAdapter) Get(url string) (*resty.Response, error) {
	return a.request.Get(url)
}

func (a *requestAdapter) Head(url string) (*resty.Response, error) {
	return a.request.Head(url)
}
func (a *requestAdapter) Post(url string) (*resty.Response, error) {
	return a.request.Post(url)
}
func (a *requestAdapter) Put(url string) (*resty.Response, error) {
	return a.request.Put(url)
}
func (a *requestAdapter) Delete(url string) (*resty.Response, error) {
	return a.request.Delete(url)
}
func (a *requestAdapter) Options(url string) (*resty.Response, error) {
	return a.request.Options(url)
}
func (a *requestAdapter) Patch(url string) (*resty.Response, error) {
	return a.request.Patch(url)
}
