package localhost

import (
	"time"
)

// MergeRequest represents a GitLab merge request as per API_Entities_MergeRequestBasic
type MergeRequest struct {
	ID                          int         `json:"id" fake:"{number:1,10000}"`
	IID                         int         `json:"iid" fake:"{number:1,1000}"`
	ProjectID                   int         `json:"project_id" fake:"{number:1,1000}"`
	Title                       string      `json:"title" fake:"{sentence:8}"`
	Description                 string      `json:"description,omitempty" fake:"{paragraph:3,5,10}"`
	State                       string      `json:"state" fake:"{randomstring:[opened,closed,locked,merged]}"`
	CreatedAt                   time.Time   `json:"created_at" fake:"{date}"`
	UpdatedAt                   time.Time   `json:"updated_at" fake:"{date}"`
	MergedBy                    *UserBasic  `json:"merged_by,omitempty"`
	MergeUser                   *UserBasic  `json:"merge_user,omitempty"`
	MergedAt                    time.Time   `json:"merged_at" fake:"{date}"`
	ClosedBy                    *UserBasic  `json:"closed_by,omitempty"`
	ClosedAt                    *string     `json:"closed_at,omitempty"`
	TitleHTML                   string      `json:"title_html,omitempty"`       // TODO extend
	DescriptionHTML             string      `json:"description_html,omitempty"` // TODO extend
	TargetBranch                string      `json:"target_branch" fake:"{randomstring:[main,master,develop,staging]}"`
	SourceBranch                string      `json:"source_branch" fake:"{word}"`
	UserNotesCount              int         `json:"user_notes_count" fake:"{number:0,100}"`
	Upvotes                     int         `json:"upvotes,omitempty" fake:"{number:0,50}"`
	Downvotes                   int         `json:"downvotes,omitempty" fake:"{number:0,10}"`
	Author                      *UserBasic  `json:"author,omitempty"`
	Assignees                   []UserBasic `json:"assignees,omitempty" fakesize:"0,3" `
	Assignee                    *UserBasic  `json:"assignee,omitempty"`
	Reviewers                   []UserBasic `json:"reviewers,omitempty fakesize:"0,3" `
	SourceProjectID             int         `json:"source_project_id,omitempty" fake:"{number:1,1000}"`
	TargetProjectID             int         `json:"target_project_id,omitempty" fake:"{number:1,1000}"`
	Labels                      []string    `json:"labels,omitempty" fakesize:"0,3" fake:"{word}"`
	Draft                       bool        `json:"draft,omitempty" fake:"{bool}"`
	Imported                    bool        `json:"imported,omitempty" fake:"{bool}"`
	ImportedFrom                string      `json:"imported_from,omitempty" fake:"{randomstring:[none,github,bitbucket,gitlab]}"`
	WorkInProgress              bool        `json:"work_in_progress,omitempty" fake:"{bool}"`
	MergeWhenPipelineSucceeds   bool        `json:"merge_when_pipeline_succeeds,omitempty" fake:"{bool}"`
	MergeStatus                 string      `json:"merge_status,omitempty" fake:"{randomstring:[can_be_merged,cannot_be_merged,unchecked]}"`
	DetailedMergeStatus         string      `json:"detailed_merge_status,omitempty" fake:"{randomstring:[mergeable,broken_status,checking,ci_must_pass,ci_still_running,discussions_not_resolved,draft_status]}"`
	MergeAfter                  string      `json:"merge_after,omitempty" fake:"{date}"`
	SHA                         string      `json:"sha,omitempty" fake:"{uuid}"`
	MergeCommitSHA              string      `json:"merge_commit_sha,omitempty" fake:"{uuid}"`
	SquashCommitSHA             string      `json:"squash_commit_sha,omitempty" fake:"{uuid}"`
	DiscussionLocked            bool        `json:"discussion_locked,omitempty" fake:"{bool}"`
	ShouldRemoveSourceBranch    bool        `json:"should_remove_source_branch,omitempty" fake:"{bool}"`
	ForceRemoveSourceBranch     bool        `json:"force_remove_source_branch,omitempty" fake:"{bool}"`
	PreparedAt                  string      `json:"prepared_at,omitempty" fake:"{date}"`
	AllowCollaboration          bool        `json:"allow_collaboration,omitempty" fake:"{bool}"`
	AllowMaintainerToPush       bool        `json:"allow_maintainer_to_push,omitempty" fake:"{bool}"`
	Reference                   string      `json:"reference,omitempty" fake:"{word}"`
	WebURL                      string      `json:"web_url,omitempty" fake:"{url}"`
	Squash                      bool        `json:"squash,omitempty" fake:"{bool}"`
	SquashOnMerge               bool        `json:"squash_on_merge,omitempty" fake:"{bool}"`
	TaskCompletionStatus        string      `json:"task_completion_status,omitempty" fake:"{word}"`
	HasConflicts                bool        `json:"has_conflicts,omitempty" fake:"{bool}"`
	BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved,omitempty" fake:"{bool}"`
	ApprovalsBeforeMerge        int         `json:"approvals_before_merge,omitempty" fake:"{number:0,10}"`
}

