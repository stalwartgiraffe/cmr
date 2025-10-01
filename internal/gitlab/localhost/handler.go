package localhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/davecgh/go-spew/spew"
)

// ProjectLinks represents the _links field in Project
type ProjectLinks struct {
	Self          string `json:"self,omitempty" fake:"{url}"`
	Issues        string `json:"issues,omitempty" fake:"{url}"`
	MergeRequests string `json:"merge_requests,omitempty" fake:"{url}"`
	RepoBranches  string `json:"repo_branches,omitempty" fake:"{url}"`
	Labels        string `json:"labels,omitempty" fake:"{url}"`
	Events        string `json:"events,omitempty" fake:"{url}"`
	Members       string `json:"members,omitempty" fake:"{url}"`
	ClusterAgents string `json:"cluster_agents,omitempty" fake:"{url}"`
}

// Project represents a GitLab project as per API_Entities_Project
type Project struct {
	ID                                        int           `json:"id" fake:"{number:100,999}"`
	Description                               string        `json:"description,omitempty" fake:"{sentence:10}"`
	Name                                      string        `json:"name" fake:"{name}"`
	NameWithNamespace                         string        `json:"name_with_namespace" fake:"{name}"`
	Path                                      string        `json:"path" fake:"{word}"`
	PathWithNamespace                         string        `json:"path_with_namespace" fake:"{word}"` // TODO extend
	CreatedAt                                 time.Time     `json:"created_at" fake:"{date}"`
	DefaultBranch                             string        `json:"default_branch,omitempty" fake:"{randomstring:[main,master,develop]}"`
	TagList                                   []string      `json:"tag_list,omitempty" fake:"{tags}"`
	Topics                                    []string      `json:"topics,omitempty" fake:"{tags}"`
	SSHURLToRepo                              string        `json:"ssh_url_to_repo,omitempty" fake:"{url}"`
	HTTPURLToRepo                             string        `json:"http_url_to_repo,omitempty" fake:"{url}"`
	WebURL                                    string        `json:"web_url,omitempty" fake:"{url}"`
	ReadmeURL                                 string        `json:"readme_url,omitempty" fake:"{url}"`
	ForksCount                                int           `json:"forks_count" fake:"{number:0,1000}"`
	LicenseURL                                string        `json:"license_url,omitempty" fake:"{url}"`
	AvatarURL                                 string        `json:"avatar_url,omitempty" fake:"{imageurl:200,200}"`
	StarCount                                 int           `json:"star_count" fake:"{number:0,5000}"`
	LastActivityAt                            time.Time     `json:"last_activity_at" fake:"{date}"`
	Visibility                                string        `json:"visibility" fake:"{randomstring:[public,private,internal]}"`
	RepositoryStorage                         string        `json:"repository_storage,omitempty" fake:"{word}"`
	ContainerRegistryImagePrefix              string        `json:"container_registry_image_prefix,omitempty" fake:"{url}"`
	Links                                     *ProjectLinks `json:"_links,omitempty"`
	MarkedForDeletionAt                       *string       `json:"marked_for_deletion_at,omitempty"` // TODO extend
	MarkedForDeletionOn                       *string       `json:"marked_for_deletion_on,omitempty"` // TODO extend
	PackagesEnabled                           bool          `json:"packages_enabled" fake:"{bool}"`
	EmptyRepo                                 bool          `json:"empty_repo" fake:"{bool}"`
	Archived                                  bool          `json:"archived" fake:"{bool}"`
	Owner                                     *UserBasic    `json:"owner,omitempty"`
	ResolveOutdatedDiffDiscussions            bool          `json:"resolve_outdated_diff_discussions" fake:"{bool}"`
	RepositoryObjectFormat                    string        `json:"repository_object_format,omitempty" fake:"{randomstring:[sha1,sha256]}"`
	IssuesEnabled                             bool          `json:"issues_enabled" fake:"{bool}"`
	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled" fake:"{bool}"`
	WikiEnabled                               bool          `json:"wiki_enabled" fake:"{bool}"`
	JobsEnabled                               bool          `json:"jobs_enabled" fake:"{bool}"`
	SnippetsEnabled                           bool          `json:"snippets_enabled" fake:"{bool}"`
	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled" fake:"{bool}"`
	ServiceDeskEnabled                        bool          `json:"service_desk_enabled" fake:"{bool}"`
	ServiceDeskAddress                        string        `json:"service_desk_address,omitempty" fake:"{email}"`
	CanCreateMergeRequestIn                   bool          `json:"can_create_merge_request_in" fake:"{bool}"`
	IssuesAccessLevel                         string        `json:"issues_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	RepositoryAccessLevel                     string        `json:"repository_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	MergeRequestsAccessLevel                  string        `json:"merge_requests_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	ForkingAccessLevel                        string        `json:"forking_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	WikiAccessLevel                           string        `json:"wiki_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	BuildsAccessLevel                         string        `json:"builds_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	SnippetsAccessLevel                       string        `json:"snippets_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	PagesAccessLevel                          string        `json:"pages_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	AnalyticsAccessLevel                      string        `json:"analytics_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	ContainerRegistryAccessLevel              string        `json:"container_registry_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	SecurityAndComplianceAccessLevel          string        `json:"security_and_compliance_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	ReleasesAccessLevel                       string        `json:"releases_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	EnvironmentsAccessLevel                   string        `json:"environments_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	FeatureFlagsAccessLevel                   string        `json:"feature_flags_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	InfrastructureAccessLevel                 string        `json:"infrastructure_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	MonitorAccessLevel                        string        `json:"monitor_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	ModelExperimentsAccessLevel               string        `json:"model_experiments_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	ModelRegistryAccessLevel                  string        `json:"model_registry_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	PackageRegistryAccessLevel                string        `json:"package_registry_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	EmailsDisabled                            bool          `json:"emails_disabled" fake:"{bool}"`
	EmailsEnabled                             bool          `json:"emails_enabled" fake:"{bool}"`
	ShowDiffPreviewInEmail                    bool          `json:"show_diff_preview_in_email" fake:"{bool}"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled" fake:"{bool}"`
	LFSEnabled                                bool          `json:"lfs_enabled" fake:"{bool}"`
	CreatorID                                 int           `json:"creator_id" fake:"{number:1,1000}"`
	MRDefaultTargetSelf                       bool          `json:"mr_default_target_self" fake:"{bool}"`
	ImportURL                                 string        `json:"import_url,omitempty" fake:"{url}"`
	ImportType                                string        `json:"import_type,omitempty" fake:"{word}"`
	ImportStatus                              string        `json:"import_status,omitempty" fake:"{randomstring:[none,scheduled,started,finished,failed]}"`
	ImportError                               string        `json:"import_error,omitempty" fake:"{sentence:5}"`
	OpenIssuesCount                           int           `json:"open_issues_count" fake:"{number:0,100}"`
	DescriptionHTML                           string        `json:"description_html,omitempty"` // TODO extend
	UpdatedAt                                 time.Time     `json:"updated_at" fake:"{date}"`
	CIDefaultGitDepth                         int           `json:"ci_default_git_depth,omitempty" fake:"{number:1,50}"`
	CIDeletePipelinesInSeconds                int           `json:"ci_delete_pipelines_in_seconds,omitempty" fake:"{number:0,604800}"`
	CIForwardDeploymentEnabled                bool          `json:"ci_forward_deployment_enabled" fake:"{bool}"`
	CIForwardDeploymentRollbackAllowed        bool          `json:"ci_forward_deployment_rollback_allowed" fake:"{bool}"`
	CIJobTokenScopeEnabled                    bool          `json:"ci_job_token_scope_enabled" fake:"{bool}"`
	CISeparatedCaches                         bool          `json:"ci_separated_caches" fake:"{bool}"`
	CIAllowForkPipelinesToRunInParentProject  bool          `json:"ci_allow_fork_pipelines_to_run_in_parent_project" fake:"{bool}"`
	CIIDTokenSubClaimComponents               []string      `json:"ci_id_token_sub_claim_components,omitempty" fakesize:"2,6" fake:"{word}"`
	BuildGitStrategy                          string        `json:"build_git_strategy,omitempty" fake:"{randomstring:[fetch,clone]}"`
	KeepLatestArtifact                        bool          `json:"keep_latest_artifact" fake:"{bool}"`
	RestrictUserDefinedVariables              bool          `json:"restrict_user_defined_variables" fake:"{bool}"`
	CIPipelineVariablesMinimumOverrideRole    string        `json:"ci_pipeline_variables_minimum_override_role,omitempty" fake:"{randomstring:[developer,maintainer,owner]}"`
	RunnerTokenExpirationInterval             int           `json:"runner_token_expiration_interval,omitempty" fake:"{number:0,31536000}"`
	GroupRunnersEnabled                       bool          `json:"group_runners_enabled" fake:"{bool}"`
	AutoCancelPendingPipelines                string        `json:"auto_cancel_pending_pipelines,omitempty" fake:"{randomstring:[enabled,disabled]}"`
	BuildTimeout                              int           `json:"build_timeout,omitempty" fake:"{number:0,86400}"`
	AutoDevOpsEnabled                         bool          `json:"auto_devops_enabled" fake:"{bool}"`
	AutoDevOpsDeployStrategy                  string        `json:"auto_devops_deploy_strategy,omitempty" fake:"{randomstring:[continuous,manual,timed_incremental]}"`
	CIPushRepositoryForJobTokenAllowed        bool          `json:"ci_push_repository_for_job_token_allowed" fake:"{bool}"`
	RunnersToken                              string        `json:"runners_token,omitempty" fake:"{uuid}"`
	CIConfigPath                              string        `json:"ci_config_path,omitempty" fake:"{filePath}"` // Does this work?
	PublicJobs                                bool          `json:"public_jobs" fake:"{bool}"`
	SharedWithGroups                          []string      `json:"shared_with_groups,omitempty" fake:"{tags}"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds" fake:"{bool}"`
	AllowMergeOnSkippedPipeline               bool          `json:"allow_merge_on_skipped_pipeline" fake:"{bool}"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled" fake:"{bool}"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved" fake:"{bool}"`
	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge" fake:"{bool}"`
	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled" fake:"{bool}"`
	MergeMethod                               string        `json:"merge_method,omitempty" fake:"{randomstring:[merge,rebase_merge,ff]}"`
	MergeRequestTitleRegex                    string        `json:"merge_request_title_regex,omitempty" fake:"{regex:[A-Z]{3}-[0-9]{3}}"`
	MergeRequestTitleRegexDescription         string        `json:"merge_request_title_regex_description,omitempty" fake:"{sentence:5}"`
	SquashOption                              string        `json:"squash_option,omitempty" fake:"{randomstring:[never,always,default_on,default_off]}"`
	EnforceAuthChecksOnUploads                bool          `json:"enforce_auth_checks_on_uploads" fake:"{bool}"`
	SuggestionCommitMessage                   string        `json:"suggestion_commit_message,omitempty" fake:"{sentence:5}"`
	MergeCommitTemplate                       string        `json:"merge_commit_template,omitempty" fake:"{sentence:5}"`
	SquashCommitTemplate                      string        `json:"squash_commit_template,omitempty" fake:"{sentence:5}"`
	IssueBranchTemplate                       string        `json:"issue_branch_template,omitempty" fake:"{sentence:5}"`
	WarnAboutPotentiallyUnwantedCharacters    bool          `json:"warn_about_potentially_unwanted_characters" fake:"{bool}"`
	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues" fake:"{bool}"`
	MaxArtifactsSize                          int           `json:"max_artifacts_size,omitempty" fake:"{number:100,10000}"`
	ApprovalsBeforeMerge                      string        `json:"approvals_before_merge,omitempty" fake:"{number:0,10}"`
	Mirror                                    string        `json:"mirror,omitempty" fake:"{bool}"`
	MirrorUserID                              string        `json:"mirror_user_id,omitempty" fake:"{number:1,1000}"`
	MirrorTriggerBuilds                       string        `json:"mirror_trigger_builds,omitempty" fake:"{bool}"`
	OnlyMirrorProtectedBranches               string        `json:"only_mirror_protected_branches,omitempty" fake:"{bool}"`
	MirrorOverwritesDivergedBranches          string        `json:"mirror_overwrites_diverged_branches,omitempty" fake:"{bool}"`
	ExternalAuthorizationClassificationLabel  string        `json:"external_authorization_classification_label,omitempty" fake:"{word}"`
	RequirementsEnabled                       string        `json:"requirements_enabled,omitempty" fake:"{bool}"`
	RequirementsAccessLevel                   string        `json:"requirements_access_level,omitempty" fake:"{randomstring:[enabled,disabled,private]}"`
	SecurityAndComplianceEnabled              string        `json:"security_and_compliance_enabled,omitempty" fake:"{bool}"`
	SecretPushProtectionEnabled               bool          `json:"secret_push_protection_enabled" fake:"{bool}"`
	PreReceiveSecretDetectionEnabled          bool          `json:"pre_receive_secret_detection_enabled" fake:"{bool}"`
	ComplianceFrameworks                      string        `json:"compliance_frameworks,omitempty" fake:"{word}"`
	IssuesTemplate                            string        `json:"issues_template,omitempty" fake:"{sentence:5}"`
	MergeRequestsTemplate                     string        `json:"merge_requests_template,omitempty" fake:"{sentence:5}"`
	CIRestrictPipelineCancellationRole        string        `json:"ci_restrict_pipeline_cancellation_role,omitempty" fake:"{randomstring:[developer,maintainer,owner]}"`
	MergePipelinesEnabled                     string        `json:"merge_pipelines_enabled,omitempty" fake:"{bool}"`
	MergeTrainsEnabled                        string        `json:"merge_trains_enabled,omitempty" fake:"{bool}"`
	MergeTrainsSkipTrainAllowed               string        `json:"merge_trains_skip_train_allowed,omitempty" fake:"{bool}"`
	OnlyAllowMergeIfAllStatusChecksPassed     string        `json:"only_allow_merge_if_all_status_checks_passed,omitempty" fake:"{bool}"`
	AllowPipelineTriggerApproveDeployment     bool          `json:"allow_pipeline_trigger_approve_deployment" fake:"{bool}"`
	PreventMergeWithoutJiraIssue              string        `json:"prevent_merge_without_jira_issue,omitempty" fake:"{bool}"`
	AutoDuoCodeReviewEnabled                  string        `json:"auto_duo_code_review_enabled,omitempty" fake:"{bool}"`
	DuoRemoteFlowsEnabled                     string        `json:"duo_remote_flows_enabled,omitempty" fake:"{bool}"`
	WebBasedCommitSigningEnabled              string        `json:"web_based_commit_signing_enabled,omitempty" fake:"{bool}"`
	SPPRepositoryPipelineAccess               bool          `json:"spp_repository_pipeline_access" fake:"{bool}"`
}

