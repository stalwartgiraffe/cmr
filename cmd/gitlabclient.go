package cmd

import (
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

func NewGitlabClientWithURL(accessToken string, baseURL string) *gitlab.Client {
	const isVerbose = false
	return NewGitlabClientWithParams(accessToken, baseURL, isVerbose)
}

// TODO simplify NewClientWithParams to use functional options

func NewGitlabClientWithParams(accessToken string, baseURL string, isVerbose bool) *gitlab.Client {
	return gitlab.NewClientWithParams(
		baseURL,
		"api/v4/",
		accessToken,
		"xlab",
		isVerbose,
	)
}
