package localhost

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// Project represents a GitLab project as per API_Entities_Project
type Project struct {
	// CLAUDE fill in the fields based on API_Entities_Project in openapi_v2.yaml
}
// Event represents a GitLab event as per API_Entities_Event
type Event struct {
	ID             int       `json:"id"`
	ProjectID      *int      `json:"project_id,omitempty"`
	ActionName     string    `json:"action_name"`
	TargetID       *int      `json:"target_id,omitempty"`
	TargetIID      *int      `json:"target_iid,omitempty"`
	TargetType     *string   `json:"target_type,omitempty"`
	AuthorID       int       `json:"author_id"`
	TargetTitle    *string   `json:"target_title,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	AuthorUsername *string   `json:"author_username,omitempty"`
	Imported       bool      `json:"imported"`
	ImportedFrom   string    `json:"imported_from"`
}

// UserBasic represents a GitLab user as per API_Entities_UserBasic
type UserBasic struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PublicEmail  string `json:"public_email,omitempty"`
	Name         string `json:"name"`
	State        string `json:"state"`
	Locked       bool   `json:"locked"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	AvatarPath   string `json:"avatar_path,omitempty"`
}

// MergeRequest represents a GitLab merge request as per API_Entities_MergeRequestBasic
type MergeRequest struct {
	ID                             int        `json:"id"`
	IID                            int        `json:"iid"`
	ProjectID                      int        `json:"project_id"`
	Title                          string     `json:"title"`
	Description                    string     `json:"description,omitempty"`
	State                          string     `json:"state"`
	CreatedAt                      time.Time  `json:"created_at"`
	UpdatedAt                      time.Time  `json:"updated_at"`
	MergedBy                       *UserBasic `json:"merged_by,omitempty"`
	MergeUser                      *UserBasic `json:"merge_user,omitempty"`
	MergedAt                       *string    `json:"merged_at,omitempty"`
	ClosedBy                       *UserBasic `json:"closed_by,omitempty"`
	ClosedAt                       *string    `json:"closed_at,omitempty"`
	TitleHTML                      string     `json:"title_html,omitempty"`
	DescriptionHTML                string     `json:"description_html,omitempty"`
	TargetBranch                   string     `json:"target_branch"`
	SourceBranch                   string     `json:"source_branch"`
	UserNotesCount                 string     `json:"user_notes_count,omitempty"`
	Upvotes                        string     `json:"upvotes,omitempty"`
	Downvotes                      string     `json:"downvotes,omitempty"`
	Author                         *UserBasic `json:"author,omitempty"`
	Assignees                      *UserBasic `json:"assignees,omitempty"`
	Assignee                       *UserBasic `json:"assignee,omitempty"`
	Reviewers                      *UserBasic `json:"reviewers,omitempty"`
	SourceProjectID                string     `json:"source_project_id,omitempty"`
	TargetProjectID                string     `json:"target_project_id,omitempty"`
	Labels                         string     `json:"labels,omitempty"`
	Draft                          string     `json:"draft,omitempty"`
	Imported                       string     `json:"imported,omitempty"`
	ImportedFrom                   string     `json:"imported_from,omitempty"`
	WorkInProgress                 string     `json:"work_in_progress,omitempty"`
	MergeWhenPipelineSucceeds      string     `json:"merge_when_pipeline_succeeds,omitempty"`
	MergeStatus                    string     `json:"merge_status,omitempty"`
	DetailedMergeStatus            string     `json:"detailed_merge_status,omitempty"`
	MergeAfter                     string     `json:"merge_after,omitempty"`
	SHA                            string     `json:"sha,omitempty"`
	MergeCommitSHA                 string     `json:"merge_commit_sha,omitempty"`
	SquashCommitSHA                string     `json:"squash_commit_sha,omitempty"`
	DiscussionLocked               string     `json:"discussion_locked,omitempty"`
	ShouldRemoveSourceBranch       string     `json:"should_remove_source_branch,omitempty"`
	ForceRemoveSourceBranch        string     `json:"force_remove_source_branch,omitempty"`
	PreparedAt                     string     `json:"prepared_at,omitempty"`
	AllowCollaboration             string     `json:"allow_collaboration,omitempty"`
	AllowMaintainerToPush          string     `json:"allow_maintainer_to_push,omitempty"`
	Reference                      string     `json:"reference,omitempty"`
	WebURL                         string     `json:"web_url,omitempty"`
	Squash                         string     `json:"squash,omitempty"`
	SquashOnMerge                  string     `json:"squash_on_merge,omitempty"`
	TaskCompletionStatus           string     `json:"task_completion_status,omitempty"`
	HasConflicts                   string     `json:"has_conflicts,omitempty"`
	BlockingDiscussionsResolved    string     `json:"blocking_discussions_resolved,omitempty"`
	ApprovalsBeforeMerge           string     `json:"approvals_before_merge,omitempty"`
}

