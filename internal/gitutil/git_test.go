package gitutil

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
)

func TestBranch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "git_test")
}

var _ = Describe("branch labels", func() {

	DescribeTable("snake_case_labels",
		func(label string, isValid bool) {
			Expect(isBranchLabel(label)).To(Equal(isValid))
		},
		Entry(nil, "x", true),
		Entry(nil, "xxx", true),
		Entry(nil, "x_x", true),
		Entry(nil, "0", true),
		Entry(nil, "0_9", true),
		Entry(nil, "mylabel", true),
		Entry(nil, "my_label", true),
		Entry(nil, "My_Label", true),
		Entry(nil, "my_label_0123456789", true),

		Entry(nil, "", false),
		Entry(nil, "my label", false),
		Entry(nil, " my_label", false),
		Entry(nil, "my_label ", false),
		Entry(nil, "_", false),
		Entry(nil, "_x", false),
		Entry(nil, "x_", false),
	)

	DescribeTable("JIRA1234",
		func(label string, isValid bool) {
			Expect(isJiraLabel(label)).To(Equal(isValid))
		},
		Entry(nil, "JIRA1234", true),
		Entry(nil, "JIRA0000", true),
		Entry(nil, "JIRA9999", true),

		Entry(nil, "", false),
		Entry(nil, " JIRA1234", false),
		Entry(nil, "JIRA1234 ", false),
		Entry(nil, "JIRA1234 ", false),
		Entry(nil, "JIRA123", false),
		Entry(nil, "JIRA123x", false),
		Entry(nil, "JIRA123_", false),
		Entry(nil, "JIRA1234x", false),
		Entry(nil, "DORA1234", false),
		Entry(nil, "_", false),
		Entry(nil, "_x", false),
		Entry(nil, "x_", false),
	)

})

var _ = Describe("Create an empty repo", func() {
	Context("in memory fs", func() {
		var rootFS billy.Filesystem
		const branchName = "main"
		const commitMsg = "Initial commit."
		const tagName = "v1.0.0"
		const authorName = "Annie Mouse"
		BeforeEach(func() {
			rootFS = memfs.New()
		})
		It("empty repo", func() {
			repo, err := MakeEmptyRepoWithBranchCommitTag(rootFS, branchName, commitMsg, tagName, authorName)
			Expect(err).To(Succeed())
			Expect(repo).ToNot(BeNil())
		})
	})
})

var _ = Describe("CheckoutCreateBranch", func() {
	var err error
	var repo *git.Repository
	var rootFS billy.Filesystem
	const branchName = "main"
	const commitMsg = "Initial commit."
	const tagName = "v1.0.0"
	const authorName = "Annie Mouse"

	BeforeEach(func() {
		rootFS = memfs.New()
		repo, err = MakeEmptyRepoWithBranchCommitTag(rootFS, branchName, commitMsg, tagName, authorName)
		Expect(err).To(Succeed())
		Expect(repo).ToNot(BeNil())
	})
	It("checkout memory fs", func() {
		const shortName = "newbranch"
		err = CheckoutCreateBranch(repo, shortName)
		Expect(err).To(Succeed())

		headRef, err := repo.Head()
		Expect(err).To(Succeed())
		Expect(headRef.Name().IsBranch()).To(Equal(true))
		Expect(headRef.Name().Short()).To(Equal(shortName))

	})
})

var _ = Describe("parse branch jira title", func() {

	DescribeTable("branch namesk",
		func(label string, issue string, description string) {
			haveIssue, haveDescription := ParseBranchJiraTitle(label)
			Expect(haveIssue).To(Equal(issue))
			Expect(haveDescription).To(Equal(description))
		},
		Entry(nil, "x", "", "x"),
		Entry(nil, "xy", "", "xy"),
		Entry(nil, "x_y", "", "x y"),
		Entry(nil, "JIRA-1234_x_y", "JIRA-1234", "x y"),
		Entry(nil, "JIRA-1234__x1?x_y_", "JIRA-1234", "x1?x y"),
		Entry(nil, "ABC-12__xx_y_", "ABC-12", "xx y"),
		Entry(nil, "AB-12__xx_y_", "", "AB-12 xx y"),
		Entry(nil, "xx_y_ABC-12", "ABC-12", "xx y"),
		Entry(nil, "xx_y_ABC-123", "ABC-123", "xx y"),
		Entry(nil, "xx_y_ABC-1234", "ABC-1234", "xx y"),
		Entry(nil, "xx_y_ABC-123456y", "ABC-1234", "xx y 56y"),
	)
})
