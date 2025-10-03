package localhost

import (
	"time"
)

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
	Owner                                     *UserBasicV0    `json:"owner,omitempty"`
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
