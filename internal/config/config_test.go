package config

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigRegex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "config_regex")
}

var _ = Describe("args tokenizer", func() {
	DescribeTable("command line arguments can be split",
		func(txt string, args []string) {
			Expect(splitCmdArgs(txt)).To(Equal(args))
		},
		Entry(nil, "x", []string{"x"}),
		Entry(nil, "xx", []string{"xx"}),
		Entry(nil, "xx  yy", []string{"xx", "yy"}),
		Entry(nil, `xx "ab cd" yy`, []string{"xx", `"ab cd"`, "yy"}),
		Entry(nil, `xx "ab cd" yy "e"`, []string{"xx", `"ab cd"`, "yy", `"e"`}),
		Entry(nil, `xx "ab cd`, []string{"xx", `"ab`, `cd`}), // improper quote
	)

	DescribeTable("bad arguments for splitting",
		func(txt string) {
			args := splitCmdArgs(txt)
			hasErr := false
			for _, a := range args {
				err := verifyQuote(a)
				hasErr = hasErr || err != nil
			}
			Expect(hasErr).To(BeTrue())
		},
		Entry(nil, `-lint run -n --enable "${GO_LINTERS}  --max-same-issues 0`),
	)
})

var _ = Describe("string is quoted correctly", func() {
	DescribeTable("command line arguments can be split",
		func(txt string, expectErr bool) {
			Expect(verifyQuote(txt) != nil).To(Equal(expectErr))
		},

		Entry(nil, ``, false),
		Entry(nil, "x", false),
		Entry(nil, "x123 -arg d", false),
		Entry(nil, `""`, false),
		Entry(nil, `"x"`, false),
		Entry(nil, `"x a b c "`, false),

		Entry(nil, `"`, true),

		Entry(nil, ` "x"`, true),
		Entry(nil, `"x" `, true),
		Entry(nil, `"x"x`, true),
		Entry(nil, `"x" x`, true),
	)
})
