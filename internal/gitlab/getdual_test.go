package gitlab

import (
	"context"
	"embed"
	"io"
	"io/ioutil"
	"log"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	"github.com/stalwartgiraffe/cmr/restclient"

	appfixtures "github.com/stalwartgiraffe/cmr/internal/app/fixtures"
)

//go:embed data/group01.json
//go:embed data/group02.json
//go:embed data/subgroup.json
var loadTestsFS embed.FS

func makeClient(
	haveResponses responseMap,
	haveErr error,
) *Client {
	baseURL := "https://gitlab.com/"
	api := "api/"
	accessToken := "looksligit"
	userAgent := "cmr"
	isVerbose := false
	client := NewClientWithRest(
		restclient.NewClientMock(
			haveResponses,
			haveErr,
		),
		baseURL,
		api,
		accessToken,
		userAgent,
		isVerbose,
	)
	Expect(client).ToNot(BeNil())
	return client
}

func makeGetPages[RespT any](
	ctx context.Context,
	channelCapacity int,
	firstQueries chan UrlQuery,
	haveResponses responseMap,
	haveErr error,
	totalPageLimit int,
) (
	<-chan CallNoError[RespT],
	<-chan UrlQuery,
	<-chan error,
) {
	client := makeClient(haveResponses, haveErr)
	callCap := channelCapacity
	queryCap := channelCapacity
	errorCap := channelCapacity

	app := appfixtures.NewApp()
	calls, queries, errors := headPageQueriesDual[RespT](
		ctx,
		app,
		client,
		firstQueries,
		callCap,
		queryCap,
		errorCap,
		totalPageLimit,
	)
	Expect(calls).ToNot(BeNil())
	Expect(queries).ToNot(BeNil())
	Expect(errors).ToNot(BeNil())
	return calls, queries, errors
}

func makeGetPageCalls[RespT any](
	ctx context.Context,
	app App,
	channelCapacity int,
	firstQueries chan UrlQuery,
	haveResponses responseMap,
	haveErr error,
) (<-chan CallNoError[RespT],
	<-chan error) {
	client := makeClient(haveResponses, haveErr)
	workersCap := channelCapacity
	errorCap := channelCapacity

	calls, errors := tailPageCallsDual[RespT](
		ctx,
		app,
		client,
		firstQueries,
		workersCap,
		errorCap,
	)
	Expect(calls).ToNot(BeNil())
	Expect(errors).ToNot(BeNil())
	return calls, errors
}

