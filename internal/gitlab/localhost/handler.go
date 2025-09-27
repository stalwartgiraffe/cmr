package localhost

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ProjectLinks represents the _links field in Project
type ProjectLinks struct {
	Self          string `json:"self,omitempty"`
	Issues        string `json:"issues,omitempty"`
	MergeRequests string `json:"merge_requests,omitempty"`
	RepoBranches  string `json:"repo_branches,omitempty"`
	Labels        string `json:"labels,omitempty"`
	Events        string `json:"events,omitempty"`
	Members       string `json:"members,omitempty"`
	ClusterAgents string `json:"cluster_agents,omitempty"`
}

// Project represents a GitLab project as per API_Entities_Project
type Project struct {
	ID                                        int           `json:"id"`
	Description                               string        `json:"description,omitempty"`
	Name                                      string        `json:"name"`
	NameWithNamespace                         string        `json:"name_with_namespace"`
	Path                                      string        `json:"path"`
	PathWithNamespace                         string        `json:"path_with_namespace"`
	CreatedAt                                 time.Time     `json:"created_at"`
	DefaultBranch                             string        `json:"default_branch,omitempty"`
	TagList                                   []string      `json:"tag_list,omitempty"`
	Topics                                    []string      `json:"topics,omitempty"`
	SSHURLToRepo                              string        `json:"ssh_url_to_repo,omitempty"`
	HTTPURLToRepo                             string        `json:"http_url_to_repo,omitempty"`
	WebURL                                    string        `json:"web_url,omitempty"`
	ReadmeURL                                 string        `json:"readme_url,omitempty"`
	ForksCount                                int           `json:"forks_count"`
	LicenseURL                                string        `json:"license_url,omitempty"`
	AvatarURL                                 string        `json:"avatar_url,omitempty"`
	StarCount                                 int           `json:"star_count"`
	LastActivityAt                            time.Time     `json:"last_activity_at"`
	Visibility                                string        `json:"visibility"`
	RepositoryStorage                         string        `json:"repository_storage,omitempty"`
	ContainerRegistryImagePrefix              string        `json:"container_registry_image_prefix,omitempty"`
	Links                                     *ProjectLinks `json:"_links,omitempty"`
	MarkedForDeletionAt                       *string       `json:"marked_for_deletion_at,omitempty"`
	MarkedForDeletionOn                       *string       `json:"marked_for_deletion_on,omitempty"`
	PackagesEnabled                           bool          `json:"packages_enabled"`
	EmptyRepo                                 bool          `json:"empty_repo"`
	Archived                                  bool          `json:"archived"`
	Owner                                     *UserBasic    `json:"owner,omitempty"`
	ResolveOutdatedDiffDiscussions            bool          `json:"resolve_outdated_diff_discussions"`
	RepositoryObjectFormat                    string        `json:"repository_object_format,omitempty"`
	IssuesEnabled                             bool          `json:"issues_enabled"`
	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled"`
	WikiEnabled                               bool          `json:"wiki_enabled"`
	JobsEnabled                               bool          `json:"jobs_enabled"`
	SnippetsEnabled                           bool          `json:"snippets_enabled"`
	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled"`
	ServiceDeskEnabled                        bool          `json:"service_desk_enabled"`
	ServiceDeskAddress                        string        `json:"service_desk_address,omitempty"`
	CanCreateMergeRequestIn                   bool          `json:"can_create_merge_request_in"`
	IssuesAccessLevel                         string        `json:"issues_access_level,omitempty"`
	RepositoryAccessLevel                     string        `json:"repository_access_level,omitempty"`
	MergeRequestsAccessLevel                  string        `json:"merge_requests_access_level,omitempty"`
	ForkingAccessLevel                        string        `json:"forking_access_level,omitempty"`
	WikiAccessLevel                           string        `json:"wiki_access_level,omitempty"`
	BuildsAccessLevel                         string        `json:"builds_access_level,omitempty"`
	SnippetsAccessLevel                       string        `json:"snippets_access_level,omitempty"`
	PagesAccessLevel                          string        `json:"pages_access_level,omitempty"`
	AnalyticsAccessLevel                      string        `json:"analytics_access_level,omitempty"`
	ContainerRegistryAccessLevel              string        `json:"container_registry_access_level,omitempty"`
	SecurityAndComplianceAccessLevel          string        `json:"security_and_compliance_access_level,omitempty"`
	ReleasesAccessLevel                       string        `json:"releases_access_level,omitempty"`
	EnvironmentsAccessLevel                   string        `json:"environments_access_level,omitempty"`
	FeatureFlagsAccessLevel                   string        `json:"feature_flags_access_level,omitempty"`
	InfrastructureAccessLevel                 string        `json:"infrastructure_access_level,omitempty"`
	MonitorAccessLevel                        string        `json:"monitor_access_level,omitempty"`
	ModelExperimentsAccessLevel               string        `json:"model_experiments_access_level,omitempty"`
	ModelRegistryAccessLevel                  string        `json:"model_registry_access_level,omitempty"`
	PackageRegistryAccessLevel                string        `json:"package_registry_access_level,omitempty"`
	EmailsDisabled                            bool          `json:"emails_disabled"`
	EmailsEnabled                             bool          `json:"emails_enabled"`
	ShowDiffPreviewInEmail                    bool          `json:"show_diff_preview_in_email"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	LFSEnabled                                bool          `json:"lfs_enabled"`
	CreatorID                                 int           `json:"creator_id"`
	MRDefaultTargetSelf                       bool          `json:"mr_default_target_self"`
	ImportURL                                 string        `json:"import_url,omitempty"`
	ImportType                                string        `json:"import_type,omitempty"`
	ImportStatus                              string        `json:"import_status,omitempty"`
	ImportError                               string        `json:"import_error,omitempty"`
	OpenIssuesCount                           int           `json:"open_issues_count"`
	DescriptionHTML                           string        `json:"description_html,omitempty"`
	UpdatedAt                                 time.Time     `json:"updated_at"`
	CIDefaultGitDepth                         int           `json:"ci_default_git_depth,omitempty"`
	CIDeletePipelinesInSeconds                int           `json:"ci_delete_pipelines_in_seconds,omitempty"`
	CIForwardDeploymentEnabled                bool          `json:"ci_forward_deployment_enabled"`
	CIForwardDeploymentRollbackAllowed        bool          `json:"ci_forward_deployment_rollback_allowed"`
	CIJobTokenScopeEnabled                    bool          `json:"ci_job_token_scope_enabled"`
	CISeparatedCaches                         bool          `json:"ci_separated_caches"`
	CIAllowForkPipelinesToRunInParentProject  bool          `json:"ci_allow_fork_pipelines_to_run_in_parent_project"`
	CIIDTokenSubClaimComponents               []string      `json:"ci_id_token_sub_claim_components,omitempty"`
	BuildGitStrategy                          string        `json:"build_git_strategy,omitempty"`
	KeepLatestArtifact                        bool          `json:"keep_latest_artifact"`
	RestrictUserDefinedVariables              bool          `json:"restrict_user_defined_variables"`
	CIPipelineVariablesMinimumOverrideRole    string        `json:"ci_pipeline_variables_minimum_override_role,omitempty"`
	RunnerTokenExpirationInterval             int           `json:"runner_token_expiration_interval,omitempty"`
	GroupRunnersEnabled                       bool          `json:"group_runners_enabled"`
	AutoCancelPendingPipelines                string        `json:"auto_cancel_pending_pipelines,omitempty"`
	BuildTimeout                              int           `json:"build_timeout,omitempty"`
	AutoDevOpsEnabled                         bool          `json:"auto_devops_enabled"`
	AutoDevOpsDeployStrategy                  string        `json:"auto_devops_deploy_strategy,omitempty"`
	CIPushRepositoryForJobTokenAllowed        bool          `json:"ci_push_repository_for_job_token_allowed"`
	RunnersToken                              string        `json:"runners_token,omitempty"`
	CIConfigPath                              string        `json:"ci_config_path,omitempty"`
	PublicJobs                                bool          `json:"public_jobs"`
	SharedWithGroups                          []string      `json:"shared_with_groups,omitempty"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	AllowMergeOnSkippedPipeline               bool          `json:"allow_merge_on_skipped_pipeline"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge"`
	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
	MergeMethod                               string        `json:"merge_method,omitempty"`
	MergeRequestTitleRegex                    string        `json:"merge_request_title_regex,omitempty"`
	MergeRequestTitleRegexDescription         string        `json:"merge_request_title_regex_description,omitempty"`
	SquashOption                              string        `json:"squash_option,omitempty"`
	EnforceAuthChecksOnUploads                bool          `json:"enforce_auth_checks_on_uploads"`
	SuggestionCommitMessage                   string        `json:"suggestion_commit_message,omitempty"`
	MergeCommitTemplate                       string        `json:"merge_commit_template,omitempty"`
	SquashCommitTemplate                      string        `json:"squash_commit_template,omitempty"`
	IssueBranchTemplate                       string        `json:"issue_branch_template,omitempty"`
	WarnAboutPotentiallyUnwantedCharacters    bool          `json:"warn_about_potentially_unwanted_characters"`
	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues"`
	MaxArtifactsSize                          int           `json:"max_artifacts_size,omitempty"`
	ApprovalsBeforeMerge                      string        `json:"approvals_before_merge,omitempty"`
	Mirror                                    string        `json:"mirror,omitempty"`
	MirrorUserID                              string        `json:"mirror_user_id,omitempty"`
	MirrorTriggerBuilds                       string        `json:"mirror_trigger_builds,omitempty"`
	OnlyMirrorProtectedBranches               string        `json:"only_mirror_protected_branches,omitempty"`
	MirrorOverwritesDivergedBranches          string        `json:"mirror_overwrites_diverged_branches,omitempty"`
	ExternalAuthorizationClassificationLabel  string        `json:"external_authorization_classification_label,omitempty"`
	RequirementsEnabled                       string        `json:"requirements_enabled,omitempty"`
	RequirementsAccessLevel                   string        `json:"requirements_access_level,omitempty"`
	SecurityAndComplianceEnabled              string        `json:"security_and_compliance_enabled,omitempty"`
	SecretPushProtectionEnabled               bool          `json:"secret_push_protection_enabled"`
	PreReceiveSecretDetectionEnabled          bool          `json:"pre_receive_secret_detection_enabled"`
	ComplianceFrameworks                      string        `json:"compliance_frameworks,omitempty"`
	IssuesTemplate                            string        `json:"issues_template,omitempty"`
	MergeRequestsTemplate                     string        `json:"merge_requests_template,omitempty"`
	CIRestrictPipelineCancellationRole        string        `json:"ci_restrict_pipeline_cancellation_role,omitempty"`
	MergePipelinesEnabled                     string        `json:"merge_pipelines_enabled,omitempty"`
	MergeTrainsEnabled                        string        `json:"merge_trains_enabled,omitempty"`
	MergeTrainsSkipTrainAllowed               string        `json:"merge_trains_skip_train_allowed,omitempty"`
	OnlyAllowMergeIfAllStatusChecksPassed     string        `json:"only_allow_merge_if_all_status_checks_passed,omitempty"`
	AllowPipelineTriggerApproveDeployment     bool          `json:"allow_pipeline_trigger_approve_deployment"`
	PreventMergeWithoutJiraIssue              string        `json:"prevent_merge_without_jira_issue,omitempty"`
	AutoDuoCodeReviewEnabled                  string        `json:"auto_duo_code_review_enabled,omitempty"`
	DuoRemoteFlowsEnabled                     string        `json:"duo_remote_flows_enabled,omitempty"`
	WebBasedCommitSigningEnabled              string        `json:"web_based_commit_signing_enabled,omitempty"`
	SPPRepositoryPipelineAccess               bool          `json:"spp_repository_pipeline_access"`
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
	ID          int    `json:"id"`
	Username    string `json:"username"`
	PublicEmail string `json:"public_email,omitempty"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Locked      bool   `json:"locked"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	AvatarPath  string `json:"avatar_path,omitempty"`
}

