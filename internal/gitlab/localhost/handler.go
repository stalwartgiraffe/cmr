package localhost

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/withstack"
)

type Handler struct {
	service  *Service
	requests []MergeRequest
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (*Handler) OnServerError(w http.ResponseWriter, msg string, err error) {
	utils.Redln(msg, "\n", err.Error())
	http.Error(w, msg, http.StatusInternalServerError)
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params, err := h.parseProjectsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock projects based on the API specification
	// In a real implementation, this would call h.service.projects.GetGroupProjects(groupID, params)
	projects := h.generateMockProjects(params)

	setOnePagedHeaders(len(projects), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		h.OnServerError(w, "Failed to encode response", err)
		return
	}
}

func (h *Handler) GetGroupsProjects(w http.ResponseWriter, r *http.Request) {
	params, err := h.parseProjectsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock projects based on the API specification
	// In a real implementation, this would call h.service.projects.GetGroupProjects(groupID, params)
	projects := h.generateMockProjects(params)

	setOnePagedHeaders(len(projects), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		h.OnServerError(w, "Failed to encode response", err)
		return
	}
}

// Extract group ID from path parameter using regex
// Expected path:
// Handle the GitLab API v4 groups endpoints:l/api/v4/groups/{id}/projects
/*
	re := regexp.MustCompile(`/api/v4/groups/([^/]+)`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	projectID := matches[1]
*/

// parseProjectsQueryParams parses the query parameters for the projects endpoint
func (h *Handler) parseProjectsQueryParams(r *http.Request) (*ProjectsQueryParams, error) {
	params := &ProjectsQueryParams{
		PageQueryParams: PageQueryParams{
			Page:    1,  // default
			PerPage: 20, // default
		},
		OrderBy:                  "created_at", // default
		Sort:                     "desc",       // default
		Simple:                   false,        // default
		Owned:                    false,        // default
		Starred:                  false,        // default
		WithIssuesEnabled:        false,        // default
		WithMergeRequestsEnabled: false,        // default
		WithShared:               true,         // default
		IncludeSubgroups:         false,        // default
		IncludeAncestorGroups:    false,        // default
		WithCustomAttributes:     false,        // default
		WithSecurityReports:      false,        // default
	}

	parsePageParams(r.URL.Query(), &params.PageQueryParams)

	// Parse archived parameter
	if archivedStr := r.URL.Query().Get("archived"); archivedStr != "" {
		if archived, err := strconv.ParseBool(archivedStr); err == nil {
			params.Archived = &archived
		}
	}

	// Parse visibility parameter
	if visibility := r.URL.Query().Get("visibility"); visibility != "" {
		validVisibilities := map[string]bool{
			"private": true, "internal": true, "public": true,
		}
		if validVisibilities[visibility] {
			params.Visibility = visibility
		}
	}

	// Parse search parameter
	params.Search = r.URL.Query().Get("search")

	// Parse order_by parameter
	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		validOrderBy := map[string]bool{
			"id": true, "name": true, "path": true, "created_at": true,
			"updated_at": true, "last_activity_at": true, "similarity": true, "star_count": true,
		}
		if validOrderBy[orderBy] {
			params.OrderBy = orderBy
		}
	}

	// Parse sort parameter
	if sort := r.URL.Query().Get("sort"); sort == "asc" || sort == "desc" {
		params.Sort = sort
	}

	// Parse simple parameter
	if simpleStr := r.URL.Query().Get("simple"); simpleStr == "true" {
		params.Simple = true
	}

	// Parse owned parameter
	if ownedStr := r.URL.Query().Get("owned"); ownedStr == "true" {
		params.Owned = true
	}

	// Parse starred parameter
	if starredStr := r.URL.Query().Get("starred"); starredStr == "true" {
		params.Starred = true
	}

	// Parse with_issues_enabled parameter
	if withIssuesStr := r.URL.Query().Get("with_issues_enabled"); withIssuesStr == "true" {
		params.WithIssuesEnabled = true
	}

	// Parse with_merge_requests_enabled parameter
	if withMergeRequestsStr := r.URL.Query().Get("with_merge_requests_enabled"); withMergeRequestsStr == "true" {
		params.WithMergeRequestsEnabled = true
	}

	// Parse with_shared parameter
	if withSharedStr := r.URL.Query().Get("with_shared"); withSharedStr == "false" {
		params.WithShared = false
	}

	// Parse include_subgroups parameter
	if includeSubgroupsStr := r.URL.Query().Get("include_subgroups"); includeSubgroupsStr == "true" {
		params.IncludeSubgroups = true
	}

	// Parse include_ancestor_groups parameter
	if includeAncestorGroupsStr := r.URL.Query().Get("include_ancestor_groups"); includeAncestorGroupsStr == "true" {
		params.IncludeAncestorGroups = true
	}

	// Parse min_access_level parameter
	if minAccessLevelStr := r.URL.Query().Get("min_access_level"); minAccessLevelStr != "" {
		if minAccessLevel, err := strconv.Atoi(minAccessLevelStr); err == nil {
			validAccessLevels := map[int]bool{
				10: true, 15: true, 20: true, 30: true, 40: true, 50: true,
			}
			if validAccessLevels[minAccessLevel] {
				params.MinAccessLevel = &minAccessLevel
			}
		}
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

	// Parse with_custom_attributes parameter
	if withCustomAttributesStr := r.URL.Query().Get("with_custom_attributes"); withCustomAttributesStr == "true" {
		params.WithCustomAttributes = true
	}

	// Parse with_security_reports parameter
	if withSecurityReportsStr := r.URL.Query().Get("with_security_reports"); withSecurityReportsStr == "true" {
		params.WithSecurityReports = true
	}

	return params, nil
}

// generateMockProjects generates mock projects for testing purposes
func (h *Handler) generateMockProjects(params *ProjectsQueryParams) []Project {
	// Generate mock projects based on the API specification

	projects := []Project{
		{
			ID:                       25,
			Name:                     "gitlab-foss",
			NameWithNamespace:        "GitLab.org / gitlab-foss",
			Path:                     "gitlab-foss",
			PathWithNamespace:        "gitlab-org/gitlab-foss",
			Description:              "GitLab Community Edition",
			CreatedAt:                time.Now().Add(-365 * 24 * time.Hour),
			UpdatedAt:                time.Now().Add(-24 * time.Hour),
			LastActivityAt:           time.Now().Add(-12 * time.Hour),
			DefaultBranch:            "main",
			TagList:                  []string{"ruby", "rails", "git"},
			Topics:                   []string{"git", "version-control", "collaboration"},
			SSHURLToRepo:             "git@gitlab.example.com:gitlab-org/gitlab-foss.git",
			HTTPURLToRepo:            "https://gitlab.example.com/gitlab-org/gitlab-foss.git",
			WebURL:                   "https://gitlab.example.com/gitlab-org/gitlab-foss",
			ReadmeURL:                "https://gitlab.example.com/gitlab-org/gitlab-foss/-/blob/main/README.md",
			AvatarURL:                "https://gitlab.example.com/uploads/project/avatar/25/gitlab_logo.png",
			StarCount:                2345,
			ForksCount:               589,
			Visibility:               "public",
			IssuesEnabled:            true,
			MergeRequestsEnabled:     true,
			WikiEnabled:              true,
			JobsEnabled:              true,
			SnippetsEnabled:          true,
			ContainerRegistryEnabled: true,
			EmptyRepo:                false,
			Archived:                 false,
			Owner: &UserBasic{
				ID:       1,
				Username: "root",
				Name:     "Administrator",
				State:    "active",
			},
		},
		{
			ID:                       30,
			Name:                     "awesome-project",
			NameWithNamespace:        "GitLab.org / awesome-project",
			Path:                     "awesome-project",
			PathWithNamespace:        "gitlab-org/awesome-project",
			Description:              "An awesome project for demonstration",
			CreatedAt:                time.Now().Add(-180 * 24 * time.Hour),
			UpdatedAt:                time.Now().Add(-48 * time.Hour),
			LastActivityAt:           time.Now().Add(-24 * time.Hour),
			DefaultBranch:            "develop",
			TagList:                  []string{"javascript", "nodejs", "react"},
			Topics:                   []string{"frontend", "web", "javascript"},
			SSHURLToRepo:             "git@gitlab.example.com:gitlab-org/awesome-project.git",
			HTTPURLToRepo:            "https://gitlab.example.com/gitlab-org/awesome-project.git",
			WebURL:                   "https://gitlab.example.com/gitlab-org/awesome-project",
			ReadmeURL:                "https://gitlab.example.com/gitlab-org/awesome-project/-/blob/develop/README.md",
			StarCount:                125,
			ForksCount:               34,
			Visibility:               "internal",
			IssuesEnabled:            true,
			MergeRequestsEnabled:     true,
			WikiEnabled:              false,
			JobsEnabled:              true,
			SnippetsEnabled:          false,
			ContainerRegistryEnabled: false,
			EmptyRepo:                false,
			Archived:                 false,
			Owner: &UserBasic{
				ID:       2,
				Username: "developer",
				Name:     "Developer User",
				State:    "active",
			},
		},
		{
			ID:                       45,
			Name:                     "archived-legacy",
			NameWithNamespace:        "GitLab.org / archived-legacy",
			Path:                     "archived-legacy",
			PathWithNamespace:        "gitlab-org/archived-legacy",
			Description:              "Legacy project that has been archived",
			CreatedAt:                time.Now().Add(-730 * 24 * time.Hour),
			UpdatedAt:                time.Now().Add(-365 * 24 * time.Hour),
			LastActivityAt:           time.Now().Add(-365 * 24 * time.Hour),
			DefaultBranch:            "master",
			TagList:                  []string{"legacy", "deprecated"},
			SSHURLToRepo:             "git@gitlab.example.com:gitlab-org/archived-legacy.git",
			HTTPURLToRepo:            "https://gitlab.example.com/gitlab-org/archived-legacy.git",
			WebURL:                   "https://gitlab.example.com/gitlab-org/archived-legacy",
			StarCount:                5,
			ForksCount:               1,
			Visibility:               "private",
			IssuesEnabled:            false,
			MergeRequestsEnabled:     false,
			WikiEnabled:              false,
			JobsEnabled:              false,
			SnippetsEnabled:          false,
			ContainerRegistryEnabled: false,
			EmptyRepo:                false,
			Archived:                 true,
			Owner: &UserBasic{
				ID:       1,
				Username: "root",
				Name:     "Administrator",
				State:    "active",
			},
		},
	}

	// Apply filtering based on parameters
	var filteredProjects []Project
	for _, project := range projects {
		// Filter by archived status
		if params.Archived != nil && project.Archived != *params.Archived {
			continue
		}

		// Filter by visibility
		if params.Visibility != "" && project.Visibility != params.Visibility {
			continue
		}

		// Filter by search in name and description
		if params.Search != "" {
			searchLower := strings.ToLower(params.Search)
			nameMatch := strings.Contains(strings.ToLower(project.Name), searchLower)
			descMatch := strings.Contains(strings.ToLower(project.Description), searchLower)
			if !nameMatch && !descMatch {
				continue
			}
		}

		// Filter by issues enabled
		if params.WithIssuesEnabled && !project.IssuesEnabled {
			continue
		}

		// Filter by merge requests enabled
		if params.WithMergeRequestsEnabled && !project.MergeRequestsEnabled {
			continue
		}

		// Apply simple mode (return only basic fields)
		if params.Simple {
			filteredProject := Project{
				ID:         project.ID,
				Name:       project.Name,
				Path:       project.Path,
				WebURL:     project.WebURL,
				Visibility: project.Visibility, // Include visibility for filtering validation
			}
			filteredProjects = append(filteredProjects, filteredProject)
		} else {
			filteredProjects = append(filteredProjects, project)
		}
	}

	// Apply pagination
	start := (params.Page - 1) * params.PerPage
	end := start + params.PerPage

	if start >= len(filteredProjects) {
		return []Project{}
	}

	if end > len(filteredProjects) {
		end = len(filteredProjects)
	}

	return filteredProjects[start:end]
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (h *Handler) GetGroupsMergeRequests(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params, err := h.parseMergeRequestsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock merge requests based on the API specification
	// In a real implementation, this would call h.service.requests.GetGroupMergeRequests(groupID, params)
	requests, err := h.generateMockMergeRequests(params)
	if err != nil {
		h.OnServerError(w, "Failed to generate response", err)
		return
	}

	setOnePagedHeaders(len(requests), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(requests); err != nil {
		h.OnServerError(w, "Failed to encode response", err)
		return
	}
}

// parseMergeRequestsQueryParams parses the query parameters for the merge requests endpoint
func (h *Handler) parseMergeRequestsQueryParams(r *http.Request) (*MergeRequestsQueryParams, error) {
	params := &MergeRequestsQueryParams{
		PageQueryParams: PageQueryParams{
			Page:    1,  // default
			PerPage: 20, // default
		},
		State:   "all",        // default
		OrderBy: "created_at", // default
		Sort:    "desc",       // default
	}

	parsePageParams(r.URL.Query(), &params.PageQueryParams)

	// Parse author_id parameter
	if authorIDStr := r.URL.Query().Get("author_id"); authorIDStr != "" {
		if authorID, err := strconv.Atoi(authorIDStr); err == nil {
			params.AuthorID = &authorID
		}
	}

	// Parse author_username parameter
	params.AuthorUsername = r.URL.Query().Get("author_username")

	// Parse assignee_id parameter
	if assigneeIDStr := r.URL.Query().Get("assignee_id"); assigneeIDStr != "" {
		if assigneeID, err := strconv.Atoi(assigneeIDStr); err == nil {
			params.AssigneeID = &assigneeID
		}
	}

	// Parse assignee_username parameter (array)
	if assigneeUsernames := r.URL.Query()["assignee_username"]; len(assigneeUsernames) > 0 {
		params.AssigneeUsername = assigneeUsernames
	}

	// Parse reviewer_username parameter
	params.ReviewerUsername = r.URL.Query().Get("reviewer_username")

	// Parse reviewer_id parameter
	if reviewerIDStr := r.URL.Query().Get("reviewer_id"); reviewerIDStr != "" {
		if reviewerID, err := strconv.Atoi(reviewerIDStr); err == nil {
			params.ReviewerID = &reviewerID
		}
	}

	// Parse labels parameter (array)
	if labels := r.URL.Query()["labels"]; len(labels) > 0 {
		params.Labels = labels
	}

	// Parse milestone parameter
	params.Milestone = r.URL.Query().Get("milestone")

	// Parse my_reaction_emoji parameter
	params.MyReactionEmoji = r.URL.Query().Get("my_reaction_emoji")

	// Parse state parameter
	if state := r.URL.Query().Get("state"); state != "" {
		validStates := map[string]bool{
			"opened": true, "closed": true, "locked": true, "merged": true, "all": true,
		}
		if validStates[state] {
			params.State = state
		}
	}

	// Parse order_by parameter
	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		validOrderBy := map[string]bool{
			"created_at": true, "label_priority": true, "milestone_due": true,
			"popularity": true, "priority": true, "title": true, "updated_at": true, "merged_at": true,
		}
		if validOrderBy[orderBy] {
			params.OrderBy = orderBy
		}
	}

	// Parse sort parameter
	if sort := r.URL.Query().Get("sort"); sort == "asc" || sort == "desc" {
		params.Sort = sort
	}

	// Parse boolean parameters
	if withLabelsDetails := r.URL.Query().Get("with_labels_details"); withLabelsDetails == "true" {
		params.WithLabelsDetails = true
	}

	if withMergeStatusRecheck := r.URL.Query().Get("with_merge_status_recheck"); withMergeStatusRecheck == "true" {
		params.WithMergeStatusRecheck = true
	}

	// Parse datetime parameters
	if createdAfterStr := r.URL.Query().Get("created_after"); createdAfterStr != "" {
		if createdAfter, err := time.Parse(time.RFC3339, createdAfterStr); err == nil {
			params.CreatedAfter = &createdAfter
		}
	}

	if createdBeforeStr := r.URL.Query().Get("created_before"); createdBeforeStr != "" {
		if createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr); err == nil {
			params.CreatedBefore = &createdBefore
		}
	}

	if updatedAfterStr := r.URL.Query().Get("updated_after"); updatedAfterStr != "" {
		if updatedAfter, err := time.Parse(time.RFC3339, updatedAfterStr); err == nil {
			params.UpdatedAfter = &updatedAfter
		}
	}

	if updatedBeforeStr := r.URL.Query().Get("updated_before"); updatedBeforeStr != "" {
		if updatedBefore, err := time.Parse(time.RFC3339, updatedBeforeStr); err == nil {
			params.UpdatedBefore = &updatedBefore
		}
	}

	// Parse view parameter
	if view := r.URL.Query().Get("view"); view == "simple" {
		params.View = view
	}

	// Parse scope parameter
	if scope := r.URL.Query().Get("scope"); scope != "" {
		validScopes := map[string]bool{
			"created-by-me": true, "assigned-to-me": true, "created_by_me": true,
			"assigned_to_me": true, "reviews_for_me": true, "all": true,
		}
		if validScopes[scope] {
			params.Scope = scope
		}
	}

	// Parse other string parameters
	params.SourceBranch = r.URL.Query().Get("source_branch")
	params.TargetBranch = r.URL.Query().Get("target_branch")
	params.Search = r.URL.Query().Get("search")
	params.In = r.URL.Query().Get("in")

	// Parse source_project_id parameter
	if sourceProjectIDStr := r.URL.Query().Get("source_project_id"); sourceProjectIDStr != "" {
		if sourceProjectID, err := strconv.Atoi(sourceProjectIDStr); err == nil {
			params.SourceProjectID = &sourceProjectID
		}
	}

	// Parse wip parameter
	if wip := r.URL.Query().Get("wip"); wip == "yes" || wip == "no" {
		params.WIP = wip
	}

	// Parse negated parameters
	if notAuthorIDStr := r.URL.Query().Get("not[author_id]"); notAuthorIDStr != "" {
		if notAuthorID, err := strconv.Atoi(notAuthorIDStr); err == nil {
			params.NotAuthorID = &notAuthorID
		}
	}

	params.NotAuthorUsername = r.URL.Query().Get("not[author_username]")

	if notAssigneeIDStr := r.URL.Query().Get("not[assignee_id]"); notAssigneeIDStr != "" {
		if notAssigneeID, err := strconv.Atoi(notAssigneeIDStr); err == nil {
			params.NotAssigneeID = &notAssigneeID
		}
	}

	if notAssigneeUsernames := r.URL.Query()["not[assignee_username]"]; len(notAssigneeUsernames) > 0 {
		params.NotAssigneeUsername = notAssigneeUsernames
	}

	return params, nil
}

// generateMockMergeRequests generates mock merge requests for testing purposes
func (h *Handler) generateMockMergeRequests(params *MergeRequestsQueryParams) ([]MergeRequest, error) {
	if h.requests == nil {
		requests := []MergeRequest{}
		for range 75 {
			request := new(MergeRequest)
			err := gofakeit.Struct(request)
			if err != nil {
				return nil, withstack.Errorf("Unmarshal error:%w", err)
			}
			requests = append(requests, *request)
		}

		h.requests = requests
	}

	// Apply filtering based on parameters
	matches := make([]MergeRequest, 0, len(h.requests))
	for _, mr := range h.requests {
		// Filter by state
		if params.State != "all" && mr.State != params.State {
			continue
		}

		// Filter by author_id
		if params.AuthorID != nil && mr.Author != nil && mr.Author.ID != *params.AuthorID {
			continue
		}

		// Filter by author_username
		if params.AuthorUsername != "" && mr.Author != nil && mr.Author.Username != params.AuthorUsername {
			continue
		}

		// Filter by source_branch
		if params.SourceBranch != "" && mr.SourceBranch != params.SourceBranch {
			continue
		}

		// Filter by target_branch
		if params.TargetBranch != "" && mr.TargetBranch != params.TargetBranch {
			continue
		}

		// Filter by created_after
		if params.CreatedAfter != nil && mr.CreatedAt.Before(*params.CreatedAfter) {
			continue
		}

		// Filter by created_before
		if params.CreatedBefore != nil && mr.CreatedAt.After(*params.CreatedBefore) {
			continue
		}

		// Filter by updated_after
		if params.UpdatedAfter != nil && mr.UpdatedAt.Before(*params.UpdatedAfter) {
			continue
		}

		// Filter by updated_before
		if params.UpdatedBefore != nil && mr.UpdatedAt.After(*params.UpdatedBefore) {
			continue
		}

		// Filter by search in title and description
		if params.Search != "" {
			searchLower := strings.ToLower(params.Search)
			titleMatch := strings.Contains(strings.ToLower(mr.Title), searchLower)
			descMatch := strings.Contains(strings.ToLower(mr.Description), searchLower)

			if params.In == "title" {
				if !titleMatch {
					continue
				}
			} else if params.In == "description" {
				if !descMatch {
					continue
				}
			} else {
				// Default: search in both title and description
				if !titleMatch && !descMatch {
					continue
				}
			}
		}

		matches = append(matches, mr)
	}

	return matches, nil
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params, err := h.parseEventsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock events based on the API specification
	// In a real implementation, this would call h.service.events.GetUserEvents(userID, params)
	events, err := h.generateMockEvents(params)
	if err != nil {
		h.OnServerError(w, "Failed to generate response", err)
		return
	}

	setOnePagedHeaders(len(events), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(events); err != nil {
		h.OnServerError(w, "Failed to encode response", err)
		return
	}
}

// setOnePagedHeaders returns just a single page, ignoring query params
//
// we are assume all the items fit in a single page
// but is actually controlled by Page and PerPage query param
// so this actually not returning fully valid paging behavor
// TODO - implement fully paged cursor through the service
func setOnePagedHeaders(total int, params PageQueryParams, header http.Header) {
	header.Set("Content-Type", "application/json")
	// force a single page
	header.Set("X-Per-Page", strconv.Itoa(params.PerPage))

	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))
	prevPage := max(0, params.Page+1)
	nextPage := min(params.Page+1, totalPages)

	header.Set("X-Page", strconv.Itoa(params.Page))
	header.Set("X-Prev-Page", strconv.Itoa(prevPage))
	header.Set("X-Next-Page", strconv.Itoa(nextPage))
	header.Set("X-Total-Pages", strconv.Itoa(totalPages))
	header.Set("X-Total", strconv.Itoa(total))
}

