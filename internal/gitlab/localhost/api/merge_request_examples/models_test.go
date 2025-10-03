package merge_request_examples

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalMergeRequests(t *testing.T) {
	// Find all JSON test files
	files, err := filepath.Glob("merge_request_*.json")
	if err != nil {
		t.Fatalf("Failed to glob JSON files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No merge request JSON files found")
	}

	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			data, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", filename, err)
			}

			var mrs MergeRequestList
			if err := json.Unmarshal(data, &mrs); err != nil {
				t.Errorf("Failed to unmarshal %s: %v", filename, err)
				return
			}

			if len(mrs) == 0 {
				t.Errorf("Expected at least one merge request in %s", filename)
				return
			}

			// Validate required fields are populated
			mr := mrs[0]
			if mr.ID == 0 {
				t.Errorf("%s: ID is zero", filename)
			}
			if mr.IID == 0 {
				t.Errorf("%s: IID is zero", filename)
			}
			if mr.ProjectID == 0 {
				t.Errorf("%s: ProjectID is zero", filename)
			}
			if mr.Title == "" {
				t.Errorf("%s: Title is empty", filename)
			}
			if mr.State == "" {
				t.Errorf("%s: State is empty", filename)
			}
			if mr.Author.ID == 0 {
				t.Errorf("%s: Author ID is zero", filename)
			}

			// Test round-trip marshaling
			marshaled, err := json.Marshal(mrs)
			if err != nil {
				t.Errorf("Failed to marshal %s: %v", filename, err)
				return
			}

			var mrs2 MergeRequestList
			if err := json.Unmarshal(marshaled, &mrs2); err != nil {
				t.Errorf("Failed to unmarshal round-trip data for %s: %v", filename, err)
			}
		})
	}
}

func TestMergeRequestStructure(t *testing.T) {
	// Test with a minimal valid merge request
	jsonData := `[{
		"id": 1,
		"iid": 1,
		"project_id": 1,
		"title": "Test MR",
		"description": "Test description",
		"state": "opened",
		"created_at": "2025-10-02T13:25:31.269-04:00",
		"updated_at": "2025-10-02T13:32:23.750-04:00",
		"merged_by": null,
		"merge_user": null,
		"merged_at": null,
		"closed_by": null,
		"closed_at": null,
		"target_branch": "main",
		"source_branch": "feature",
		"user_notes_count": 0,
		"upvotes": 0,
		"downvotes": 0,
		"author": {
			"id": 1,
			"username": "test",
			"name": "Test User",
			"state": "active",
			"avatar_url": null,
			"web_url": "https://example.com/test"
		},
		"assignees": [],
		"assignee": null,
		"reviewers": [],
		"source_project_id": 1,
		"target_project_id": 1,
		"labels": [],
		"draft": false,
		"work_in_progress": false,
		"milestone": null,
		"merge_when_pipeline_succeeds": false,
		"merge_status": "unchecked",
		"detailed_merge_status": "unchecked",
		"sha": "0000000000000000000000000000000000000000",
		"merge_commit_sha": null,
		"squash_commit_sha": null,
		"discussion_locked": null,
		"should_remove_source_branch": null,
		"force_remove_source_branch": false,
		"reference": "!1",
		"references": {
			"short": "!1",
			"relative": "!1",
			"full": "test/project!1"
		},
		"web_url": "https://example.com/mr/1",
		"time_stats": {
			"time_estimate": 0,
			"total_time_spent": 0,
			"human_time_estimate": null,
			"human_total_time_spent": null
		},
		"squash": false,
		"squash_on_merge": false,
		"task_completion_status": {
			"count": 0,
			"completed_count": 0
		},
		"has_conflicts": false,
		"blocking_discussions_resolved": true,
		"approvals_before_merge": null
	}]`

	var mrs MergeRequestList
	if err := json.Unmarshal([]byte(jsonData), &mrs); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	if len(mrs) != 1 {
		t.Fatalf("Expected 1 merge request, got %d", len(mrs))
	}

	mr := mrs[0]
	if mr.ID != 1 {
		t.Errorf("Expected ID 1, got %d", mr.ID)
	}
	if mr.Title != "Test MR" {
		t.Errorf("Expected title 'Test MR', got %s", mr.Title)
	}
	if mr.Author.Username != "test" {
		t.Errorf("Expected author username 'test', got %s", mr.Author.Username)
	}
	if mr.MergedBy != nil {
		t.Errorf("Expected MergedBy to be nil, got %+v", mr.MergedBy)
	}
}
