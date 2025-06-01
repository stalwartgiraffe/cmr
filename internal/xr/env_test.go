package xr

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "xr_env_test")
}

var _ = Describe("expandVar", func() {
	DescribeTable("expandable vars",
		func(env []string, v string, want string) {
			Expect(expandVar(env, v)).To(Equal(want))
		},
		Entry(nil, []string{"a=123", "x=456"}, "x", "456"),
		Entry(nil, []string{"a=123", "xyz=456"}, "xyz", "456"),
		Entry(nil, []string{}, "x", ""),
		Entry(nil, []string{"a=123", "xyzx=456"}, "xyz", ""),
		Entry(nil, []string{"a=123", "xyzx=456"}, "xyz", ""),
		Entry(nil, []string{"a=123", "x=456"}, "z", ""),
	)
})

var _ = Describe("expandArgs", func() {
	DescribeTable("expandable args",
		func(env []string, args []string, want []string) {
			Expect(expandArgs(env, args)).To(Equal(want))
		},
		Entry(nil, []string{"a=123", "x=456"}, []string{"ab", "${x}"}, []string{"ab", "456"}),
		Entry(nil, []string{"a=123", "xyz=456"}, []string{"ab", "${xyz}", "cd"}, []string{"ab", "456", "cd"}),
		Entry(nil, []string{"xyz=456", "a=123"}, []string{"ab", "cd", "${xyz}"}, []string{"ab", "cd", "456"}),
		Entry(nil, []string{"xyz=456", "a=123"}, []string{"ab", "cd", "${zzz}"}, []string{"ab", "cd", ""}),
		Entry(nil, []string{"xyz=456", "a=123"}, []string{"ab", "cd", "${a123}"}, []string{"ab", "cd", ""}),
		Entry(nil, []string{"xyz=456", "a123=456"}, []string{"ab", "cd", "${a123}"}, []string{"ab", "cd", "456"}),
		Entry(nil, []string{"xyz=456", "GO_LINTERS=abc,def"}, []string{"ab", "cd", "${GO_LINTERS}"}, []string{"ab", "cd", "abc,def"}),
		Entry(nil, []string{"xyz=456", "GO_LINTERS=abc,def"}, []string{"ab", "cd", `"${GO_LINTERS}"`}, []string{"ab", "cd", "abc,def"}),
	)
})