var _ = Describe("for each data file, loadtests", func() {
	It("test every file", func() {
		file, err := loadTestsFS.Open("data/group01.json")
		Expect(err).To(BeNil())
		defer file.Close()

		buf, err := ioutil.ReadAll(file)
		Expect(err).To(BeNil())
		Expect(buf).To(Not(BeNil()))
		err = utils.WalkFileReaders(loadTestsFS, func(path string, file io.Reader) {
		})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("client test of get page queries kam", func() {
	var _ = Describe("ok response", func() {
		groupBuf, err := loadTestsFS.ReadFile("data/group01.json")
		Expect(err).To(BeNil())
		group01txt := string(groupBuf)

		channelCapacity := 5
		var haveErr error
		const totalPageLimit = 5
		It("one responses", func(ctx SpecContext) {
			p := "groups"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: makePagedResponse("api/"+p, group01txt),
			}
			calls, queries, errors := makeGetPages[kam.JSONValue](ctx, channelCapacity, firstQueries, haveResponses, haveErr, totalPageLimit)
			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)
			Consistently(ctx, queries).ShouldNot(Receive())
			Eventually(ctx, queries).Should(BeClosed())
			Consistently(ctx, errors).ShouldNot(Receive())
			Eventually(ctx, errors).Should(BeClosed())

			c := <-calls
			ar := c.Val.AnyVal.(kam.Array)
			Expect(ar).ToNot(BeNil())
			Expect(len(ar)).To(Equal(2))

			m0, err := kam.AsMap(ar[0])
			Expect(err).To(Succeed())
			id, ok := m0["id"]
			Expect(ok).To(BeTrue())
			Expect(id).To(BeNumerically("==", 101))

			m1, err := kam.AsMap(ar[1])
			Expect(err).To(Succeed())
			id, ok = m1["id"]
			Expect(ok).To(BeTrue())
			Expect(id).To(BeNumerically("==", 102))

			Eventually(ctx, calls).Should(BeClosed())
		})
	})
})

var _ = Describe("client test of get page queries RespT", func() {
	var _ = Describe("ok response", func() {
		groupBuf, err := loadTestsFS.ReadFile("data/group01.json")
		Expect(err).To(BeNil())
		group01txt := string(groupBuf)

		channelCapacity := 5
		var haveErr error
		const totalPageLimit = 5
		It("one responses", func(ctx SpecContext) {
			p := "groups"
			firstQueries := make(chan UrlQuery)
			haveResponses := map[string]*resty.Response{
				"api/" + p: makePagedResponse("api/"+p, group01txt),
			}

			calls, queries, errors := makeGetPages[[]GroupModel](ctx, channelCapacity, firstQueries, haveResponses, haveErr, totalPageLimit)
			firstQueries <- UrlQuery{
				Path: p,
			}
			close(firstQueries)

			Consistently(ctx, queries).ShouldNot(Receive())
			Eventually(ctx, queries).Should(BeClosed())
			Consistently(ctx, errors).ShouldNot(Receive())
			Eventually(ctx, errors).Should(BeClosed())

			c := <-calls
			ar := c.Val
			Expect(ar).ToNot(BeNil())
			Expect(len(ar)).To(Equal(2))

			Expect(ar[0].ID).To(Equal(101))
			Expect(ar[1].ID).To(Equal(102))

			Eventually(ctx, calls).Should(BeClosed())
		})
	})
})

func makeDataResponses() map[string]*resty.Response {
	groupBuf, err := loadTestsFS.ReadFile("data/group01.json")
	Expect(err).To(BeNil())
	group01txt := string(groupBuf)
	groupBuf, err = loadTestsFS.ReadFile("data/group02.json")
	Expect(err).To(BeNil())
	group02txt := string(groupBuf)
	return map[string]*resty.Response{
		"api/group01": makePagedResponse("api/group01", group01txt),
		"api/group02": makePagedResponse("api/group02", group02txt),
	}
}

var _ = Describe("test gather page calls RespT", func() {
	var _ = Describe("ok response", func() {
		channelCapacity := 5
		It("two responses", func(ctx SpecContext) {
			firstQueries := make(chan UrlQuery)
			var haveErr error
			app := appfixtures.NewApp()
			calls, errors := makeGetPageCalls[[]GroupModel](
				ctx,
				app,
				channelCapacity,
				firstQueries,
				makeDataResponses(),
				haveErr,
			)
			firstQueries <- UrlQuery{
				Path: "group01",
			}
			firstQueries <- UrlQuery{
				Path: "group02",
			}

			Consistently(ctx, errors).ShouldNot(Receive())

			resp := make(map[string]CallNoError[[]GroupModel])
			for i := 0; i < 2; i++ {
				var c CallNoError[[]GroupModel]
				Eventually(calls).Should(Receive(&c))
				_, ok := resp[c.Query.Path]
				Expect(ok).To(BeFalse())
				resp[c.Query.Path] = c
			}

			Expect(len(resp)).To(Equal(2))
			Expect(len(resp["group01"].Val)).To(Equal(2))
			Expect(resp["group01"].Val[0].ID).To(Equal(101))
			Expect(resp["group01"].Val[1].ID).To(Equal(102))
			Expect(len(resp["group02"].Val)).To(Equal(2))
			Expect(resp["group02"].Val[0].ID).To(Equal(201))
			Expect(resp["group02"].Val[1].ID).To(Equal(202))

			close(firstQueries)
		})
	})
})

func makeGatherAllCalls[RespT any](
	ctx context.Context,
	app App,
	channelCapacity int,
	firstQueries chan UrlQuery,
	haveResponses responseMap,
	haveErr error,
	totalPagesLimit int,
) (<-chan CallNoError[RespT],
	<-chan error) {

	client := makeClient(haveResponses, haveErr)

	callCap := channelCapacity
	queryCap := channelCapacity
	workersCap := channelCapacity
	errorCap := channelCapacity

	calls, errors := GatherPageCallsWithDual[RespT](
		ctx,
		app,
		client,
		firstQueries,
		callCap,
		queryCap,
		workersCap,
		errorCap,
		totalPagesLimit, // 0 means no limit
	)

	Expect(calls).ToNot(BeNil())
	Expect(errors).ToNot(BeNil())
	return calls, errors
}

// go test -v --ginkgo.focus "gather all calls"
var _ = Describe("test gather all calls RespT", func() {
	var _ = Describe("ok response", func() {
		const channelCapacity = 3
		const totalPagesLimit = 5
		It("two responses", func(ctx SpecContext) {
			firstQueries := make(chan UrlQuery)
			app := appfixtures.NewApp()
			var haveErr error
			calls, errors := makeGatherAllCalls[[]GroupModel](
				ctx,
				app,
				channelCapacity,
				firstQueries,
				makeDataResponses(),
				haveErr,
				totalPagesLimit,
			)

			Expect(errors).ToNot(BeNil())
			firstQueries <- UrlQuery{
				Path: "group01",
			}
			firstQueries <- UrlQuery{
				Path: "group02",
			}
			close(firstQueries)

			resp := make(map[string]CallNoError[[]GroupModel])
			for i := 0; i < 2; i++ {
				var c CallNoError[[]GroupModel]
				Eventually(calls).Should(Receive(&c))
				_, ok := resp[c.Query.Path]
				Expect(ok).To(BeFalse())
				resp[c.Query.Path] = c
			}

			Expect(len(resp)).To(Equal(2))
			Expect(len(resp["group01"].Val)).To(Equal(2))
			Expect(resp["group01"].Val[0].ID).To(Equal(101))
			Expect(resp["group01"].Val[1].ID).To(Equal(102))
			Expect(len(resp["group02"].Val)).To(Equal(2))
			Expect(resp["group02"].Val[0].ID).To(Equal(201))
			Expect(resp["group02"].Val[1].ID).To(Equal(202))

			Consistently(ctx, errors).ShouldNot(Receive())
			Eventually(ctx, errors).Should(BeClosed())
			log.Println("calls are closed")

		})
	})
})
