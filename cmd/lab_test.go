package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app/fixtures"
	"github.com/stalwartgiraffe/cmr/internal/gitlab/localhost"
)

func TestGetProjects(t *testing.T) {
	server := localhost.NewServer()
	defer server.Close()
	app := fixtures.NewApp()
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	accessToken := "local"
	client := NewProjectsClient(accessToken, server.URL())
	projects,errs := client.getProjects(
		ctx,
		app)

	require.NoError(t, errs)
	require.NotNil(t, projects)
	require.Equal(t, 3, len(projects))
}