// Event represents a GitLab event as per API_Entities_Event
type Event struct {
	ID             int       `json:"id" fake:"{number:100,999}"`
	ProjectID      *int      `json:"project_id,omitempty" fake:"{number:100,999}"`
	ActionName     string    `json:"action_name" fake:"{randomstring:[created,updated,closed,reopened,pushed,commented,merged,joined,left,destroyed,expired]}"`
	TargetID       *int      `json:"target_id,omitempty"  fake:"{number:100,999}"`
	TargetIID      *int      `json:"target_iid,omitempty"   fake:"{number:100,999}"`
	TargetType     *string   `json:"target_type,omitempty" fake:"{randomstring:[issue,milestone,merge_request,note,project,snippet,user,wiki,design]}"`
	AuthorID       int       `json:"author_id" fake:"{number:100,999}"`
	TargetTitle    *string   `json:"target_title,omitempty" fake:"{sentence:10}"`
	CreatedAt      time.Time `json:"created_at" fake:"{date}"`
	AuthorUsername *string   `json:"author_username,omitempty" fake:"{name}"`
	Imported       bool      `json:"imported" fake:"{bool}"`
	ImportedFrom   string    `json:"imported_from" fake:"{randomstring:[none,github,bitbucket,gitlab,gitea]}"`
}

