package gitlab

import (
	"fmt"

	"github.com/aarondl/opt/omitnull"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

type EventModel struct {
	ID             int    `json:"id"`
	ProjectID      int    `json:"project_id"`
	TargetID       int    `json:"target_id"`
	TargetIid      int    `json:"target_iid"`
	AuthorID       int    `json:"author_id"`
	AuthorUsername string `json:"author_username"`

	Title       omitnull.Val[string] `json:"title,omitempty"`
	ActionName  string               `json:"action_name"`
	TargetType  string               `json:"target_type"`
	TargetTitle omitnull.Val[string] `json:"target_title,omitempty"`
	CreatedAt   Time                 `json:"created_at"`

	Data omitnull.Val[string] `json:"data,omitempty"`

	Author *AuthorModel `json:"author,omitempty"`
	Note   *Note        `json:"note,omitempty"`

	Imported     omitnull.Val[bool]   `json:"imported,omitempty"`
	ImportedFrom omitnull.Val[string] `json:"imported_from,omitempty"`

	// this does not round trip, use raw pointer
	//	panic: unsupported Scan, storing driver.Value
	//PushData     omitnull.Val[PushDataModel] `json:"push_data,omitempty"`
	PushData *PushDataModel `json:"push_data,omitempty"`
}

type PushDataModel struct {
	CommitCount int    `json:"commit_count"`
	Action      string `json:"action"`
	RefType     string `json:"ref_type"`
	CommitFrom  string `json:"commit_from"`
	CommitTo    string `json:"commit_to"`
	Ref         string `json:"ref"`
	CommitTitle string `json:"commit_title"`
}

// fix
//	panic: unsupported Scan, storing driver.Value
//func (p *PushDataModel) MarshalText() (text []byte, err error) {
//	return []byte(p.String()), nil
//}

//easyjson:json
type EventModelSlice []EventModel

var EventModelFieldNames = getJsonKeys(EventModel{})

func (m *PushDataModel) String() string {
	if m == nil {
		return ""
	}

	return utils.Join(
		fmt.Sprint(m.CommitCount),
		m.Action,
		m.CommitTitle,
		m.RefType,
		m.CommitFrom,
		m.CommitTo,
		m.Ref,
	)
}
