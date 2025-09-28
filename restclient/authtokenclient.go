package restclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
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

type AuthTokenClient struct {
	client    Client
	baseURL   string
	userAgent string
	api       string
	authToken string // the manually managed bearer token.
	isVerbose bool
	isDebug   bool
	headers   map[string]string
}

func NewWithParams(
	baseURL string,
	api string,
	authToken string,
	userAgent string,
	isVerbose bool,
) *AuthTokenClient {
	client := newClientAdapter()
	c := &AuthTokenClient{
		client:    client,
		baseURL:   baseURL,
		userAgent: userAgent,
		api:       api,
		isVerbose: isVerbose,
		headers:   make(map[string]string),

		// Note that by default the resty.Client uses a golang CookieJar.
		// The cookie jar manager the session cookies
		// https://pkg.go.dev/net/http/cookiejar
		authToken: authToken,
	}
	return c
}

type Option func(*AuthTokenClient)

func WithBaseURL(baseURL string) Option {
	return func(c *AuthTokenClient) {
		c.baseURL = baseURL
	}
}

func WithAPI(api string) Option {
	return func(c *AuthTokenClient) {
		c.api = api
	}
}

func WithAuthToken(authToken string) Option {
	return func(c *AuthTokenClient) {
		c.authToken = authToken
	}
}

func WithIsVerbose(isVerbose bool) Option {
	return func(c *AuthTokenClient) {
		c.isVerbose = isVerbose
	}
}

func WithIsDebug(isDebug bool) Option {
	return func(c *AuthTokenClient) {
		c.isDebug = isDebug
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *AuthTokenClient) {
		c.userAgent = userAgent
	}
}

func WithHeader(k, v string) Option {
	return func(c *AuthTokenClient) {
		c.headers[k] = v
	}
}

func WithClient(client Client) Option {
	return func(c *AuthTokenClient) {
		if client == nil {
			panic("WithClient passed nil client")
		}

		c.client = client
		if len(c.baseURL) < 1 {
			panic("baseURL not set")
		}
		c.client.SetBaseURL(c.baseURL)
		if len(c.userAgent) < 1 {
			panic("userAgent not set")
		}
		if 0 < len(c.userAgent) {
			c.client.SetHeader("User-Agent", c.userAgent)
		}

		for k, v := range c.headers {
			c.client.SetHeader(k, v)
		}
	}
}

func ConnectClient(opts ...Option) *AuthTokenClient {
	c := &AuthTokenClient{}
	for _, opt := range opts {
		opt(c)
	}

	if c.client == nil {
		WithClient(newClientAdapter())(c)
	}
	return c
}

type App interface {
	Tracer
	Logger
}

type Tracer interface {
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption) (
		context.Context,
		trace.Span)
}

type Logger interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
}

//func (c *AuthTokenClient) BaseApiPath() string {
// 	s, err := url.JoinPath(c.Client.GetBaseURL(), c.API)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to join base URL %s with API %s: %v	", c.Client.GetBaseURL(), c.API, err))
// 	}
// 	return s
// }

// Generic member functions are not natively support in Go1.19
// see https://github.com/golang/go/issues/49085
func getResponse(
	ctx context.Context,
	app App,
	tokenClient *AuthTokenClient,
	path string,
	queries string) (
	*resty.Response,
	error,
) {
	// just accept json from the web server for now
	const accept = "application/json"
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

	r := tokenClient.client.Request()
	r.SetContext(ctx).
		SetHeader("Accept", accept)
	if tokenClient.authToken != "" {
		r = r.SetAuthToken(tokenClient.authToken)
	}
	if queries != "" {
		r = r.SetQueryString(queries)
	}

	// TODO safely join the paths here
	// https://stackoverflow.com/questions/34668012/combine-url-paths-with-path-join
	resp, err := r.Get(tokenClient.api + path)
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
	const includeBody = true
	bb, err := httputil.DumpResponse(resp.RawResponse, includeBody)
	if err == nil {
		return string(bb)
	} else {
		return err.Error()
	}
}

func Get[RespT any](ctx context.Context,
	app App,
	tokenClient *AuthTokenClient,
	path string,
	query string) (
	*RespT, error,
) {
	resp, err := getResponse(ctx, app, tokenClient, path, query)
	if err != nil {
		return nil, err
	}
	if tokenClient.isVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
	}
	resp.Header()
	return Unmarshal[RespT](resp)
}

func GetWithHeader[RespT any](
	ctx context.Context,
	app App,
	tokenClient *AuthTokenClient,
	path string,
	queries string) (
	*RespT, http.Header,
	error,
) {
	ctx, span := app.StartSpan(ctx, "GetWithHeader")
	defer span.End()

	resp, err := getResponse(ctx, app, tokenClient, path, queries)
	if err != nil {
		return nil, nil, err
	}
	if tokenClient.isVerbose {
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
	tokenClient *AuthTokenClient,
	path string,
	queries string,
	unmarshal func(context.Context, App, *resty.Response) (*RespT, error),
) (
	*RespT, http.Header,
	error,
) {
	ctx, span := app.StartSpan(ctx, "GetWithUnmarshal")
	defer span.End()

	resp, err := getResponse(ctx, app, tokenClient, path, queries)
	if err != nil {
		return nil, nil, err
	}
	if tokenClient.isVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
	}
	r, err := unmarshal(ctx, app, resp)
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

func Update[BodyT any, RespT any](ctx context.Context, op int, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	r := tokenClient.client.Request().
		SetContext(ctx).
		SetHeader("Accept", "application/json")
	if tokenClient.authToken != "" {
		r = r.SetAuthToken(tokenClient.authToken)
	}

	rb := r.SetBody(b)
	p := tokenClient.api + path
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

	if tokenClient.isVerbose {
		fmt.Println(color.Ize(rstClr, SprintRequestQuiet(resp)))
		if tokenClient.isDebug {
			cookies := resp.Cookies()
			for _, c := range cookies {
				fmt.Println(color.Ize(rstClr, fmt.Sprintf("%+v", *c)))
			}
		}
	}
	return Unmarshal[RespT](resp)
}

func Head[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, HEAD, tokenClient, path, b)
}
func Post[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, POST, tokenClient, path, b)
}

func Put[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, PUT, tokenClient, path, b)
}
func Delete[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, DELETE, tokenClient, path, b)
}
func Options[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, OPTIONS, tokenClient, path, b)
}
func Patch[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, error) {
	return Update[BodyT, RespT](ctx, PATCH, tokenClient, path, b)
}

func PostReturnCookies[BodyT any, RespT any](ctx context.Context, tokenClient *AuthTokenClient, path string, b *BodyT) (*RespT, []*http.Cookie, error) {
	r := tokenClient.client.Request().
		SetContext(ctx).
		SetHeader("Accept", "application/json")
	if tokenClient.authToken != "" {
		r = r.SetAuthToken(tokenClient.authToken)
	}

	resp, err := r.
		SetBody(b).
		Post(tokenClient.api + path)
	if err != nil {
		return nil, nil, withstack.Errorf("POST error:%w", err)
	}

	if tokenClient.isVerbose {

		fmt.Println(color.Ize(color.Cyan, SprintRequestQuiet(resp)))
		cookies := resp.Cookies()
		for _, c := range cookies {
			fmt.Println(color.Ize(rstClr, fmt.Sprintf("%+v", *c)))
		}
	}
	tt, err := Unmarshal[RespT](resp)
	return tt, resp.Cookies(), err
}
