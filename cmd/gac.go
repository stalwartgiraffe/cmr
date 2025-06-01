package cmd

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/gitutil"
	"github.com/stalwartgiraffe/cmr/internal/prompts"
	"github.com/stalwartgiraffe/cmr/withstack"
)

// initCmd represents the init command
func NewGacCommand(cfg *CmdConfig, repo *git.Repository) *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "gac",
		Short: "git add and commit",
		Long:  `git add and commit.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			RunGac()
		},
	}

	return pushCmd
}
func RunGac() {
	repo, err := gitutil.OpenCwd()
	if err != nil {
		fmt.Println(err)
	}
	err = runRepoGac(repo)
	if err != nil {
		fmt.Println(err)
	}
}

// files with these status in worktree get added to staging index
var worktreeFilter = []git.StatusCode{
	git.Untracked,
	git.Modified,
	git.Deleted,
	git.Renamed,
	git.Copied,
}

// files with these status in staging index get commited to local repo
var stagingFilter = []git.StatusCode{
	git.Modified,
	git.Added,
	git.Deleted,
	git.Renamed,
	git.Copied,
}

func runRepoGac(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return withstack.Errorf("Could not get worktree: %w", err)
	}

	workTreeMatches, err := gitutil.FindByWorktreeStatus(worktree, worktreeFilter)
	if err != nil {
		return err
	}

	if 0 < len(workTreeMatches) {
		worktreeFiles, err := prompts.SelectFiles(workTreeMatches, prompts.ToWorkTreeStatus)
		if err != nil {
			return err
		}
		if len(worktreeFiles) < 1 {
			return nil
		}

		if err = gitutil.AddAll(worktree, worktreeFiles); err != nil {
			return err
		}
	}

	filePaths, commitMsg, err := getCommit(repo, worktree)
	if err != nil {
		return err
	}
	if len(filePaths) < 1 {
		return nil
	}

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	_, err = worktree.Commit(commitMsg,
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
	if err != nil {
		return err
	}

	// push using default options
	//return repo.Gac(&git.GacOptions{})
	return nil
}

func getCommit(repo *git.Repository, worktree *git.Worktree) ([]string, string, error) {
	issue, description, err := getJiraTitleFromBranch(repo)
	if err != nil {
		return nil, "", err
	}
	stagingMatches, err := gitutil.FindByStagingStatus(worktree, stagingFilter)
	if err != nil {
		return nil, "", err
	}
	filePaths, opTxt, descriptionTxt, issueTxt, err := prompts.CommitToLocal(stagingMatches, issue, description)
	if err != nil {
		return nil, "", err
	}
	if len(filePaths) < 1 {
		return nil, "", nil
	}

	msg := formatConventionalCommit(opTxt, descriptionTxt, issueTxt)
	return filePaths, msg, nil
}

func formatConventionalCommit(opTxt, descriptionTxt, issueTxt string) string {
	return fmt.Sprintf("%s: %s [%s]", opTxt, descriptionTxt, issueTxt)
}

// see https://nomix.gpages.indexexchange.com/arc3/doc/pol/commit-message-guidelines/
// see https://gitlab.indexexchange.com/exchange-node/gitlab-ci-modules/blob/main/validate_mainline_mr.yaml
// https://commitlint.js.org/#/
// https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional

// ['sentence-case', 'start-case', 'pascal-case', 'upper-case']
// very annoying at this time 2023-12-19  support for acronyms is still a feature request
// https://github.com/conventional-changelog/commitlint/issues/3312
func getJiraTitleFromBranch(repo *git.Repository) (string, string, error) {
	branchName, err := gitutil.BranchShortName(repo)
	if err != nil {
		return "", "", withstack.Errorf("Could not get branch name: %w", err)
	}

	issue, description := gitutil.ParseBranchJiraTitle(branchName)
	return issue, description, nil
}
