package gitlab

import (
	"github.com/aarondl/opt/omitnull"
)

//easyjson:json
type MergeRequestModelSlice []MergeRequestModel
type MergeRequestModel struct {
	ID                          int                  `json:"id"`
	Iid                         int                  `json:"iid"`
	ProjectID                   int                  `json:"project_id"`
	Title                       string               `json:"title"`
	Description                 string               `json:"description"`
	State                       string               `json:"state"`
	Imported                    bool                 `json:"imported"`
	ImportedFrom                string               `json:"imported_from"`
	MergedBy                    *UserModel           `json:"merged_by,omitempty"`
	MergeUser                   *UserModel           `json:"merge_user,omitempty"`
	MergedAt                    Time                 `json:"merged_at"`
	PreparedAt                  Time                 `json:"prepared_at"`
	ClosedBy                    *UserModel           `json:"closed_by,omitempty"`
	ClosedAt                    omitnull.Val[string] `json:"closed_at,omitempty"`
	CreatedAt                   Time                 `json:"created_at"`
	UpdatedAt                   Time                 `json:"updated_at"`
	TargetBranch                string               `json:"target_branch"`
	SourceBranch                string               `json:"source_branch"`
	Upvotes                     int                  `json:"upvotes"`
	Downvotes                   int                  `json:"downvotes"`
	Author                      *UserModel           `json:"author,omitempty"`
	Assignee                    *UserModel           `json:"assignee,omitempty"`
	Assignees                   []UserModel          `json:"assignees"`
	Reviewers                   []UserModel          `json:"reviewers"`
	SourceProjectID             int                  `json:"source_project_id"`
	TargetProjectID             int                  `json:"target_project_id"`
	Labels                      []string             `json:"labels"`
	Draft                       bool                 `json:"draft"`
	WorkInProgress              bool                 `json:"work_in_progress"`
	Milestone                   omitnull.Val[string] `json:"milestone,omitempty"`
	MergeWhenPipelineSucceeds   bool                 `json:"merge_when_pipeline_succeeds"`
	MergeStatus                 string               `json:"merge_status"`
	DetailedMergeStatus         string               `json:"detailed_merge_status"`
	Sha                         string               `json:"sha"`
	MergeCommitSha              omitnull.Val[string] `json:"merge_commit_sha,omitempty"`
	SquashCommitSha             omitnull.Val[string] `json:"squash_commit_sha,omitempty"`
	UserNotesCount              int                  `json:"user_notes_count"`
	DiscussionLocked            omitnull.Val[bool]   `json:"discussion_locked,omitempty"`
	ShouldRemoveSourceBranch    omitnull.Val[bool]   `json:"should_remove_source_branch,omitempty"`
	ForceRemoveSourceBranch     bool                 `json:"force_remove_source_branch"`
	AllowCollaboration          bool                 `json:"allow_collaboration"`
	AllowMaintainerToPush       bool                 `json:"allow_maintainer_to_push"`
	WebURL                      string               `json:"web_url"`
	References                  *ReferencesModel     `json:"references,omitempty"`
	TimeStats                   *TimeStatsModel      `json:"time_stats,omitempty"`
	Squash                      bool                 `json:"squash"`
	SquashOnMerge               bool                 `json:"squash_on_merge"`
	TaskCompletionStatus        *UserBasic           `json:"task_completion_status"`
	HasConflicts                bool                 `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                 `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        omitnull.Val[int]    `json:"approvals_before_merge,omitempty"`
}

// comment
type ReferencesModel struct {
	Short    string `json:"short"`
	Relative string `json:"relative"`
	Full     string `json:"full"`
}

type TimeStatsModel struct {
	TimeEstimate        int                  `json:"time_estimate"`
	TotalTimeSpent      int                  `json:"total_time_spent"`
	HumanTimeEstimate   omitnull.Val[string] `json:"human_time_estimate,omitempty"`
	HumanTotalTimeSpent omitnull.Val[string] `json:"human_total_time_spent,omitempty"`
}

type UserBasic struct {
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}
