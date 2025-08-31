package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app"
)

func TestLoadModels(t *testing.T) {

	appErr := app.NewApp()

	jsonBlob, err := readJsonBlob()
	require.True(t, 0 < len(jsonBlob))
	require.NoError(t, err)

	models, err := LoadModels(appErr.App, jsonBlob)
	require.NoError(t, err)
	require.NotNil(t, models)
}

func readJsonBlob() ([]byte, error) {
	return os.ReadFile("/home/karl.meissner/dev/cmr/cmd/jsonblob.txt")
}
