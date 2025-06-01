package gitutil

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"

	"github.com/go-git/go-billy/v5"
	//"github.com/go-git/go-billy/v5/osfs"

	"github.com/stalwartgiraffe/cmr/withstack"
)

func OpenCwd() (*git.Repository, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, withstack.Errorf("Error getting current working directory: %w", err)
	}

	return PlainOpen(path)
}

func PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func MakeEmptyRepo(rootFS billy.Filesystem) (*git.Repository, error) {
	const defaultBranch = "main"
	const commitMsg = "Initial commit."
	const tagName = "v1.0.0"
	const authorName = "Annie Mouse"
	return MakeEmptyRepoWithBranchCommitTag(rootFS, defaultBranch, commitMsg, tagName, authorName)
}

//git ls-remote <repository-url>
//git clone --depth=1 <repository-url>

// NewEmptyCommitOptions sets commit options for an empty commit.
func NewEmptyCommitOptions(authorName string) *git.CommitOptions {
	return &git.CommitOptions{
		// empty commit means no refs have been added to worktree
		AllowEmptyCommits: true, // override default
		Author: &object.Signature{
			Name: authorName,
			When: time.Now(),
		}}
}

func Clone(directory, url, token string, progress io.Writer) error {
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth: &http.BasicAuth{
			Username: "auto", // yes, this can be anything except an empty string
			Password: token,
		},
		Progress: progress,
	})
	return err
}

func Pull(directory, url, token string, progress io.Writer) error {
	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	fmt.Println("git pull origin")
	err = w.Pull(
		&git.PullOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: "auto", // yes, this can be anything except an empty string
				Password: token,
			},
			Progress: progress,
		})
	if err != nil {
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	fmt.Println(commit)
	return nil
}

func MakeEmptyRepoWithBranchCommitTag(
	rootFS billy.Filesystem,
	branchName string,
	commitMsg string,
	tagName string,
	authorName string,
) (*git.Repository, error) {

	if err := rootFS.MkdirAll("/.git", os.ModeDir); err != nil {
		return nil, withstack.Errorf("%w", err)
	}
	dotFS, _ := rootFS.Chroot("/.git")
	storage := filesystem.NewStorage(dotFS, cache.NewObjectLRUDefault())

	var closeErr error
	defer func() {
		closeErr = storage.Close()
		if closeErr != nil {
			closeErr = withstack.Errorf("%w", closeErr)
		}
	}()

	//  git init
	repo, err := git.Init(storage, rootFS)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	head := plumbing.NewSymbolicReference(plumbing.HEAD, refName)
	err = storage.SetReference(head)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	cfg, err := repo.Config()
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	cfg.Init.DefaultBranch = branchName
	if err := storage.SetConfig(cfg); err != nil {
		return nil, withstack.Errorf("%w", err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	hash, err := worktree.Commit(
		commitMsg,
		NewEmptyCommitOptions(authorName),
	)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}
	_, err = repo.CreateTag(tagName, hash, nil)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}
	return repo, closeErr
}

var oldJiraRE = regexp.MustCompile(`^JIRA\d{4}$`)

// isJiraLabel return true on JIRA1234
func isJiraLabel(s string) bool {
	return oldJiraRE.MatchString(s)
}

var branchLabelRE = regexp.MustCompile(`^[a-zA-Z0-9]+(_[a-zA-Z0-9]+)*$`)

// isBranchLabel return true on alpha numeric words separated by underscore: snake_case_Labels01
func isBranchLabel(s string) bool {
	return branchLabelRE.MatchString(s)
}

func CheckoutCreateBranchJiraLabel(repo *git.Repository, jiraIssue, branchLabel string) error {
	if isJiraLabel(jiraIssue) {
		return fmt.Errorf("not in expected jira format:%s", jiraIssue)
	}
	if !isBranchLabel(branchLabel) {
		return fmt.Errorf("not in expected branch label format:%s", branchLabel)
	}
	return CheckoutCreateBranch(repo, fmt.Sprintf("%s_%s", jiraIssue, branchLabel))
}