// UserBasic represents a GitLab user as per API_Entities_UserBasic
type UserBasic struct {
	ID          int    `json:"id" fake:"{number:100,999}"`
	Username    string `json:"username" fake:"{username}"`
	PublicEmail string `json:"public_email,omitempty" fake:"{email}"`
	Name        string `json:"name" fake:"{name}"`
	State       string `json:"state" fake:"{randomstring:[active,blocked,deactivated]}"`
	Locked      bool   `json:"locked" fake:"{bool}"`
	AvatarURL   string `json:"avatar_url,omitempty" fake:"{imageurl:200,200}"`
	AvatarPath  string `json:"avatar_path,omitempty" fake:"{filePath}"`
}

// MergeRequest represents a GitLab merge request as per API_Entities_MergeRequestBasic
type MergeRequest struct {
	ID                          int        `json:"id" fake:"{number:1,10000}"`
	IID                         int        `json:"iid" fake:"{number:1,1000}"`
	ProjectID                   int        `json:"project_id" fake:"{number:1,1000}"`
	Title                       string     `json:"title" fake:"{sentence:8}"`
	Description                 string     `json:"description,omitempty" fake:"{paragraph:3,5,10}"`
	State                       string     `json:"state" fake:"{randomstring:[opened,closed,locked,merged]}"`
	CreatedAt                   time.Time  `json:"created_at" fake:"{date}"`
	UpdatedAt                   time.Time  `json:"updated_at" fake:"{date}"`
	MergedBy                    *UserBasic `json:"merged_by,omitempty" fake:"{username}"`
	MergeUser                   *UserBasic `json:"merge_user,omitempty" fake:"{username}"`
	MergedAt                    *string    `json:"merged_at,omitempty" fake:"{username}"`
	ClosedBy                    *UserBasic `json:"closed_by,omitempty" fake:"{username}"`
	ClosedAt                    *string    `json:"closed_at,omitempty" fake:"{username}"`
	TitleHTML                   string     `json:"title_html,omitempty"` // TODO extend
	DescriptionHTML             string     `json:"description_html,omitempty"` // TODO extend
	TargetBranch                string     `json:"target_branch" fake:"{randomstring:[main,master,develop,staging]}"`
	SourceBranch                string     `json:"source_branch" fake:"{word}"`
	UserNotesCount              string     `json:"user_notes_count,omitempty" fake:"{number:0,100}"`
	Upvotes                     string     `json:"upvotes,omitempty" fake:"{number:0,50}"`
	Downvotes                   string     `json:"downvotes,omitempty" fake:"{number:0,10}"`
	Author                      *UserBasic `json:"author,omitempty"`
	Assignees                   *UserBasic `json:"assignees,omitempty"`
	Assignee                    *UserBasic `json:"assignee,omitempty"`
	Reviewers                   *UserBasic `json:"reviewers,omitempty"`
	SourceProjectID             string     `json:"source_project_id,omitempty" fake:"{number:1,1000}"`
	TargetProjectID             string     `json:"target_project_id,omitempty" fake:"{number:1,1000}"`
	Labels                      string     `json:"labels,omitempty" fake:"{word}"`
	Draft                       string     `json:"draft,omitempty" fake:"{bool}"`
	Imported                    string     `json:"imported,omitempty" fake:"{bool}"`
	ImportedFrom                string     `json:"imported_from,omitempty" fake:"{randomstring:[none,github,bitbucket,gitlab]}"`
	WorkInProgress              string     `json:"work_in_progress,omitempty" fake:"{bool}"`
	MergeWhenPipelineSucceeds   string     `json:"merge_when_pipeline_succeeds,omitempty" fake:"{bool}"`
	MergeStatus                 string     `json:"merge_status,omitempty" fake:"{randomstring:[can_be_merged,cannot_be_merged,unchecked]}"`
	DetailedMergeStatus         string     `json:"detailed_merge_status,omitempty" fake:"{randomstring:[mergeable,broken_status,checking,ci_must_pass,ci_still_running,discussions_not_resolved,draft_status]}"`
	MergeAfter                  string     `json:"merge_after,omitempty" fake:"{date}"`
	SHA                         string     `json:"sha,omitempty" fake:"{uuid}"`
	MergeCommitSHA              string     `json:"merge_commit_sha,omitempty" fake:"{uuid}"`
	SquashCommitSHA             string     `json:"squash_commit_sha,omitempty" fake:"{uuid}"`
	DiscussionLocked            string     `json:"discussion_locked,omitempty" fake:"{bool}"`
	ShouldRemoveSourceBranch    string     `json:"should_remove_source_branch,omitempty" fake:"{bool}"`
	ForceRemoveSourceBranch     string     `json:"force_remove_source_branch,omitempty" fake:"{bool}"`
	PreparedAt                  string     `json:"prepared_at,omitempty" fake:"{date}"`
	AllowCollaboration          string     `json:"allow_collaboration,omitempty" fake:"{bool}"`
	AllowMaintainerToPush       string     `json:"allow_maintainer_to_push,omitempty" fake:"{bool}"`
	Reference                   string     `json:"reference,omitempty" fake:"{word}"`
	WebURL                      string     `json:"web_url,omitempty" fake:"{url}"`
	Squash                      string     `json:"squash,omitempty" fake:"{bool}"`
	SquashOnMerge               string     `json:"squash_on_merge,omitempty" fake:"{bool}"`
	TaskCompletionStatus        string     `json:"task_completion_status,omitempty" fake:"{word}"`
	HasConflicts                string     `json:"has_conflicts,omitempty" fake:"{bool}"`
	BlockingDiscussionsResolved string     `json:"blocking_discussions_resolved,omitempty" fake:"{bool}"`
	ApprovalsBeforeMerge        string     `json:"approvals_before_merge,omitempty" fake:"{number:0,10}"`
}

