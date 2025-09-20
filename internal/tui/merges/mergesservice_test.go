package merges

import (
	//"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/find"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

func TestUpdateFind(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "test_case_name",
			pattern: "add installing jq to ansible task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//filepath := "ignore/my_recent_merge_request.yaml"
			filepath := "/home/karl.meissner/dev/cmr/ignore/my_recent_merge_request.yaml"

			mergesMap, err := gitlab.NewMergeRequestMapFromYaml(filepath)
			require.NoError(t, err)

			filepath = "/home/karl.meissner/dev/cmr/ignore/projects.yaml"
			projectsSlice, err := gitlab.ReadProjectsSlice(filepath)
			require.NoError(t, err)
			projects := gitlab.MakeProjectMap(projectsSlice)

			require.NoError(t, err)
			contents := NewMergeRequestContents()
			table := NewRecordTable(contents, mergesMap, projects)
			view := find.NewTableView(table)
			require.NotNil(t, view)
			view.UpdateFind(tt.pattern)
		})
	}
}
