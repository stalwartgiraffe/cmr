package gitlab

import (
	"net/http"
	"strconv"

	"github.com/stalwartgiraffe/cmr/withstack"
)

type pageCursor struct {
	page       *int
	nextPage   *int
	totalPages *int
	perPage    *int
	prevPage   *int
	totalItems *int
}

func (c *pageCursor) remainingPageIndex() []int {
	if c.page == nil || c.totalPages == nil {
		return nil
	}
	idx := []int{}
	p := *c.page + 1
	n := *c.totalPages
	for ; p <= n; p++ {
		idx = append(idx, p)
	}
	return idx
}

func parseHeaderInts(h http.Header, key string) ([]int, error) {
	strVals, ok := h[key]
	if !ok {
		return nil, withstack.Errorf("Header key %s not found", key)
	}

	vals := []int{}
	for _, s := range strVals {
		if s == "" {
			continue
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}
func parseOneHeaderInt(err error, h http.Header, key string, val **int) error {
	if err != nil {
		return err
	}

	vals, err := parseHeaderInts(h, key)
	if err != nil {
		return err
	}
	n := len(vals)
	if n == 0 {
		return nil
	}
	if 1 < n {
		return withstack.Errorf("unexpected additional values in %s", key)
	}
	i := vals[0]
	*val = &i
	return nil
}

func parsePageCursor(h http.Header) (pageCursor, error) {
	const (
		page       = "X-Page"
		nextPage   = "X-Next-Page"
		prevPage   = "X-Prev-Page"
		totalPages = "X-Total-Pages"
		perPage    = "X-Per-Page"
		total      = "X-Total"
	)

	var err error
	p := pageCursor{}
	err = parseOneHeaderInt(err, h, page, &p.page)
	err = parseOneHeaderInt(err, h, nextPage, &p.nextPage)
	err = parseOneHeaderInt(err, h, prevPage, &p.prevPage)
	err = parseOneHeaderInt(err, h, totalPages, &p.totalPages)
	err = parseOneHeaderInt(err, h, perPage, &p.perPage)
	err = parseOneHeaderInt(err, h, total, &p.totalItems)
	return p, err
}
