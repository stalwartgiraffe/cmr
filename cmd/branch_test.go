package cmd

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

type BranchSuite struct {
	suite.Suite
}

// https://gianarb.it/blog/golang-mockmania-cli-command-with-cobra
func (s *BranchSuite) SetupSuite() {
	r := s.Require()
	// test can you even see the internet
	ok, err := IsURLStatusOK("http://github.com", 3*time.Second)
	r.True(ok, "URL does not respond OK: %v", err)
}

func TestRunBranchSuite(t *testing.T) {
	suite.Run(t, new(BranchSuite))
}

func IsURLStatusOK(url string, timeout time.Duration) (bool, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	if response, err := client.Head(url); err != nil {
		return false, err
	} else {
		return response.StatusCode == http.StatusOK, nil
	}
}

func dumpDirTree(fs billy.Filesystem, parent string) error {
	// Iterate over the files and directories in the root directory.
	dirs, err := fs.ReadDir(parent)
	if err != nil {
		return err
	}

	for _, entry := range dirs {
		child := filepath.Join(parent, entry.Name())
		if entry.IsDir() {
			fmt.Println("Directory:", child)
			err := dumpDirTree(fs, child)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("File:", child)
		}
	}

	return nil
}

func (s *BranchSuite) TestClone() {
	r := s.Require()
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	storer := memory.NewStorage()

	// Clones the repository into the worktree (fs) and stores all the .git
	// content into the storer
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: "https://github.com/git-fixtures/basic.git",
	})

	r.NoError(err)

	dumpDirTree(fs, "/")
}

func setupGitRepo(t *testing.T, rootFS billy.Filesystem) {
	r := require.New(t)
	rootFS.MkdirAll("/repo/.git", os.ModeDir)

	testFS, _ := rootFS.Chroot("/repo")
	dotFS, _ := testFS.Chroot("/.git")
	storage := filesystem.NewStorage(dotFS, cache.NewObjectLRUDefault())
	repo, err := git.Init(storage, testFS)
	r.NoError(err)
	r.NotNil(repo)

	// default to main
	head := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName("refs/heads/main"))
	err = storage.SetReference(head)
	r.NoError(err)
	// config needs to be created after setting up a "normal" fs repo
	// this is possibly a bug in git-go?
	cfg, err := repo.Config()
	r.NoError(err)

	cfg.Init.DefaultBranch = "main"
	err = storage.SetConfig(cfg)
	r.NoError(err)
	worktree, err := repo.Worktree()
	r.NoError(err)

	testFS.MkdirAll("/foo/bar", os.ModeDir)
	f, err := testFS.Create("/foo/bar/hi.txt")
	r.NoError(err)
	_, err = f.Write([]byte("hello world"))
	r.NoError(err)
	_, err = worktree.Add(f.Name())
	r.NoError(err)
	hash, err := worktree.Commit("initial commit", &git.CommitOptions{Author: &object.Signature{}})
	r.NoError(err)
	ref, err := repo.CreateTag("v1", hash, nil)
	r.NoError(err)
	testHashes := map[string]string{}
	testHashes["v1"] = hash.String()
	branchName := plumbing.NewBranchReferenceName("mybranch")
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchName,
		Hash:   ref.Hash(),
		Create: true,
	})
	r.NoError(err)

	f, err = testFS.Create("/secondfile.txt")
	r.NoError(err)
	_, err = f.Write([]byte("another file\n"))
	r.NoError(err)

	n := f.Name()
	_, err = worktree.Add(n)
	r.NoError(err)
	hash, err = worktree.Commit("second commit", &git.CommitOptions{
		Author: &object.Signature{
			Name: "John Doe",
		},
	})
	ref = plumbing.NewHashReference(branchName, hash)
	r.NoError(err)
	testHashes["mybranch"] = ref.Hash().String()

	// make the repo dirty
	_, err = f.Write([]byte("dirty file"))
	r.NoError(err)
	// set up a bare repo
	rootFS.MkdirAll("/bare.git", os.ModeDir)
	rootFS.MkdirAll("/barewt", os.ModeDir)
	bareFS, _ := rootFS.Chroot("/barewt")
	dotFS, _ = rootFS.Chroot("/bare.git")
	storage = filesystem.NewStorage(dotFS, nil)
	repo, err = git.Init(storage, bareFS)
	r.NoError(err)
	worktree, err = repo.Worktree()
	r.NoError(err)

	f, err = bareFS.Create("/hello.txt")
	r.NoError(err)
	f.Write([]byte("hello world"))
	worktree.Add(f.Name())

	_, err = worktree.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	r.NoError(err)
}
func TestSetupGitRepo(t *testing.T) {
	rootFS := memfs.New()
	setupGitRepo(t, rootFS)
}

// const defaultPublicMode = 0744 // rwxr--r--
func (s *BranchSuite) TestInitLocalRepo() {
	/*
			r := s.Require()
			fs := memfs.New()
			err := fs.MkdirAll("code", defaultPublicMode)
			r.NoError(err)
			storer := memory.NewStorage()
			repo, err := git.Init(storer, fs)
			r.NoError(err)
			r.NotNil(repo)

			worktree, err := repo.Worktree()
			r.NoError(err)
			err = worktree.Checkout(&git.CheckoutOptions{
				Create: true,
				Branch: "master",
			})
			r.NoError(err)

			//testBranch := &config.Branch{
			//Name:   "foo",
			//Remote: "origin",
			//Merge:  "refs/heads/foo",
			//}
			//err := repo.CreateBranch(testBranch)

			// Get the HEAD reference.
			//err = repo.Storer.SetReference(nil)
			//err = repo.Storer.C
			//r.NoError(err)
			headRef, err := repo.Head()
			r.NoError(err)

			jiraIssue := "JIRA-1234"
			branchLabel := "first_commit"
			name := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s_%s", jiraIssue, branchLabel))
			ref := plumbing.NewHashReference(name, headRef.Hash())
			err = repo.Storer.SetReference(ref)
			r.NoError(err)

			worktree, err = repo.Worktree()
			r.NoError(err)

			filepath := "/code/README.md"
			f, err := fs.Create(filepath)
			r.NoError(err)
			fmt.Fprintln(f,
				`# Test Repo

		Repository for test code.`)
			f.Close()

			_, err = worktree.Add(filepath)
			r.NoError(err)

			status, err := worktree.Status()
			r.NoError(err)

			fmt.Println(status)
			fmt.Println("---------------")

			dumpDirTree(fs, "/")
	*/

	/*
		// Filesystem abstraction based on memory
		fs := memfs.New()
		r, _ := git.Init(memory.NewStorage(), fs)


		// Add a new remote, with the default fetch refspec
		_, err :=repo.CreateRemote(&config.RemoteConfig{
			Name: "example",
			URLs: []string{"https://github.com/git-fixtures/basic.git"},
		})

		if err != nil {
			log.Fatal(err)
		}
	*/

	//list, err :=repo.Remotes()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, r := range list {
	//	fmt.Println("my tst", r)
	//}

	/*
		r := require.New(t)
		cases := []struct {
			name    string
			txt     string
			want    JiraIssue
			wantErr bool
		}{}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_ = tc
				r.Equal(2+2, 4)
			})
		}
	*/
}
