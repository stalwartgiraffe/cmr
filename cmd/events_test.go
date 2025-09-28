package cmd

import (
	"context"
	"testing"

	"github.com/stalwartgiraffe/cmr/internal/app/fixtures"
	"github.com/stalwartgiraffe/cmr/internal/gitlab/localhost"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	server := localhost.NewServer()
	defer server.Close()

	app := fixtures.NewApp()
	ctx, cancel := context.WithCancel(context.Background())

	lastDateStr := "2025-01-01"

	accessToken := "local"
	client := NewEventClient(accessToken, server.URL())

	//http://127.0.0.1:46067/api/v4//api/v4/users/1/events
	route := "users/1/events"
	recentEvents, err := client.getEvents(
		ctx,
		app,
		cancel,
		route,
		lastDateStr)

	require.NoError(t, err)
	require.NotNil(t, recentEvents)
}
