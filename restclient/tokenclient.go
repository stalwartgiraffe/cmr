package restclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/stalwartgiraffe/cmr/withstack"
)

var (
	rstClr = color.Purple
)

type App interface {
	Tracer
}

type Tracer interface {
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption) (
		context.Context,
		trace.Span)
}

type TokenClient struct {
	Client      Client
	Api         string
	AccessToken string // the manually managed bearer token.
	IsVerbose   bool
	IsDebug     bool
}

func New(
	baseURL string,
	api string,
	accessToken string,
	userAgent string,
	isVerbose bool) *TokenClient {

	client := newClientAdapter()
	return NewWithClient(
		client,
		baseURL,
		api,
		accessToken,
		userAgent,
		isVerbose)
}

func NewWithClient(
	client Client,
	baseURL string,
	api string,
	accessToken string,
	userAgent string,
	isVerbose bool) *TokenClient {

	// Note that by default the resty.Client uses a golang CookieJar.
	// The cookie jar manager the session cookies
	// https://pkg.go.dev/net/http/cookiejar
	tokenClient := TokenClient{
		Client:      client,
		Api:         api,
		AccessToken: accessToken,
		IsVerbose:   isVerbose,
	}
	tokenClient.Client.SetBaseURL(baseURL)
	tokenClient.Client.SetHeader("User-Agent", userAgent)
	tokenClient.Client.SetHeader("ix-custom", "testix")
	return &tokenClient
}

func (c *TokenClient) WithAPI(api string) *TokenClient {
	wrap := *c
	wrap.Api = api
	return &wrap
}

func (c *TokenClient) BaseApiPath() string {
	s, err := url.JoinPath(c.Client.GetBaseURL(), c.Api)
	if err != nil {
		panic(fmt.Sprintf("Failed to join base URL %s with API %s: %v	", c.Client.GetBaseURL(), c.Api, err))
	}
	return s
}

// Generic member functions are not natively support in Go1.19
// see https://github.com/golang/go/issues/49085
func GetResponse(
	ctx context.Context,
	tokenClient *TokenClient,
	path string,
	query string) (
	*resty.Response,
	error,
) {
	// just accept json from the web server for now
	return GetResponseWithAccept(ctx, nil, tokenClient, path, query, "application/json")
}

// Generic member functions are not natively support in Go1.19
// see https://github.com/golang/go/issues/49085
func GetResponseWithApp(
	ctx context.Context,
	app App,
	tokenClient *TokenClient,
	path string,
	query string) (
	*resty.Response,
	error,
) {
	// just accept json from the web server for now
	return GetResponseWithAccept(ctx, app, tokenClient, path, query, "application/json")
}

func GetResponseWithAccept(
	ctx context.Context,
	app App,
	tokenClient *TokenClient,
	path string,
	queries string,
	accept string) (
	*resty.Response,
	error,
) {
	if app != nil {
		var span trace.Span
		ctx, span = app.StartSpan(ctx, "GetResponseWithAccept")
		defer span.End()

		attributes := []attribute.KeyValue{
			attribute.String("path", path),
		}
		if queries != "" {
			attributes = append(attributes, queryAsKV(queries)...)
		}
		span.SetAttributes(attributes...)
	}

	r := tokenClient.Client.Request()
	r.SetContext(ctx).
		SetHeader("Accept", accept)
	if tokenClient.AccessToken != "" {
		r = r.SetAuthToken(tokenClient.AccessToken)
	}
	if queries != "" {
		r = r.SetQueryString(queries)
	}

	// TODO safely join the paths here
	// https://stackoverflow.com/questions/34668012/combine-url-paths-with-path-join
	resp, err := r.Get(tokenClient.Api + path)
	if err != nil {
		return resp, withstack.Errorf("Path error:%w", err)
	}
	return resp, nil
}

func queryAsKV(queries string) []attribute.KeyValue {
	attributes := []attribute.KeyValue{}
	for _, query := range strings.Split(queries, "&") {
		kv := strings.Split(query, "=")
		if len(kv) == 2 {

			attributes = append(attributes,
				attribute.String(kv[0], kv[1]),
			)
		}
	}
	return attributes
}

func GetString(
	ctx context.Context,
	tokenClient *TokenClient,
	path string,
	query string) (
	string, error) {
	resp, err := GetResponse(ctx, tokenClient, path, query)
	if err != nil {
		return "", err
	}
	if tokenClient.IsVerbose {
		fmt.Println("GET string on", color.Ize(rstClr, resp.Request.URL))
	}
	return string(resp.Body()), nil
}

// SprintRequestQuiet dumps the raw request that generated this response.
// Quietly hide headers Authorization and Cookie because they are usually too noisy.
func SprintRequestQuiet(resp *resty.Response) string {
	ctx := context.Background()
	req := resp.Request.RawRequest.Clone(ctx)
	req.Header.Del("Authorization") // hide
	req.Header.Del("Cookie")        // hide
	bb, err := httputil.DumpRequest(req, true)
	if err == nil {
		return string(bb)
	} else {
		return err.Error()
	}
}
func SprintResponse(resp *resty.Response) string {
	bb, err := httputil.DumpResponse(resp.RawResponse, true)
	if err == nil {
		return string(bb)
	} else {
		return err.Error()
	}
}

