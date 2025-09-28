package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app/fixtures"
	"github.com/stalwartgiraffe/cmr/internal/gitlab/localhost"
	rc "github.com/stalwartgiraffe/cmr/restclient"
)

func TestGetProjects(t *testing.T) {
	server := localhost.NewServer()
	defer server.Close()
	client := NewProjectsClient(
		rc.WithBaseURL(server.URL()),
		rc.WithAuthToken("local"),
	)
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	app := fixtures.NewApp()
	projects, errs := client.getProjects(
		ctx,
		app)

	require.NoError(t, errs)
	require.NotNil(t, projects)
	require.Equal(t, 3, len(projects))
}
