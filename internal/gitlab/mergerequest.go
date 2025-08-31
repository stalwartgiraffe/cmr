package gitlab

import (
	"time"

	"github.com/mailru/easyjson/opt"
)

//easyjson:json
type BadMergeRequestModelSlice []BadMergeRequestModel
type BadMergeRequestModel struct {
	ID             int    `json:"id"`
	Iid            int    `json:"iid"`
	ProjectID      int    `json:"project_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	State          string `json:"state"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	MergedBy       any    `json:"merged_by"`
	MergeUser      any    `json:"merge_user"`
	MergedAt       any    `json:"merged_at"`
	ClosedBy       any    `json:"closed_by"`
	ClosedAt       any    `json:"closed_at"`
	TargetBranch   string `json:"target_branch"`
	SourceBranch   string `json:"source_branch"`
	UserNotesCount int    `json:"user_notes_count"`
	Upvotes        int    `json:"upvotes"`
	Downvotes      int    `json:"downvotes"`
	Author         struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	Assignees []struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignees"`
	Assignee struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignee"`
	Reviewers                 []any    `json:"reviewers"`
	SourceProjectID           int      `json:"source_project_id"`
	TargetProjectID           int      `json:"target_project_id"`
	Labels                    []string `json:"labels"`
	Draft                     bool     `json:"draft"`
	WorkInProgress            bool     `json:"work_in_progress"`
	Milestone                 any      `json:"milestone"`
	MergeWhenPipelineSucceeds bool     `json:"merge_when_pipeline_succeeds"`
	MergeStatus               string   `json:"merge_status"`
	DetailedMergeStatus       string   `json:"detailed_merge_status"`
	Sha                       string   `json:"sha"`
	MergeCommitSha            any      `json:"merge_commit_sha"`
	SquashCommitSha           any      `json:"squash_commit_sha"`
	DiscussionLocked          any      `json:"discussion_locked"`
	ShouldRemoveSourceBranch  any      `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch   bool     `json:"force_remove_source_branch"`
	Reference                 string   `json:"reference"`
	References                struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	WebURL    string `json:"web_url"`
	TimeStats struct {
		TimeEstimate        int `json:"time_estimate"`
		TotalTimeSpent      int `json:"total_time_spent"`
		HumanTimeEstimate   any `json:"human_time_estimate"`
		HumanTotalTimeSpent any `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Squash               bool `json:"squash"`
	SquashOnMerge        bool `json:"squash_on_merge"`
	TaskCompletionStatus struct {
		Count          int `json:"count"`
		CompletedCount int `json:"completed_count"`
	} `json:"task_completion_status"`
	HasConflicts                bool `json:"has_conflicts"`
	BlockingDiscussionsResolved bool `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        any  `json:"approvals_before_merge"`
}

//easyjson:json
type MergeRequestModelSlice []MergeRequestModel
type MergeRequestModel struct {
	ID                          int                        `json:"id"`
	Iid                         int                        `json:"iid"`
	ProjectID                   int                        `json:"project_id"`
	Title                       string                     `json:"title"`
	Description                 string                     `json:"description"`
	State                       string                     `json:"state"`
	Imported                    bool                       `json:"imported"`
	ImportedFrom                string                     `json:"imported_from"`
	MergedBy                    *AuthorModel               `json:"merged_by,omitempty"`
	MergeUser                   *AuthorModel               `json:"merge_user,omitempty"`
	MergedAt                    Time                       `json:"merged_at"`
	PreparedAt                  Time                       `json:"prepared_at"`
	ClosedBy                    opt.String                 `json:"closed_by"`
	ClosedAt                    opt.String                 `json:"closed_at,omitempty"`
	CreatedAt                   Time                       `json:"created_at"`
	UpdatedAt                   Time                       `json:"updated_at"`
	TargetBranch                string                     `json:"target_branch"`
	SourceBranch                string                     `json:"source_branch"`
	Upvotes                     int                        `json:"upvotes"`
	Downvotes                   int                        `json:"downvotes"`
	Author                      *AuthorModel               `json:"author,omitempty"`
	Assignee                    *AuthorModel               `json:"assignee,omitempty"`
	Assignees                   []AuthorModel              `json:"assignees"`
	Reviewers                   []AuthorModel              `json:"reviewers"`
	SourceProjectID             int                        `json:"source_project_id"`
	TargetProjectID             int                        `json:"target_project_id"`
	Labels                      []string                   `json:"labels"`
	Draft                       bool                       `json:"draft"`
	WorkInProgress              bool                       `json:"work_in_progress"`
	Milestone                   *MilestoneModel            `json:"milestone,omitempty"`
	MergeWhenPipelineSucceeds   bool                       `json:"merge_when_pipeline_succeeds"`
	MergeStatus                 string                     `json:"merge_status"`
	DetailedMergeStatus         string                     `json:"detailed_merge_status"`
	Sha                         string                     `json:"sha"`
	MergeCommitSha              opt.String                 `json:"merge_commit_sha,omitempty"`
	SquashCommitSha             opt.String                 `json:"squash_commit_sha,omitempty"`
	UserNotesCount              int                        `json:"user_notes_count"`
	DiscussionLocked            opt.String                 `json:"discussion_locked,omitempty"`
	ShouldRemoveSourceBranch    opt.Bool                   `json:"should_remove_source_branch,omitempty"`
	ForceRemoveSourceBranch     bool                       `json:"force_remove_source_branch"`
	AllowCollaboration          bool                       `json:"allow_collaboration"`
	AllowMaintainerToPush       bool                       `json:"allow_maintainer_to_push"`
	WebURL                      string                     `json:"web_url"`
	References                  *ReferencesModel           `json:"references,omitempty"`
	TimeStats                   *TimeStatsModel            `json:"time_stats,omitempty"`
	Squash                      bool                       `json:"squash"`
	SquashOnMerge               bool                       `json:"squash_on_merge"`
	TaskCompletionStatus        *TaskCompletionStatusModel `json:"task_completion_status"`
	HasConflicts                bool                       `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                       `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        opt.Bool                   `json:"approvals_before_merge,omitempty"`
}

type ReferencesModel struct {
	Short    string `json:"short"`
	Relative string `json:"relative"`
	Full     string `json:"full"`
}

type TimeStatsModel struct {
	TimeEstimate        int        `json:"time_estimate"`
	TotalTimeSpent      int        `json:"total_time_spent"`
	HumanTimeEstimate   opt.String `json:"human_time_estimate,omitempty"`
	HumanTotalTimeSpent opt.String `json:"human_total_time_spent,omitempty"`
}

type TaskCompletionStatusModel struct {
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}

type MilestoneModel struct {
	ID          int       `json:"id"`
	Iid         int       `json:"iid"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DueDate     string    `json:"due_date"`
	StartDate   string    `json:"start_date"`
	WebURL      string    `json:"web_url"`
}