func Get[RespT any](ctx context.Context,
	tokenClient *TokenClient,
	path string,
	query string) (
	*RespT, error,
) {
	resp, err := GetResponse(ctx, tokenClient, path, query)
	if err != nil {
		return nil, err
	}
	if tokenClient.IsVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
	}
	resp.Header()
	return Unmarshal[RespT](resp)
}

func GetWithHeader[RespT any](ctx context.Context,
	tokenClient *TokenClient,
	path string,
	queries string) (
	*RespT, http.Header,
	error,
) {
	return GetWithHeaderWithApp[RespT](
		ctx,
		nil,
		tokenClient,
		path,
		queries)
}

func GetWithHeaderWithApp[RespT any](
	ctx context.Context,
	app App,
	tokenClient *TokenClient,
	path string,
	query string) (
	*RespT, http.Header,
	error,
) {
	resp, err := GetResponseWithApp(ctx, app, tokenClient, path, query)
	if err != nil {
		return nil, nil, err
	}
	if tokenClient.IsVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
	}
	r, err := Unmarshal[RespT](resp)
	if err != nil {
		return nil, nil, err
	}
	return r, resp.Header(), err
}

func GetWithUnmarshal[RespT any](ctx context.Context,
	app App,
	tokenClient *TokenClient,
	path string,
	query string,
	unmarshal func(*resty.Response) (*RespT, error),
) (
	*RespT, http.Header,
	error,
) {
	resp, err := GetResponseWithApp(ctx, app, tokenClient, path, query)
	if err != nil {
		return nil, nil, err
	}
	if tokenClient.IsVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
	}
	//r, err := Unmarshal[RespT](resp)
	r, err := unmarshal(resp)
	if err != nil {
		return nil, nil, err
	}
	return r, resp.Header(), err
}

type VerbFn func(req *resty.Request, url string) (*resty.Response, error)

const (
	HEAD int = iota
	POST
	PUT
	DELETE
	OPTIONS
	PATCH
)

func Update[BodyT any, RespT any](ctx context.Context, op int, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	r := tokenClient.Client.Request().
		SetContext(ctx).
		SetHeader("Accept", "application/json")
	if tokenClient.AccessToken != "" {
		r = r.SetAuthToken(tokenClient.AccessToken)
	}

	rb := r.SetBody(b)
	p := tokenClient.Api + path
	var err error
	var resp *resty.Response
	switch op {
	case HEAD:
		resp, err = rb.Head(p)
	case POST:
		resp, err = rb.Post(p)
	case PUT:
		resp, err = rb.Put(p)
	case DELETE:
		resp, err = rb.Delete(p)
	case OPTIONS:
		resp, err = rb.Options(p)
	case PATCH:
		resp, err = rb.Patch(p)
	default:
		panic(fmt.Sprintf("Unknown op %d", op))
	}
	if err != nil {
		return nil, withstack.Errorf("POST error:%w", err)
	}

	if tokenClient.IsVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
		if tokenClient.IsDebug {
			cookies := resp.Cookies()
			for _, c := range cookies {
				fmt.Println(color.Ize(rstClr, fmt.Sprintf("%+v", *c)))
			}
		}
	}
	return Unmarshal[RespT](resp)
}

func Head[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, HEAD, tokenClient, path, b)
}
func Post[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, POST, tokenClient, path, b)
}

func Put[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, PUT, tokenClient, path, b)
}
func Delete[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, DELETE, tokenClient, path, b)
}
func Options[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, OPTIONS, tokenClient, path, b)
}
func Patch[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, PATCH, tokenClient, path, b)
}

func PostReturnCookies[BodyT any, RespT any](ctx context.Context, tokenClient *TokenClient, path string, b *BodyT) (*RespT, []*http.Cookie, error) {
	r := tokenClient.Client.Request().
		SetContext(ctx).
		SetHeader("Accept", "application/json")
	if tokenClient.AccessToken != "" {
		r = r.SetAuthToken(tokenClient.AccessToken)
	}

	resp, err := r.
		SetBody(b).
		Post(tokenClient.Api + path)
	if err != nil {
		return nil, nil, withstack.Errorf("POST error:%w", err)
	}

	if tokenClient.IsVerbose {

		fmt.Println(color.Ize(color.Cyan, SprintRequestQuiet(resp)))
		cookies := resp.Cookies()
		for _, c := range cookies {
			fmt.Println(color.Ize(rstClr, fmt.Sprintf("%+v", *c)))
		}
	}
	tt, err := Unmarshal[RespT](resp)
	return tt, resp.Cookies(), err
}