// EventsQueryParams represents the query parameters for the events endpoint
type EventsQueryParams struct {
	Page       int
	PerPage    int
	Action     string
	TargetType string
	Before     *time.Time
	After      *time.Time
	Sort       string
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path parameter using regex
	// Expected path: /api/v4/users/{id}/events
	re := regexp.MustCompile(`/api/v4/users/([^/]+)/events`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userID := matches[1]

	// Parse query parameters
	params, err := h.parseEventsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock events based on the API specification
	// In a real implementation, this would call h.service.events.GetUserEvents(userID, params)
	events := h.generateMockEvents(userID, params)

	w.Header().Set("Content-Type", "application/json")
	// parsePageCursor() requires these http headers
	w.Header().Set("X-Page", "1")
	w.Header().Set("X-Next-Page", "2")
	w.Header().Set("X-Prev-Page", "0")
	w.Header().Set("X-Total-Pages", "1")
	w.Header().Set("X-Per-Page", "20")
	w.Header().Set("X-Total", "3")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// parseEventsQueryParams parses the query parameters for the events endpoint
func (h *Handler) parseEventsQueryParams(r *http.Request) (*EventsQueryParams, error) {
	params := &EventsQueryParams{
		Page:    1,      // default
		PerPage: 20,     // default
		Sort:    "desc", // default
	}

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse per_page parameter
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			params.PerPage = perPage
		}
	}

	// Parse action parameter
	params.Action = r.URL.Query().Get("action")

	// Parse target_type parameter
	targetType := r.URL.Query().Get("target_type")
	validTargetTypes := map[string]bool{
		"issue": true, "milestone": true, "merge_request": true, "note": true,
		"project": true, "snippet": true, "user": true, "wiki": true, "design": true,
	}
	if targetType != "" && validTargetTypes[targetType] {
		params.TargetType = targetType
	}

	// Parse before parameter
	if beforeStr := r.URL.Query().Get("before"); beforeStr != "" {
		if before, err := time.Parse("2006-01-02", beforeStr); err == nil {
			params.Before = &before
		}
	}

	// Parse after parameter
	if afterStr := r.URL.Query().Get("after"); afterStr != "" {
		if after, err := time.Parse("2006-01-02", afterStr); err == nil {
			params.After = &after
		}
	}

	// Parse sort parameter
	if sort := r.URL.Query().Get("sort"); sort == "asc" || sort == "desc" {
		params.Sort = sort
	}

	return params, nil
}

// generateMockEvents generates mock events for testing purposes
func (h *Handler) generateMockEvents(userID string, params *EventsQueryParams) []Event {
	// Generate mock events based on the API specification
	events := []Event{
		{
			ID:             1,
			ProjectID:      intPtr(2),
			ActionName:     "closed",
			TargetID:       intPtr(160),
			TargetIID:      intPtr(157),
			TargetType:     stringPtr("Issue"),
			AuthorID:       25,
			TargetTitle:    stringPtr("Public project search field"),
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			AuthorUsername: stringPtr("test_user"),
			Imported:       false,
			ImportedFrom:   "none",
		},
		{
			ID:             2,
			ProjectID:      intPtr(3),
			ActionName:     "opened",
			TargetID:       intPtr(161),
			TargetIID:      intPtr(158),
			TargetType:     stringPtr("MergeRequest"),
			AuthorID:       25,
			TargetTitle:    stringPtr("Add new feature"),
			CreatedAt:      time.Now().Add(-12 * time.Hour),
			AuthorUsername: stringPtr("test_user"),
			Imported:       false,
			ImportedFrom:   "none",
		},
		{
			ID:             3,
			ProjectID:      intPtr(1),
			ActionName:     "pushed",
			TargetID:       nil,
			TargetIID:      nil,
			TargetType:     stringPtr("Project"),
			AuthorID:       25,
			TargetTitle:    stringPtr("main"),
			CreatedAt:      time.Now().Add(-6 * time.Hour),
			AuthorUsername: stringPtr("test_user"),
			Imported:       false,
			ImportedFrom:   "none",
		},
	}

	// Apply filtering based on parameters
	var filteredEvents []Event
	for _, event := range events {
		// Filter by action
		if params.Action != "" && event.ActionName != params.Action {
			continue
		}

		// Filter by target_type
		if params.TargetType != "" && (event.TargetType == nil || *event.TargetType != params.TargetType) {
			continue
		}

		// Filter by before date
		if params.Before != nil && event.CreatedAt.After(*params.Before) {
			continue
		}

		// Filter by after date
		if params.After != nil && event.CreatedAt.Before(*params.After) {
			continue
		}

		filteredEvents = append(filteredEvents, event)
	}

	// Apply pagination
	start := (params.Page - 1) * params.PerPage
	end := start + params.PerPage

	if start >= len(filteredEvents) {
		return []Event{}
	}

	if end > len(filteredEvents) {
		end = len(filteredEvents)
	}

	return filteredEvents[start:end]
}

// Helper functions for pointer creation
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
