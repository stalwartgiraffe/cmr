package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/gitutil"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func addBranchCommand(parent *cobra.Command) {
	parent.AddCommand(NewBranchCommand(nil))
}

func NewBranchCommand(repo *git.Repository) *cobra.Command {
	return &cobra.Command{
		Use:   "branch [jira] [name]",
		Short: "Start a new branch with a jira and name",
		Long: `A new branch.
FIXME: detailed explanation.
FIXME: include examples.
`,

		Args: func(cmd *cobra.Command, args []string) error {
			_, _, err := parseBranchArgs(args)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRepoBranch(cmd, args, repo)
		},
	}
}

func runRepoBranch(cmd *cobra.Command, args []string, repo *git.Repository) error {
	jiraIssue, branchLabel, err := parseBranchArgs(args)
	if err != nil {
		return err
	}

	if repo == nil {
		repo, err = gitutil.OpenCwd()
		if err != nil {
			return err
		}
	}

	// Get the HEAD reference.
	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	// Create a new plumbing.HashReference object with the name of the branch
	// and the hash from the HEAD. The reference name should be a full reference
	// name and not an abbreviated one, as is used on the git cli.
	//
	// For tags we should use `refs/tags/%s` instead of `refs/heads/%s` used
	// for branches.
	branchShortName := fmt.Sprintf("%s_%s", jiraIssue, branchLabel)
	branchRefName := plumbing.NewBranchReferenceName(branchShortName)
	branchRef := plumbing.NewHashReference(branchRefName, headRef.Hash())

	// The created reference is saved in the storage.
	if err = repo.Storer.SetReference(branchRef); err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	// git checkout branchShortName
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRefName,
		Create: false,
	})
	if err != nil {
		return err
	}

	// Print the HEAD reference.
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "make branch %s %s\n", jiraIssue, branchLabel)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), branchRef)
	if err != nil {
		return err
	}

	return nil
}

func match[T ~string](re *regexp.Regexp, txt string, p *T) error {
	if !re.MatchString(txt) {
		return fmt.Errorf("%s is not in the expected format", txt)
	}
	*p = T(txt)
	return nil
}

func parseBranchArgs(args []string) (JiraIssue, BranchLabel, error) {
	if len(args) != 2 {
		return "", "", errors.New("requires two arguments")
	}

	var jiraIssue JiraIssue
	if err := match(jiraIssueRegexp, args[0], &jiraIssue); err != nil {
		return "", "", err
	}
	var branchLabel BranchLabel
	if err := match(branchLabelRegexp, args[1], &branchLabel); err != nil {
		return "", "", err
	}
	return jiraIssue, branchLabel, nil
}

type JiraIssue string

// jiraIssueRegexp is legal jira issue:
// DEALS-1234
// OPD-12345
var jiraIssueRegexp = regexp.MustCompile(`^[A-Z]+\-\d{4,5}$`)

type BranchLabel string

// branchLabelRegexp is a legal label
// foo
// foo_bar
// a_important_pr
var branchLabelRegexp = regexp.MustCompile(`^([a-z]+\_)*([a-z]+)$`)
