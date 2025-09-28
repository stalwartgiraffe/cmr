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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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

func TestHandleMergeRequests(t *testing.T) {
	// Create test server
	server := NewServer()
	defer server.Close()

	// Test basic endpoint
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/merge_requests", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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

	// Test invalid path (missing group ID and endpoint)
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d for invalid path, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestHandleProjects(t *testing.T) {
	// Create test server
	server := NewServer()
	defer server.Close()

	// Test basic endpoint
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/projects", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

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
	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify we got projects
	if len(projects) == 0 {
		t.Error("Expected at least one project in response")
	}

	// Verify project structure
	if len(projects) > 0 {
		project := projects[0]
		if project.ID == 0 {
			t.Error("Expected project to have non-zero ID")
		}
		if project.Name == "" {
			t.Error("Expected project to have name")
		}
		if project.Path == "" {
			t.Error("Expected project to have path")
		}
		if project.Visibility == "" {
			t.Error("Expected project to have visibility")
		}
		if project.CreatedAt.IsZero() {
			t.Error("Expected project to have created_at timestamp")
		}
	}
}

func TestHandleProjectsWithQueryParams(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test with query parameters
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/projects?visibility=public&archived=false&simple=true", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// With visibility=public filter, we should get only public projects
	for _, project := range projects {
		if project.Visibility != "public" {
			t.Errorf("Expected only public projects, got visibility: %s", project.Visibility)
		}
		// In simple mode, only basic fields should be populated
		if project.Description != "" {
			t.Error("Expected simple mode to exclude description field")
		}
	}
}

func TestHandleProjectsArchivedFilter(t *testing.T) {
	server := NewServer()
	defer server.Close()

	// Test archived=true filter
	resp, err := http.Get(fmt.Sprintf("%s/api/v4/groups/123/projects?archived=true", server.URL()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// With archived=true filter, we should get only archived projects
	for _, project := range projects {
		if !project.Archived {
			t.Error("Expected only archived projects")
		}
	}
}

