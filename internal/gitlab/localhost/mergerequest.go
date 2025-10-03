package localhost

import (
	"github.com/aarondl/opt/omitnull"
	"time"
)

// MergeRequest represents a GitLab merge request
type MergeRequest struct {
	ID                          int                  `json:"id" fake:"{number:1,10000}"`
	IID                         int                  `json:"iid" fake:"{number:1,1000}"`
	ProjectID                   int                  `json:"project_id" fake:"{number:1,1000}"`
	Title                       string               `json:"title" fake:"{sentence:8}"`
	Description                 string               `json:"description" fake:"{paragraph:3,5,10}"`
	State                       string               `json:"state" fake:"{randomstring:[opened,closed,locked,merged]}"`
	CreatedAt                   time.Time            `json:"created_at" fake:"{date}"`
	UpdatedAt                   time.Time            `json:"updated_at" fake:"{date}"`
	MergedBy                    *UserBasic           `json:"merged_by"`
	MergeUser                   *UserBasic           `json:"merge_user"`
	MergedAt                    *time.Time           `json:"merged_at"`
	ClosedBy                    *UserBasic           `json:"closed_by"`
	ClosedAt                    *time.Time           `json:"closed_at"`
	TargetBranch                string               `json:"target_branch" fake:"{randomstring:[main,master,develop,staging]}"`
	SourceBranch                string               `json:"source_branch" fake:"{word}"`
	UserNotesCount              int                  `json:"user_notes_count" fake:"{number:0,100}"`
	Upvotes                     int                  `json:"upvotes" fake:"{number:0,50}"`
	Downvotes                   int                  `json:"downvotes" fake:"{number:0,10}"`
	Author                      *UserBasic           `json:"author"`
	Assignees                   []UserBasic          `json:"assignees" fakesize:"0,3"`
	Assignee                    *UserBasic           `json:"assignee"`
	Reviewers                   []UserBasic          `json:"reviewers" fakesize:"0,3"`
	SourceProjectID             int                  `json:"source_project_id" fake:"{number:1,1000}"`
	TargetProjectID             int                  `json:"target_project_id" fake:"{number:1,1000}"`
	Labels                      []string             `json:"labels" fakesize:"0,3" fake:"{word}"`
	Draft                       bool                 `json:"draft" fake:"{bool}"`
	WorkInProgress              bool                 `json:"work_in_progress" fake:"{bool}"`
	Milestone                   *Milestone           `json:"milestone"`
	MergeWhenPipelineSucceeds   bool                 `json:"merge_when_pipeline_succeeds" fake:"{bool}"`
	MergeStatus                 string               `json:"merge_status" fake:"{randomstring:[can_be_merged,cannot_be_merged,unchecked,cannot_be_merged_recheck]}"`
	DetailedMergeStatus         string               `json:"detailed_merge_status" fake:"{randomstring:[mergeable,broken_status,checking,ci_must_pass,ci_still_running,discussions_not_resolved,draft_status,not_approved,not_open,unchecked]}"`
	SHA                         string               `json:"sha" fake:"{uuid}"`
	MergeCommitSHA              string               `json:"merge_commit_sha,omitempty" fake:"{uuid}"`
	SquashCommitSHA             string               `json:"squash_commit_sha,omitempty" fake:"{uuid}"`
	DiscussionLocked            bool                 `json:"discussion_locked,omitempty" fake:"{bool}"`
	ShouldRemoveSourceBranch    bool                 `json:"should_remove_source_branch,omitempty" fake:"{bool}"`
	ForceRemoveSourceBranch     bool                 `json:"force_remove_source_branch" fake:"{bool}"`
	Reference                   string               `json:"reference" fake:"{word}"`
	References                  References           `json:"references"`
	WebURL                      string               `json:"web_url" fake:"{url}"`
	TimeStats                   TimeStats            `json:"time_stats"`
	Squash                      bool                 `json:"squash" fake:"{bool}"`
	SquashOnMerge               bool                 `json:"squash_on_merge" fake:"{bool}"`
	TaskCompletionStatus        TaskCompletionStatus `json:"task_completion_status"`
	HasConflicts                bool                 `json:"has_conflicts" fake:"{bool}"`
	BlockingDiscussionsResolved bool                 `json:"blocking_discussions_resolved" fake:"{bool}"`
	ApprovalsBeforeMerge        omitnull.Val[int]    `json:"approvals_before_merge"`
}

// MergeRequestList represents a list of merge requests
type MergeRequestList []MergeRequest

// UserBasic represents a GitLab user
type UserBasic struct {
	ID        int     `json:"id" fake:"{number:1,10000}"`
	Username  string  `json:"username" fake:"{username}"`
	Name      string  `json:"name" fake:"{name}"`
	State     string  `json:"state" fake:"{randomstring:[active,blocked,deactivated]}"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	WebURL    string  `json:"web_url" fake:"{url}"`
}

// References contains reference strings for the merge request
type References struct {
	Short    string `json:"short" fake:"{word}"`
	Relative string `json:"relative" fake:"{word}"`
	Full     string `json:"full" fake:"{word}"`
}

// TimeStats contains time tracking statistics
type TimeStats struct {
	TimeEstimate        int     `json:"time_estimate" fake:"{number:0,100000}"`
	TotalTimeSpent      int     `json:"total_time_spent" fake:"{number:0,100000}"`
	HumanTimeEstimate   *string `json:"human_time_estimate,omitempty"`
	HumanTotalTimeSpent *string `json:"human_total_time_spent,omitempty"`
}

// TaskCompletionStatus tracks task completion in the merge request
type TaskCompletionStatus struct {
	Count          int `json:"count" fake:"{number:0,20}"`
	CompletedCount int `json:"completed_count" fake:"{number:0,20}"`
}

type Milestone struct {
	ID          int       `json:"id" fake:"{number:100,999}"`
	Iid         int       `json:"iid" fake:"{number:100,999}"`
	GroupID     int       `json:"group_id" fake:"{number:100,999}"`
	ProjectID   int       `json:"project_id" fake:"{number:100,999}"`
	CreatedAt   time.Time `json:"created_at" fake:"{date}"`
	DueDate     time.Time `json:"due_date" fake:"{date}"`
	StartDate   time.Time `json:"start_date" fake:"{date}"`
	UpdateAt    time.Time `json:"updated_at" fake:"{date}"`
	Expired     bool      `json:"expired" fake:"{bool}"`
	Title       string    `json:"title" fake:"{sentence:8}"`
	Description string    `json:"description" fake:"{sentence:8}"`
	State       string    `json:"state" fake:"{word}"`
	WebURL      string    `json:"web_url" fake:"{url}"`
}
