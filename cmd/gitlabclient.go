package cmd

import (
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

func NewGitlabClientWithURL(authToken string, baseURL string) *gitlab.Client {
	const isVerbose = false
	return NewGitlabClientWithParams(authToken, baseURL, isVerbose)
}

// TODO simplify NewClientWithParams to use functional options

func NewGitlabClientWithParams(authToken string, baseURL string, isVerbose bool) *gitlab.Client {
	return gitlab.NewClientWithParams(
		baseURL,
		"api/v4/",
		authToken,
		"xlab",
		isVerbose,
	)
}
