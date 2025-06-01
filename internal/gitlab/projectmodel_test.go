package gitlab

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RepoFilePath", func() {
	DescribeTable("test RepoFilePath",
		func(root string, home string, p ProjectModel, want string) {
			have := RepoFilePath(root, home, p)
			Expect(want).To(Equal(have))
		},
		Entry(nil, "/home/annie", "", ProjectModel{Path: "worker", PathWithNamespace: "dev/code/worker"}, "/home/annie/repos/dev/code/worker"),
		Entry(nil, "/home/annie", "cmr", ProjectModel{Path: "worker", PathWithNamespace: "dev/code/worker"}, "/home/annie/cmr/repos/dev/code/worker"),
		Entry(nil, "/home/annie", "/topdir/cmr", ProjectModel{Path: "worker", PathWithNamespace: "dev/code/worker"}, "/topdir/cmr/repos/dev/code/worker"),
		Entry(nil, "/home/annie", "cmr/foo", ProjectModel{Path: "worker", PathWithNamespace: "dev/code/worker"}, "/home/annie/cmr/foo/repos/dev/code/worker"),
		Entry(nil, "/home/annie", "cmr/foo", ProjectModel{Path: "xxx", PathWithNamespace: "dev/code/yyy"}, "/home/annie/cmr/foo/repos/dev/code/yyy"),
		Entry(nil, "/home/annie", "cmr/foo", ProjectModel{Path: "xxx", PathWithNamespace: "dev/code/"}, "/home/annie/cmr/foo/repos/dev/code"),
	)
})