// TODO get declarative query parser package?
// MergeRequestsQueryParams represents the query parameters for the merge requests endpoint
type MergeRequestsQueryParams struct {
	PageQueryParams

	AuthorID               *int       `json:"author_id,omitempty" fake:"{number:100,999}"`
	AuthorUsername         string     `json:"author_username,omitempty" fake:"{username}"`
	AssigneeID             *int       `json:"assignee_id,omitempty" fake:"{number:100,999}"`
	AssigneeUsername       []string   `json:"assignee_username,omitempty" fake:"{tags}"`
	ReviewerUsername       string     `json:"reviewer_username,omitempty" fake:"{username}"`
	ReviewerID             *int       `json:"reviewer_id,omitempty" fake:"{number:100,999}"`
	Labels                 []string   `json:"labels,omitempty" fake:"{tags}"`
	Milestone              string     `json:"milestone,omitempty" fake:"{word}"`
	MyReactionEmoji        string     `json:"my_reaction_emoji,omitempty" fake:"{emoji}"`
	State                  string     `json:"state,omitempty" fake:"{randomstring:[opened,closed,locked,merged,all]}"`
	OrderBy                string     `json:"order_by,omitempty" fake:"{randomstring:[created_at,updated_at,title]}"`
	Sort                   string     `json:"sort,omitempty" fake:"{randomstring:[asc,desc]}"`
	WithLabelsDetails      bool       `json:"with_labels_details,omitempty" fake:"{bool}"`
	WithMergeStatusRecheck bool       `json:"with_merge_status_recheck,omitempty" fake:"{bool}"`
	CreatedAfter           *time.Time `json:"created_after,omitempty" fake:"{date}"`
	CreatedBefore          *time.Time `json:"created_before,omitempty" fake:"{date}"`
	UpdatedAfter           *time.Time `json:"updated_after,omitempty" fake:"{date}"`
	UpdatedBefore          *time.Time `json:"updated_before,omitempty" fake:"{date}"`
	View                   string     `json:"view,omitempty" fake:"{randomstring:[simple,normal]}"`
	Scope                  string     `json:"scope,omitempty" fake:"{randomstring:[created-by-me,assigned-to-me,all]}"`
	SourceBranch           string     `json:"source_branch,omitempty" fake:"{word}"`
	SourceProjectID        *int       `json:"source_project_id,omitempty" fake:"{number:100,999}"`
	TargetBranch           string     `json:"target_branch,omitempty" fake:"{randomstring:[main,master,develop]}"`
	Search                 string     `json:"search,omitempty" fake:"{word}"`
	In                     string     `json:"in,omitempty" fake:"{randomstring:[title,description]}"`
	WIP                    string     `json:"wip,omitempty" fake:"{randomstring:[yes,no]}"`
	NotAuthorID            *int       `json:"not_author_id,omitempty"`
	NotAuthorUsername      string     `json:"not_author_username,omitempty" fake:"{username}"`
	NotAssigneeID          *int       `json:"not_assignee_id,omitempty" fake:"{number:100,999}"`
	NotAssigneeUsername    []string   `json:"not_assignee_username,omitempty" fake:"{tags}"`
}
