package prompts

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("jira issue", func() {
	DescribeTable("strings that match jira",
		func(txt string, want bool) {
			b := []byte(txt)
			Expect(jiraIssueRE.Match(b)).To(Equal(want))
		},
		Entry(nil, "X-1", true),
		Entry(nil, "JIRAS-12345", true),

		Entry(nil, "", false),
		Entry(nil, "j", false),
		Entry(nil, "jra-12345", false),
		Entry(nil, "JIRA 12345", false),
		Entry(nil, "JIRA 12345", false),
		Entry(nil, "JIRASX-12345", false),
		Entry(nil, "JIRAs-12345", false),
		Entry(nil, "JIRAS-X2345", false),
		Entry(nil, "JIRAS-12345 ", false),
	)
})

var _ = Describe("sentence case", func() {
	DescribeTable("strings that match sentence casing",
		func(txt string, want bool) {
			b := []byte(txt)
			Expect(sentenceCaseRE.Match(b)).To(Equal(want))
		},
		Entry(nil, "x", true),
		Entry(nil, "X", true),
		Entry(nil, "xxx", true),
		Entry(nil, "Xxx", true),
		Entry(nil, "Xxx yy0 0123456789  : aaa", true),
		Entry(nil, "42", true),
		Entry(nil, "42   abc   123", true),

		//Entry(nil, "", false),
		Entry(nil, "XYZ", false),
		Entry(nil, "Xxx Yyy", false),
	)
})

var _ = Describe("start case", func() {
	DescribeTable("strings that match start casing",
		func(txt string, want bool) {
			b := []byte(txt)
			Expect(startRe.Match(b)).To(Equal(want))
		},
		Entry(nil, "X", true),
		Entry(nil, "Xx Yy", true),
		Entry(nil, "Xx   Yy   Zzz", true),
		Entry(nil, "Xxx Yy0 0123456789  :  Aaa", true),

		Entry(nil, "x", false),
		Entry(nil, "Xx yy Zz", false),
	)
})

var _ = Describe("pascal case", func() {
	DescribeTable("strings that match upper casing",
		func(txt string, want bool) {
			b := []byte(txt)
			Expect(pascalRe.Match(b)).To(Equal(want))
		},
		Entry(nil, "Xx", true),
		Entry(nil, "XxYyZz", true),
		Entry(nil, "Xx0123aa", true),

		Entry(nil, "X", false),
	)
})

var _ = Describe("upper case", func() {
	DescribeTable("strings that match upper casing",
		func(txt string, want bool) {
			b := []byte(txt)
			Expect(upperRe.Match(b)).To(Equal(want))
		},
		Entry(nil, "XX", true),
		Entry(nil, "XX YYY Z", true),
		Entry(nil, "XZ 0123 A9Z", true),

		Entry(nil, "x", false),
	)
})