type PageQueryParams struct {
	Page    int `json:"page,omitempty" fake:"{number:1,100}"`
	PerPage int `json:"per_page,omitempty" fake:"{number:10,100}"`
}

// EventsQueryParams represents the query parameters for the events endpoint
type EventsQueryParams struct {
	PageQueryParams

	Action     string
	TargetType string
	Before     *time.Time
	After      *time.Time
	Sort       string
}

// MergeRequestsQueryParams represents the query parameters for the merge requests endpoint
type MergeRequestsQueryParams struct {
	PageQueryParams

	AuthorID               *int       `json:"author_id,omitempty"`
	AuthorUsername         string     `json:"author_username,omitempty" fake:"{username}"`
	AssigneeID             *int       `json:"assignee_id,omitempty"`
	AssigneeUsername       []string   `json:"assignee_username,omitempty"`
	ReviewerUsername       string     `json:"reviewer_username,omitempty" fake:"{username}"`
	ReviewerID             *int       `json:"reviewer_id,omitempty"`
	Labels                 []string   `json:"labels,omitempty"`
	Milestone              string     `json:"milestone,omitempty" fake:"{word}"`
	MyReactionEmoji        string     `json:"my_reaction_emoji,omitempty" fake:"{emoji}"`
	State                  string     `json:"state,omitempty" fake:"{randomstring:[opened,closed,locked,merged,all]}"`
	OrderBy                string     `json:"order_by,omitempty" fake:"{randomstring:[created_at,updated_at,title]}"`
	Sort                   string     `json:"sort,omitempty" fake:"{randomstring:[asc,desc]}"`
	WithLabelsDetails      bool       `json:"with_labels_details,omitempty" fake:"{bool}"`
	WithMergeStatusRecheck bool       `json:"with_merge_status_recheck,omitempty" fake:"{bool}"`
	CreatedAfter           *time.Time `json:"created_after,omitempty"`
	CreatedBefore          *time.Time `json:"created_before,omitempty"`
	UpdatedAfter           *time.Time `json:"updated_after,omitempty"`
	UpdatedBefore          *time.Time `json:"updated_before,omitempty"`
	View                   string     `json:"view,omitempty" fake:"{randomstring:[simple,normal]}"`
	Scope                  string     `json:"scope,omitempty" fake:"{randomstring:[created-by-me,assigned-to-me,all]}"`
	SourceBranch           string     `json:"source_branch,omitempty" fake:"{word}"`
	SourceProjectID        *int       `json:"source_project_id,omitempty"`
	TargetBranch           string     `json:"target_branch,omitempty" fake:"{randomstring:[main,master,develop]}"`
	Search                 string     `json:"search,omitempty" fake:"{word}"`
	In                     string     `json:"in,omitempty" fake:"{randomstring:[title,description]}"`
	WIP                    string     `json:"wip,omitempty" fake:"{randomstring:[yes,no]}"`
	NotAuthorID            *int       `json:"not_author_id,omitempty"`
	NotAuthorUsername      string     `json:"not_author_username,omitempty" fake:"{username}"`
	NotAssigneeID          *int       `json:"not_assignee_id,omitempty"`
	NotAssigneeUsername    []string   `json:"not_assignee_username,omitempty"`
}

