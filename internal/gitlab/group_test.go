package gitlab

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stalwartgiraffe/cmr/kam"
)

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal_group_test")
}

var _ = Describe("getVal", func() {
	DescribeTable("get vals from kam",
		func(m kam.Map, k string, wantVal int, wantOk bool) {
			var have int
			ok := getVal(m, k, &have)
			Expect(ok).To(Equal(wantOk))
			if wantOk {
				Expect(have).To(Equal(wantVal))
			}
		},
		Entry(nil, kam.Map{"key": 123}, "key", 123, true),
		Entry(nil, kam.Map{"key": 123, "x": 2}, "key", 123, true),
		Entry(nil, kam.Map{"nope": 123}, "key", 123, false),
		Entry(nil, kam.Map{}, "key", 123, false),
	)
})
