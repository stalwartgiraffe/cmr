package gitlab

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stalwartgiraffe/cmr/kam"
)

type responseMap map[string]*resty.Response

func makePagedResponse(route, body string) *resty.Response {
	r := &resty.Response{
		Request: &resty.Request{
			Method:     "GET",
			URL:        route,
			QueryParam: url.Values{},
			FormData:   url.Values{},
		},
		RawResponse: &http.Response{
			Status:     "ok",
			StatusCode: http.StatusOK,
			Header: http.Header{
				"X-Page":        []string{"1"},
				"X-Next-Page":   []string{"1"},
				"X-Prev-Page":   []string{"0"},
				"X-Total-Pages": []string{"1"},
				"X-Per-Page":    []string{"5"},
				"X-Total":       []string{"5"},
			},
		},
	}
	r.SetBody([]byte(body))
	return r
}

var _ = Describe("client test of get page queries", func() {
	var _ = Describe("malformed client response", func() {
		eventuallyOneErr := func(
			cfg context.Context,
			calls <-chan CallNoError[kam.JSONValue],
			queries <-chan UrlQuery,
			errors <-chan error) {
			Consistently(cfg, queries).ShouldNot(Receive())
			Eventually(cfg, queries).Should(BeClosed())
			Consistently(cfg, calls).ShouldNot(Receive())
			Eventually(cfg, calls).Should(BeClosed())
			Eventually(cfg, errors).Should(Not(BeNil()))
			Eventually(cfg, errors).Should(BeClosed())
		}
		channelCapacity := 5
		var haveErr error
		It("no responses", func(ctx SpecContext) {
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)

			firstQueries <- UrlQuery{
				Path: "noresp",
			}
			close(firstQueries)
			eventuallyOneErr(ctx, calls, queries, errors)
		})

		It("one response missing request", func(ctx SpecContext) {
			p := "missingrequest"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: &resty.Response{},
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)

			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			eventuallyOneErr(ctx, calls, queries, errors)
		})

		It("one response bad request", func(ctx SpecContext) {
			p := "badrequest"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: &resty.Response{
					Request: &resty.Request{
						Method:     "GET",
						URL:        "api/" + p,
						QueryParam: url.Values{},
						FormData:   url.Values{},
					},
					RawResponse: &http.Response{
						Status:     "Bad Request",
						StatusCode: http.StatusBadRequest,
					},
				},
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)

			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			eventuallyOneErr(ctx, calls, queries, errors)
		})
		It("one response missing response body", func(ctx SpecContext) {
			p := "missingbody"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: &resty.Response{
					Request: &resty.Request{
						Method:     "GET",
						URL:        "api/" + p,
						QueryParam: url.Values{},
						FormData:   url.Values{},
					},
					RawResponse: &http.Response{
						Status:     "ok",
						StatusCode: http.StatusOK,
					},
				},
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)

			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			eventuallyOneErr(ctx, calls, queries, errors)
		})
		It("one response empty response body missing headers", func(ctx SpecContext) {
			p := "emptybody"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: func() *resty.Response {
					r := &resty.Response{
						Request: &resty.Request{
							Method:     "GET",
							URL:        "api/" + p,
							QueryParam: url.Values{},
							FormData:   url.Values{},
						},
						RawResponse: &http.Response{
							Status:     "ok",
							StatusCode: http.StatusOK,
						},
					}
					r.SetBody([]byte("{}"))
					return r
				}(),
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)

			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			Consistently(ctx, queries).ShouldNot(Receive())
			Eventually(ctx, queries).Should(BeClosed())
			Eventually(ctx, calls).Should(Not(BeNil()))
			Eventually(ctx, calls).Should(BeClosed())
			Eventually(ctx, errors).Should(Not(BeNil()))
			Eventually(ctx, errors).Should(BeClosed())
		})

		It("one response got headers, empty body", func(ctx SpecContext) {
			p := "emptybody"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: makePagedResponse("api/"+p, "{}"),
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr)
			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			Consistently(ctx, queries).ShouldNot(Receive())
			Eventually(ctx, queries).Should(BeClosed())
			Eventually(ctx, calls).Should(Not(BeNil()))
			Eventually(ctx, calls).Should(BeClosed())
			Consistently(ctx, errors).ShouldNot(Receive())
			Eventually(ctx, errors).Should(BeClosed())
		})
	})
})