// ProjectsQueryParams represents the query parameters for the projects endpoint
type ProjectsQueryParams struct {
	PageQueryParams

	Archived                 *bool  `json:"archived,omitempty"`
	Visibility               string `json:"visibility,omitempty" fake:"{randomstring:[public,private,internal]}"`
	Search                   string `json:"search,omitempty" fake:"{word}"`
	OrderBy                  string `json:"order_by,omitempty" fake:"{randomstring:[id,name,path,created_at,updated_at]}"`
	Sort                     string `json:"sort,omitempty" fake:"{randomstring:[asc,desc]}"`
	Simple                   bool   `json:"simple,omitempty" fake:"{bool}"`
	Owned                    bool   `json:"owned,omitempty" fake:"{bool}"`
	Starred                  bool   `json:"starred,omitempty" fake:"{bool}"`
	WithIssuesEnabled        bool   `json:"with_issues_enabled,omitempty" fake:"{bool}"`
	WithMergeRequestsEnabled bool   `json:"with_merge_requests_enabled,omitempty" fake:"{bool}"`
	WithShared               bool   `json:"with_shared,omitempty" fake:"{bool}"`
	IncludeSubgroups         bool   `json:"include_subgroups,omitempty" fake:"{bool}"`
	IncludeAncestorGroups    bool   `json:"include_ancestor_groups,omitempty" fake:"{bool}"`
	MinAccessLevel           *int   `json:"min_access_level,omitempty"`
	WithCustomAttributes     bool   `json:"with_custom_attributes,omitempty" fake:"{bool}"`
	WithSecurityReports      bool   `json:"with_security_reports,omitempty" fake:"{bool}"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetGroupsProjects(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

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
	// In a real implementation, this would call h.service.mergeRequests.GetGroupMergeRequests(groupID, params)
	mergeRequests := h.generateMockMergeRequests(params)

	setOnePagedHeaders(len(mergeRequests), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(mergeRequests); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
func (h *Handler) generateMockMergeRequests(params *MergeRequestsQueryParams) []MergeRequest {
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
				Username: "merge_user",
				Name:     "merge user",
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
				Username: "merge_user",
				Name:     "merge user",
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
				Username: "merge_user",
				Name:     "merge user",
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

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path parameter using regex
	// Expected path: /api/v4/users/{id}/events
	/*
		re := regexp.MustCompile(`/api/v4/users/([^/]+)/events`)
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) < 2 {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}
		userID := matches[1]
	*/

	// Parse query parameters
	params, err := h.parseEventsQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, return mock events based on the API specification
	// In a real implementation, this would call h.service.events.GetUserEvents(userID, params)
	events := h.generateMockEvents(params)

	setOnePagedHeaders(len(events), params.PageQueryParams, w.Header())
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
	header.Set("X-Page", "1")
	header.Set("X-Next-Page", "2")
	header.Set("X-Prev-Page", "0")
	header.Set("X-Per-Page", strconv.Itoa(params.PerPage))
	if params.PerPage < total {
		panic("count of collection exceeds PerPage")
	}
	/*
		// actual pagination headers would need to spilt the collection
		header.Set("X-Page", strconv.Itoa(params.Page))
		header.Set("X-Next-Page", strconv.Itoa(params.Page+1))
		header.Set("X-Prev-Page", strconv.Itoa(max(0, params.Page-1)))
		header.Set("X-Per-Page", strconv.Itoa(params.PerPage))
	*/

	header.Set("X-Total-Pages", "1")
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
func (h *Handler) generateMockEvents(params *EventsQueryParams) []Event {
	// Generate mock events based on the API specification

	events := []Event{}
	for range 20 {
		event := new(Event)
		err := gofakeit.Struct(event)
		if err != nil {
			fmt.Printf("could not make struct: %+v\n", err)
		}
		fmt.Println("---------------------")
		spew.Dump(event)
		events = append(events, *event)
	}
	/*
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
				AuthorUsername: stringPtr("event_user"),
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
				AuthorUsername: stringPtr("event_user"),
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
				AuthorUsername: stringPtr("event_user"),
				Imported:       false,
				ImportedFrom:   "none",
			},
		}
	*/

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
