package restclient

import "github.com/go-resty/resty/v2"

var _ Client = &ClientMock{}

type ClientMock struct {
	BaseURL string
	Header  string
	Value   string

	HaveRequest RequestMock
}

func NewClientMock(
	haveResponses map[string]*resty.Response,
	haveErr error,
) *ClientMock {
	return &ClientMock{
		HaveRequest: RequestMock{
			Responses: haveResponses,
			Err:       haveErr,
		},
	}
}

func (a *ClientMock) SetBaseURL(url string) Client {
	a.BaseURL = url
	return a
}

func (a *ClientMock) GetBaseURL() string {
	return a.BaseURL
}

func (a *ClientMock) SetHeader(header, value string) Client {
	a.Header = header
	a.Value = value
	return a
}

func (a *ClientMock) Request() Request {
	return &a.HaveRequest
}
