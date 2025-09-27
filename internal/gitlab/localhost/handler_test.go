package localhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestHandleEvents(t *testing.T) {
	// Create test server
	server := NewServer()
	defer server.Close()

	// Test basic endpoint
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/users/123/events", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify we got events
	if len(events) == 0 {
		t.Error("Expected at least one event in response")
	}

	// Verify event structure
	if len(events) > 0 {
		event := events[0]
		if event.ID == 0 {
			t.Error("Expected event to have non-zero ID")
		}
		if event.ActionName == "" {
			t.Error("Expected event to have action_name")
		}
		if event.AuthorID == 0 {
			t.Error("Expected event to have non-zero author_id")
		}
		if event.CreatedAt.IsZero() {
			t.Error("Expected event to have created_at timestamp")
		}
	}
}

func TestHandleEventsWithQueryParams(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test with query parameters
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/users/123/events?page=1&per_page=2&action=closed&target_type=Issue", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// With action=closed filter, we should get fewer events
	if len(events) > 2 {
		t.Errorf("Expected at most 2 events with per_page=2, got %d", len(events))
	}
}

func TestHandleEventsInvalidUserID(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test invalid path (missing user ID)
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/users/", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid path, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestHandleMergeRequests(t *testing.T) {
	// Create test server
	server := NewServer()
	defer server.Close()

	// Test basic endpoint
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/merge_requests", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Check pagination headers
	if resp.Header.Get("X-Page") != "1" {
		t.Errorf("Expected X-Page header to be 1, got %s", resp.Header.Get("X-Page"))
	}

	// Parse response
	var mergeRequests []MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mergeRequests); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify we got merge requests
	if len(mergeRequests) == 0 {
		t.Error("Expected at least one merge request in response")
	}

	// Verify merge request structure
	if len(mergeRequests) > 0 {
		mr := mergeRequests[0]
		if mr.ID == 0 {
			t.Error("Expected merge request to have non-zero ID")
		}
		if mr.IID == 0 {
			t.Error("Expected merge request to have non-zero IID")
		}
		if mr.Title == "" {
			t.Error("Expected merge request to have title")
		}
		if mr.State == "" {
			t.Error("Expected merge request to have state")
		}
		if mr.Author == nil {
			t.Error("Expected merge request to have author")
		}
		if mr.CreatedAt.IsZero() {
			t.Error("Expected merge request to have created_at timestamp")
		}
	}
}

func TestHandleMergeRequestsWithQueryParams(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test with query parameters
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/merge_requests?state=opened&author_id=25&sort=asc", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var mergeRequests []MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mergeRequests); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// With state=opened filter, we should get only opened merge requests
	for _, mr := range mergeRequests {
		if mr.State != "opened" {
			t.Errorf("Expected only opened merge requests, got state: %s", mr.State)
		}
	}
}

func TestHandleMergeRequestsInvalidGroupID(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test invalid path (missing group ID)
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid path, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}