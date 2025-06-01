package restclient

import "github.com/go-resty/resty/v2"

type Client interface {
	GetBaseURL() string
	SetBaseURL(url string) Client
	SetHeader(header, value string) Client

	Request() Request
}

type clientAdapter struct {
	client *resty.Client
}

func newClientAdapter() *clientAdapter {
	return &clientAdapter{
		client: resty.New(),
	}
}

func (a *clientAdapter) SetBaseURL(url string) Client {

	a.client.SetBaseURL(url)
	return a
}

func (a *clientAdapter) GetBaseURL() string {
	return a.client.BaseURL
}

func (a *clientAdapter) SetHeader(header, value string) Client {

	a.client.SetHeader(header, value)
	return a
}

func (a *clientAdapter) Request() Request {
	return &requestAdapter{
		request: a.client.R(),
	}
}
