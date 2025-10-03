package localhost

import (
	"time"
)

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

// UserBasicV0 represents a GitLab user as per API_Entities_UserBasic
type UserBasicV0 struct {
	ID          int    `json:"id" fake:"{number:100,999}"`
	Username    string `json:"username" fake:"{username}"`
	PublicEmail string `json:"public_email,omitempty" fake:"{email}"`
	Name        string `json:"name" fake:"{name}"`
	State       string `json:"state" fake:"{randomstring:[active,blocked,deactivated]}"`
	Locked      bool   `json:"locked" fake:"{bool}"`
	AvatarURL   string `json:"avatar_url,omitempty" fake:"{imageurl:200,200}"`
	AvatarPath  string `json:"avatar_path,omitempty" fake:"{filePath}"`
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
