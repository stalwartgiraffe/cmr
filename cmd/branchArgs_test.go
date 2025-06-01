package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJiraIssue(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		name    string
		txt     string
		want    JiraIssue
		wantErr bool
	}{
		{"deals", "DEALS-1234", "DEALS-1234", false},
		{"opd", "OPD-12345", "OPD-12345", false},

		{"missing error", "", "", true},
		{"casing error", "deals-1234", "", true},
		{"no txt error", "-1234", "", true},
		{"bad txt error", "d!!!s-1234", "", true},
		{"space error", "bad deals-1234", "", true},
		{"missing digits error", "DEALS-", "", true},
		{"few digits error", "DEALS-123", "", true},
		{"excess digits error", "DEALS-123456", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var have JiraIssue
			err := match(jiraIssueRegexp, tc.txt, &have)
			r.Equal(tc.wantErr, err != nil)
			r.Equal(tc.want, have)
		})
	}
}

func TestBrandLabel(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		name    string
		txt     string
		want    BranchLabel
		wantErr bool
	}{
		{"one", "word", "word", false},
		{"two", "aaa_bbb", "aaa_bbb", false},
		{"three", "a_b_c", "a_b_c", false},

		{"underscore 1", "_", "", true},
		{"underscore 2", "a_", "", true},
		{"underscore 3", "_b", "", true},

		{"bad txt", "w!!!d", "", true},

		{"space 1", " wordd", "", true},
		{"space 2", "wo d", "", true},
		{"space 2", "wod ", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var have BranchLabel
			err := match(branchLabelRegexp, tc.txt, &have)
			r.Equal(tc.wantErr, err != nil)
			r.Equal(tc.want, have)
		})
	}
}

func TestParseBranchArgs(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		name       string
		args       []string
		wantJira   JiraIssue
		wantBranch BranchLabel
		wantErr    bool
	}{
		{"ok", []string{"DEALS-1234", "a_label"}, "DEALS-1234", "a_label", false},

		{"nil", nil, "", "", true},
		{"empty", []string{}, "", "", true},
		{"one", []string{"DEALS-1234"}, "", "", true},
		{"three", []string{"DEALS-1234", "a_label", "bad"}, "", "", true},

		{"bad jira", []string{"DEALS", "a_label"}, "", "", true},
		{"bad branch", []string{"DEALS-1234", "_"}, "", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			haveJira, haveBranch, err := parseBranchArgs(tc.args)
			r.Equal(tc.wantJira, haveJira)
			r.Equal(tc.wantBranch, haveBranch)
			r.Equal(tc.wantErr, err != nil)
		})
	}
}