// MergeRequest represents a GitLab merge request as per API_Entities_MergeRequestBasic
type MergeRequest struct {
	ID                          int        `json:"id"`
	IID                         int        `json:"iid"`
	ProjectID                   int        `json:"project_id"`
	Title                       string     `json:"title"`
	Description                 string     `json:"description,omitempty"`
	State                       string     `json:"state"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
	MergedBy                    *UserBasic `json:"merged_by,omitempty"`
	MergeUser                   *UserBasic `json:"merge_user,omitempty"`
	MergedAt                    *string    `json:"merged_at,omitempty"`
	ClosedBy                    *UserBasic `json:"closed_by,omitempty"`
	ClosedAt                    *string    `json:"closed_at,omitempty"`
	TitleHTML                   string     `json:"title_html,omitempty"`
	DescriptionHTML             string     `json:"description_html,omitempty"`
	TargetBranch                string     `json:"target_branch"`
	SourceBranch                string     `json:"source_branch"`
	UserNotesCount              string     `json:"user_notes_count,omitempty"`
	Upvotes                     string     `json:"upvotes,omitempty"`
	Downvotes                   string     `json:"downvotes,omitempty"`
	Author                      *UserBasic `json:"author,omitempty"`
	Assignees                   *UserBasic `json:"assignees,omitempty"`
	Assignee                    *UserBasic `json:"assignee,omitempty"`
	Reviewers                   *UserBasic `json:"reviewers,omitempty"`
	SourceProjectID             string     `json:"source_project_id,omitempty"`
	TargetProjectID             string     `json:"target_project_id,omitempty"`
	Labels                      string     `json:"labels,omitempty"`
	Draft                       string     `json:"draft,omitempty"`
	Imported                    string     `json:"imported,omitempty"`
	ImportedFrom                string     `json:"imported_from,omitempty"`
	WorkInProgress              string     `json:"work_in_progress,omitempty"`
	MergeWhenPipelineSucceeds   string     `json:"merge_when_pipeline_succeeds,omitempty"`
	MergeStatus                 string     `json:"merge_status,omitempty"`
	DetailedMergeStatus         string     `json:"detailed_merge_status,omitempty"`
	MergeAfter                  string     `json:"merge_after,omitempty"`
	SHA                         string     `json:"sha,omitempty"`
	MergeCommitSHA              string     `json:"merge_commit_sha,omitempty"`
	SquashCommitSHA             string     `json:"squash_commit_sha,omitempty"`
	DiscussionLocked            string     `json:"discussion_locked,omitempty"`
	ShouldRemoveSourceBranch    string     `json:"should_remove_source_branch,omitempty"`
	ForceRemoveSourceBranch     string     `json:"force_remove_source_branch,omitempty"`
	PreparedAt                  string     `json:"prepared_at,omitempty"`
	AllowCollaboration          string     `json:"allow_collaboration,omitempty"`
	AllowMaintainerToPush       string     `json:"allow_maintainer_to_push,omitempty"`
	Reference                   string     `json:"reference,omitempty"`
	WebURL                      string     `json:"web_url,omitempty"`
	Squash                      string     `json:"squash,omitempty"`
	SquashOnMerge               string     `json:"squash_on_merge,omitempty"`
	TaskCompletionStatus        string     `json:"task_completion_status,omitempty"`
	HasConflicts                string     `json:"has_conflicts,omitempty"`
	BlockingDiscussionsResolved string     `json:"blocking_discussions_resolved,omitempty"`
	ApprovalsBeforeMerge        string     `json:"approvals_before_merge,omitempty"`
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

// MergeRequestsQueryParams represents the query parameters for the merge requests endpoint
type MergeRequestsQueryParams struct {
	AuthorID                  *int       `json:"author_id,omitempty"`
	AuthorUsername            string     `json:"author_username,omitempty"`
	AssigneeID                *int       `json:"assignee_id,omitempty"`
	AssigneeUsername          []string   `json:"assignee_username,omitempty"`
	ReviewerUsername          string     `json:"reviewer_username,omitempty"`
	ReviewerID                *int       `json:"reviewer_id,omitempty"`
	Labels                    []string   `json:"labels,omitempty"`
	Milestone                 string     `json:"milestone,omitempty"`
	MyReactionEmoji           string     `json:"my_reaction_emoji,omitempty"`
	State                     string     `json:"state,omitempty"`
	OrderBy                   string     `json:"order_by,omitempty"`
	Sort                      string     `json:"sort,omitempty"`
	WithLabelsDetails         bool       `json:"with_labels_details,omitempty"`
	WithMergeStatusRecheck    bool       `json:"with_merge_status_recheck,omitempty"`
	CreatedAfter              *time.Time `json:"created_after,omitempty"`
	CreatedBefore             *time.Time `json:"created_before,omitempty"`
	UpdatedAfter              *time.Time `json:"updated_after,omitempty"`
	UpdatedBefore             *time.Time `json:"updated_before,omitempty"`
	View                      string     `json:"view,omitempty"`
	Scope                     string     `json:"scope,omitempty"`
	SourceBranch              string     `json:"source_branch,omitempty"`
	SourceProjectID           *int       `json:"source_project_id,omitempty"`
	TargetBranch              string     `json:"target_branch,omitempty"`
	Search                    string     `json:"search,omitempty"`
	In                        string     `json:"in,omitempty"`
	WIP                       string     `json:"wip,omitempty"`
	NotAuthorID               *int       `json:"not_author_id,omitempty"`
	NotAuthorUsername         string     `json:"not_author_username,omitempty"`
	NotAssigneeID             *int       `json:"not_assignee_id,omitempty"`
	NotAssigneeUsername       []string   `json:"not_assignee_username,omitempty"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleMergeRequests(w http.ResponseWriter, r *http.Request) {
	// Extract group ID from path parameter using regex
	// Expected path: /api/v4/groups/{id}/merge_requests
	re := regexp.MustCompile(`/api/v4/groups/([^/]+)/merge_requests`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	groupID := matches[1]

	// Parse query parameters
	params, err := h.parseMergeRequestsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock merge requests based on the API specification
	// In a real implementation, this would call h.service.mergeRequests.GetGroupMergeRequests(groupID, params)
	mergeRequests := h.generateMockMergeRequests(groupID, params)

	w.Header().Set("Content-Type", "application/json")
	// Add pagination headers like GitLab API
	w.Header().Set("X-Page", "1")
	w.Header().Set("X-Next-Page", "2")
	w.Header().Set("X-Prev-Page", "0")
	w.Header().Set("X-Total-Pages", "1")
	w.Header().Set("X-Per-Page", "20")
	w.Header().Set("X-Total", "3")
	if err := json.NewEncoder(w).Encode(mergeRequests); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// parseMergeRequestsQueryParams parses the query parameters for the merge requests endpoint
func (h *Handler) parseMergeRequestsQueryParams(r *http.Request) (*MergeRequestsQueryParams, error) {
	params := &MergeRequestsQueryParams{
		State:   "all",        // default
		OrderBy: "created_at", // default
		Sort:    "desc",       // default
	}

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
func (h *Handler) generateMockMergeRequests(groupID string, params *MergeRequestsQueryParams) []MergeRequest {
	// Generate mock merge requests based on the API specification
	mergeRequests := []MergeRequest{
		{
			ID:          84,
			IID:         14,
			ProjectID:   4,
			Title:       "Impedit et ut et dolores vero provident ullam est",
			Description: "Repellendus impedit et vel velit dignissimos.",
			State:       "opened",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
			Author: &UserBasic{
				ID:       25,
				Username: "test_user",
				Name:     "Test User",
				State:    "active",
			},
			TargetBranch: "main",
			SourceBranch: "feature-branch",
			WebURL:       "https://gitlab.example.com/group/project/-/merge_requests/14",
		},
		{
			ID:          85,
			IID:         15,
			ProjectID:   5,
			Title:       "Add new feature for user management",
			Description: "This MR adds comprehensive user management features.",
			State:       "merged",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
			Author: &UserBasic{
				ID:       26,
				Username: "developer",
				Name:     "Developer User",
				State:    "active",
			},
			MergedBy: &UserBasic{
				ID:       27,
				Username: "maintainer",
				Name:     "Maintainer User",
				State:    "active",
			},
			TargetBranch: "main",
			SourceBranch: "user-management",
			WebURL:       "https://gitlab.example.com/group/project/-/merge_requests/15",
		},
		{
			ID:          86,
			IID:         16,
			ProjectID:   6,
			Title:       "Fix critical security vulnerability",
			Description: "Addresses CVE-2023-1234 security issue.",
			State:       "closed",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
			UpdatedAt:   time.Now().Add(-48 * time.Hour),
			Author: &UserBasic{
				ID:       28,
				Username: "security_expert",
				Name:     "Security Expert",
				State:    "active",
			},
			TargetBranch: "main",
			SourceBranch: "security-fix",
			WebURL:       "https://gitlab.example.com/group/project/-/merge_requests/16",
		},
	}

	// Apply filtering based on parameters
	var filteredMergeRequests []MergeRequest
	for _, mr := range mergeRequests {
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

		filteredMergeRequests = append(filteredMergeRequests, mr)
	}

	return filteredMergeRequests
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