// CheckoutCreateBranch
// Info("git branch my-branch")
func CheckoutCreateBranch(repo *git.Repository, branchShortName string) error {
	// Get the HEAD reference.
	headRef, err := repo.Head()
	if err != nil {
		return withstack.Errorf("%w", err)
	}

	// Create a new plumbing.HashReference object with the name of the branch
	// and the hash from the HEAD. The reference name should be a full reference
	// name and not an abbreviated one, as is used on the git cli.
	//
	// For tags we should use `refs/tags/%s` instead of `refs/heads/%s` used
	// for branches.
	branchRefName := plumbing.NewBranchReferenceName(branchShortName)
	branchRef := plumbing.NewHashReference(branchRefName, headRef.Hash())

	// The created reference is saved in the storage.
	err = repo.Storer.SetReference(branchRef)
	if err != nil {
		return withstack.Errorf("%w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRefName,
		Create: false,
	})
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	return nil
}

func Push(fs billy.Filesystem) error {
	repo, err := OpenRepo(fs)
	if err != nil {
		return err
	}
	// push using default options
	return repo.Push(&git.PushOptions{})
}

func OpenRepo(fs billy.Filesystem) (*git.Repository, error) {
	if _, err := fs.Stat(git.GitDirName); err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	dotFS, err := fs.Chroot(git.GitDirName)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}

	storage := filesystem.NewStorageWithOptions(dotFS,
		cache.NewObjectLRUDefault(),
		filesystem.Options{KeepDescriptors: true})
	repo, err := git.Open(storage, dotFS)
	if err != nil {
		return nil, withstack.Errorf("%w", err)
	}
	return repo, nil
}

func BranchShortName(repo *git.Repository) (string, error) {
	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	return headRef.Name().Short(), nil
}

type statusMatchFn func(status *git.FileStatus, codes []git.StatusCode) bool

func matchStatus(code git.StatusCode, codes []git.StatusCode) bool {
	for _, c := range codes {
		if code == c {
			return true
		}
	}
	return false
}

func FindByWorktreeStatus(worktree *git.Worktree, codes []git.StatusCode) (git.Status, error) {
	return findByStatus(worktree, codes, func(status *git.FileStatus, codes []git.StatusCode) bool {
		return matchStatus(status.Worktree, codes)
	})
}

func FindByStagingStatus(worktree *git.Worktree, codes []git.StatusCode) (git.Status, error) {
	return findByStatus(worktree, codes, func(status *git.FileStatus, codes []git.StatusCode) bool {
		return matchStatus(status.Staging, codes)
	})
}

// findByStatus return the git.Status  map of files that match.
func findByStatus(worktree *git.Worktree, codes []git.StatusCode, match statusMatchFn) (git.Status, error) {
	all, err := worktree.Status()
	if err != nil {
		return nil, withstack.Errorf("Could not get status: %w", err)
	}
	//fmt.Println("------------------")
	//spew.Dump(all)

	matches := git.Status{}
	for filepath, fileStatus := range all {
		if match(fileStatus, codes) {
			cp := *fileStatus
			matches[filepath] = &cp
		}
	}

	return matches, nil
}

func FileStatusString(fs git.FileStatus) string {
	return fmt.Sprintf("{ staging: %c,  worktree: %c,  extra: %s }",
		fs.Staging,  // status of this file in staging area
		fs.Worktree, // status  of this file in the worktree area
		fs.Extra,
	)
}

var jiraRE = regexp.MustCompile(`[A-Z]{3,}-[0-9]{2,4}`)

func ParseBranchJiraTitle(branchName string) (string, string) {
	issue := jiraRE.FindString(branchName)
	if 0 < len(issue) {
		branchName = strings.Replace(branchName, issue, "", 1)
	}
	branchName = strings.Replace(branchName, "_", " ", -1)
	words := strings.Fields(branchName)
	description := strings.Join(words, " ")
	return issue, description
}

// AddAll adds all the files or directory to the staging index.
func AddAll(worktree *git.Worktree, filePaths []string) error {
	// FIXME? on error use git reset to rollback?
	for _, f := range filePaths {
		_, err := worktree.Add(f)
		if err != nil {
			return withstack.Errorf("Could not add %s to staging index: %w", f, err)
		}
	}
	return nil
}
