package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stalwartgiraffe/cmr/internal/app/fixtures"
	"github.com/stalwartgiraffe/cmr/internal/gitlab/localhost"
)

func TestGetEvents(t *testing.T) {

	server := localhost.NewServer()
	defer server.Close()

	app := fixtures.NewApp()
	ctx, cancel := context.WithCancel(context.Background())

	lastDateStr := "2025-01-01"

	accessToken := "local"
	ec := NewEventClientWithURL(accessToken, server.URL())

	route := "/api/v4/users/1/events"
	recentEvents, err := ec.getEvents(
		ctx,
		app,
		cancel,
		route,
		lastDateStr)

	require.NoError(t, err)
	require.NotNil(t, recentEvents)
}
