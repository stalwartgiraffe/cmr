package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/gitutil"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func TestBranch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Branch Suite")
}

var _ = BeforeSuite(func() {
	ok, err := IsURLStatusOK("http://github.com", 3*time.Second)
	Expect(err).To(Succeed())
	Expect(ok).To(Equal(true))
})

var _ = Describe("Create an empty repo", func() {
	Context("set up", func() {
		rootFS := memfs.New()

		It("make a git repo, add a file, check out branch", func() {

			rootFS.MkdirAll("/.git", os.ModeDir)
			dotFS, _ := rootFS.Chroot("/.git")
			storage := filesystem.NewStorage(dotFS, cache.NewObjectLRUDefault())
			repo, err := git.Init(storage, rootFS)
			Expect(err).To(Succeed())

			// default to main
			head := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName("refs/heads/main"))
			Expect(storage.SetReference(head)).To(BeNil())
			cfg, err := repo.Config()
			Expect(err).To(Succeed())
			cfg.Init.DefaultBranch = "main"
			Expect(storage.SetConfig(cfg)).To(Succeed())
			worktree, err := repo.Worktree()
			Expect(err).To(Succeed())

			rootFS.MkdirAll("/foo/bar", os.ModeDir)
			f, err := rootFS.Create("/foo/bar/hi.txt")
			Expect(err).To(Succeed())
			_, err = f.Write([]byte("hello world"))
			Expect(err).To(Succeed())
			_, err = worktree.Add(f.Name())
			Expect(err).To(Succeed())

			hash, err := worktree.Commit(
				"initial commit",
				gitutil.NewEmptyCommitOptions("Annie Mouse"),
			)
			Expect(err).To(Succeed())

			ref, err := repo.CreateTag("v1", hash, nil)
			Expect(err).To(Succeed())
			testHashes := map[string]string{}
			testHashes["v1"] = hash.String()
			myBranchRefName := plumbing.NewBranchReferenceName("mybranch")
			err = worktree.Checkout(&git.CheckoutOptions{
				Branch: myBranchRefName,
				Hash:   ref.Hash(),
				Create: true,
			})
			Expect(err).To(Succeed())
		})

	})
})

func initMemRepo(rootFS billy.Filesystem) *git.Repository {
	rootFS.MkdirAll("/.git", os.ModeDir)
	dotFS, _ := rootFS.Chroot("/.git")
	storage := filesystem.NewStorage(dotFS, cache.NewObjectLRUDefault())
	repo, err := git.Init(storage, rootFS)
	Expect(err).To(Succeed())

	// default to main
	head := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName("refs/heads/main"))
	Expect(storage.SetReference(head)).To(BeNil())
	cfg, err := repo.Config()
	Expect(err).To(Succeed())
	cfg.Init.DefaultBranch = "main"
	Expect(storage.SetConfig(cfg)).To(Succeed())

	worktree, err := repo.Worktree()
	Expect(err).To(Succeed())

	_, err = worktree.Commit(
		"initial commit",
		gitutil.NewEmptyCommitOptions("Annie Mouse"),
	)
	Expect(err).To(Succeed())

	return repo
}

func bufContainsSubstring(b *bytes.Buffer, txt string) {
	ba, err := io.ReadAll(b)
	Expect(err).To(Succeed())
	Expect(string(ba)).To(ContainSubstring(txt))
}

var _ = Describe("Branch Cmd", func() {

	var rootFS billy.Filesystem
	var repo *git.Repository
	var cmd *cobra.Command
	var outBuf *bytes.Buffer
	var errBuf *bytes.Buffer

	BeforeEach(func() {
		rootFS = memfs.New()
		repo = initMemRepo(rootFS)
		cmd = NewBranchCommand(repo)
		outBuf = &bytes.Buffer{}
		cmd.SetOut(outBuf)
		errBuf = &bytes.Buffer{}
		cmd.SetErr(errBuf)
	})

	Describe("make a branch", func() {

		It("bad args", func() {
			cmd.SetArgs([]string{"--badarg", "testisawesome"})
			err := cmd.Execute()
			Expect(err).NotTo(BeNil())

			branchUseTxt := "Usage"
			bufContainsSubstring(outBuf, branchUseTxt)

		})

		It("big deal branch", func() {
			jiraIssue := "DEALS-1234"
			branchLabel := "big_deal"
			branchShortName := fmt.Sprintf("%s_%s", jiraIssue, branchLabel)
			cmd.SetArgs([]string{jiraIssue, branchLabel})
			err := cmd.Execute()
			Expect(err).To(Succeed())

			headRef, err := repo.Head()
			Expect(err).To(Succeed())
			Expect(headRef.Name().IsBranch()).To(Equal(true))
			Expect(headRef.Name().Short()).To(Equal(branchShortName))

		})

	})
})