// parseEventsQueryParams parses the query parameters for the events endpoint
func (h *Handler) parseEventsQueryParams(r *http.Request) (*EventsQueryParams, error) {
	params := &EventsQueryParams{
		PageQueryParams: PageQueryParams{
			Page:    1,  // default
			PerPage: 20, // default
		},
		Sort: "desc", // default
	}

	parsePageParams(r.URL.Query(), &params.PageQueryParams)

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

func parsePageParams(query url.Values, params *PageQueryParams) {
	// Parse page parameter
	//description: Current page number
	if pageStr := query.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse per_page parameter
	// description: Number of items per page
	if perPageStr := query.Get("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			params.PerPage = perPage
		}
	}
}

// generateMockEvents generates mock events for testing purposes
func (h *Handler) generateMockEvents(params *EventsQueryParams) ([]Event, error) {
	// Generate mock events based on the API specification
	events := []Event{}
	for range 20 {
		event := new(Event)
		err := gofakeit.Struct(event)
		if err != nil {
			return nil, withstack.Errorf("Unmarshal error:%w", err)
		}
		events = append(events, *event)
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
		return []Event{}, nil
	}

	if end > len(filteredEvents) {
		end = len(filteredEvents)
	}

	return filteredEvents[start:end], nil
}

// Helper functions for pointer creation
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
