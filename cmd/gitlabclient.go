package cmd

import (
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

func NewGitlabClientWithURL(accessToken string, baseURL string)  *gitlab.Client {
	const isVerbose = false
		return gitlab.NewClientWithParams(
		baseURL,
		"api/v4/",
		accessToken,
		"xlab",
		isVerbose,
	)
}
