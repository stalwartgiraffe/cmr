package gitlab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	rc "github.com/stalwartgiraffe/cmr/restclient"
)

type Client struct {
	client *rc.AuthTokenClient
}

func NewClient(overrides ...rc.Option) *Client {
	opts := []rc.Option{
		rc.WithBaseURL( "https://gitlab.indexexchange.com/"),
		rc.WithAPI("api/v4/"),
		rc.WithAuthToken("local"),
		rc.WithUserAgent("xlab"),
		rc.WithIsVerbose(false),
	}
	opts = append(opts, overrides...)
	return &Client{
		client: rc.ConnectClient(
			opts...,
		),
	}
}

func (c *Client) Get(ctx context.Context, app App, q UrlQuery) (kam.JSONValue, http.Header, error) {
	return c.GetPathParams(ctx, app, q.Path, q.Params)
}

func (c *Client) GetPathParams(ctx context.Context, app App, path string, params kam.Map) (kam.JSONValue, http.Header, error) {
	v, header, err := rc.GetWithHeader[kam.JSONValue](ctx, app, c.client, path, params.ToQueryParameters())
	if err != nil {
		return kam.JSONValue{}, nil, err
	}
	if v == nil {
		return kam.JSONValue{}, nil, fmt.Errorf("no JSONValue value was returned")
	}
	return *v, header, nil
}

func GetWithHeader[RespT any](
	ctx context.Context,
	app App,
	c *Client,
	path string,
	params kam.Map) (
	*RespT,
	http.Header, error) {
	return rc.GetWithHeader[RespT](ctx, app, c.client, path, params.ToQueryParameters())
}

func GetWithUnmarshal[RespT any](
	ctx context.Context,
	app App,
	c *Client,
	path string,
	params kam.Map,
	unmarshal func(context.Context, rc.App, *resty.Response) (*RespT, error),
) (*RespT, http.Header, error) {
	return rc.GetWithUnmarshal[RespT](
		ctx,
		app,
		c.client,
		path,
		params.ToQueryParameters(),
		unmarshal,
	)
}

type UrlQuery struct {
	Path   string
	Params kam.Map
}

func (q *UrlQuery) Clone() *UrlQuery {
	return &UrlQuery{
		Path:   q.Path,
		Params: q.Params.Clone(),
	}
}

func (q *UrlQuery) String() string {
	return fmt.Sprint(q.Path, q.Params.ToQueryParameters())
}

func NewPageQuery(path string, page int) *UrlQuery {
	return &UrlQuery{
		Path:   path,
		Params: NewPageParams(page),
	}
}

func NewPageParams(page int) kam.Map {
	const per_page = 200
	return kam.Map{
		"order_by":               "id",
		"owned":                  false,
		"page":                   page,
		"per_page":               per_page,
		"sort":                   "asc",
		"statistics":             false,
		"with_custom_attributes": false,
	}
}

type UrlQueryError struct {
	err   error
	query UrlQuery
}

func (e *UrlQueryError) Error() string {
	return fmt.Sprintf("%s\nUrlQuery:\n%s",
		e.err.Error(),
		utils.YamlString(e.query),
	)
}

type JSONCall struct {
	// Request
	Query UrlQuery

	// Response
	Header http.Header
	Val    kam.JSONValue
}

type CallNoError[RespT any] struct {
	// Request
	Query UrlQuery

	// Response
	Header http.Header
	Val    RespT
}

type Call[RespT any] struct {
	// Request
	Query UrlQuery
	Error error

	// Response
	Header http.Header
	Val    RespT
}
