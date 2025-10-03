package merge_request_examples

import (
	"github.com/mailru/easyjson/opt"
	"time"
)

// User represents a GitLab user
type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	AvatarURL opt.String `json:"avatar_url"`
	WebURL    string     `json:"web_url"`
}

// References contains reference strings for the merge request
type References struct {
	Short    string `json:"short"`
	Relative string `json:"relative"`
	Full     string `json:"full"`
}

// TimeStats contains time tracking statistics
type TimeStats struct {
	TimeEstimate        int        `json:"time_estimate"`
	TotalTimeSpent      int        `json:"total_time_spent"`
	HumanTimeEstimate   opt.String `json:"human_time_estimate"`
	HumanTotalTimeSpent opt.String `json:"human_total_time_spent"`
}

// TaskCompletionStatus tracks task completion in the merge request
type TaskCompletionStatus struct {
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}

// Milestone represents a GitLab milestone (kept as generic object since it's nullable and structure not specified)
type Milestone map[string]interface{}

// MergeRequest represents a GitLab merge request
type MergeRequest struct {
	ID                          int                  `json:"id"`
	IID                         int                  `json:"iid"`
	ProjectID                   int                  `json:"project_id"`
	Title                       string               `json:"title"`
	Description                 string               `json:"description"`
	State                       string               `json:"state"`
	CreatedAt                   time.Time            `json:"created_at"`
	UpdatedAt                   time.Time            `json:"updated_at"`
	MergedBy                    *User                `json:"merged_by"`
	MergeUser                   *User                `json:"merge_user"`
	MergedAt                    *time.Time           `json:"merged_at"`
	ClosedBy                    *User                `json:"closed_by"`
	ClosedAt                    *time.Time           `json:"closed_at"`
	TargetBranch                string               `json:"target_branch"`
	SourceBranch                string               `json:"source_branch"`
	UserNotesCount              int                  `json:"user_notes_count"`
	Upvotes                     int                  `json:"upvotes"`
	Downvotes                   int                  `json:"downvotes"`
	Author                      User                 `json:"author"`
	Assignees                   []User               `json:"assignees"`
	Assignee                    *User                `json:"assignee"`
	Reviewers                   []User               `json:"reviewers"`
	SourceProjectID             int                  `json:"source_project_id"`
	TargetProjectID             int                  `json:"target_project_id"`
	Labels                      []string             `json:"labels"`
	Draft                       bool                 `json:"draft"`
	WorkInProgress              bool                 `json:"work_in_progress"`
	Milestone                   *Milestone           `json:"milestone"`
	MergeWhenPipelineSucceeds   bool                 `json:"merge_when_pipeline_succeeds"`
	MergeStatus                 string               `json:"merge_status"`
	DetailedMergeStatus         string               `json:"detailed_merge_status"`
	SHA                         string               `json:"sha"`
	MergeCommitSHA              opt.String           `json:"merge_commit_sha"`
	SquashCommitSHA             opt.String           `json:"squash_commit_sha"`
	DiscussionLocked            opt.Bool             `json:"discussion_locked"`
	ShouldRemoveSourceBranch    opt.Bool             `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch     bool                 `json:"force_remove_source_branch"`
	Reference                   string               `json:"reference"`
	References                  References           `json:"references"`
	WebURL                      string               `json:"web_url"`
	TimeStats                   TimeStats            `json:"time_stats"`
	Squash                      bool                 `json:"squash"`
	SquashOnMerge               bool                 `json:"squash_on_merge"`
	TaskCompletionStatus        TaskCompletionStatus `json:"task_completion_status"`
	HasConflicts                bool                 `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                 `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        opt.Int              `json:"approvals_before_merge"`
}

// MergeRequestList represents a list of merge requests
type MergeRequestList []MergeRequest
