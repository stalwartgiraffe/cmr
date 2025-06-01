package gitlab

import (
	"fmt"

	"github.com/aarondl/opt/omitnull"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

type Note struct {
	ID int `json:"id"`

	NoteableID   int    `json:"noteable_id"`
	NoteableIid  int    `json:"noteable_iid"`
	NoteableType string `json:"noteable_type"`

	Body string `json:"body"`

	CreatedAt Time               `json:"created_at"`
	UpdatedAt omitnull.Val[Time] `json:"updated_at,omitempty"`
	System    bool               `json:"system"`

	Attachment   omitnull.Val[string] `json:"attachment,omitempty"`
	ProjectID    omitnull.Val[int]    `json:"project_id,omitempty"`
	Resolvable   omitnull.Val[bool]   `json:"resolvable,omitempty"`
	Confidential omitnull.Val[bool]   `json:"confidential,omitempty"`
	Internal     omitnull.Val[bool]   `json:"internal,omitempty"`

	// we wish we could use this but it does not work seemlessly with yaml
	//Author null.Val[*AuthorModel] `json:"author,omitempty"`
	Author *AuthorModel `json:"author,omitempty"`
}

func (a *Note) String() string {
	if a == nil {
		return ""
	}

	author := ""
	p := a.Author
	if p != nil {
		author = fmt.Sprintf("{%s}", p.String())
	}
	return utils.Join(
		fmt.Sprint(a.ID),
		a.Body,
		a.Attachment.GetOr(""),
		fmt.Sprint(a.CreatedAt),
		fmt.Sprint(a.System),
		fmt.Sprint(a.NoteableID),
		fmt.Sprint(a.NoteableType),
		fmt.Sprint(a.NoteableIid),
		author,
	)
}

type AuthorModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	State    string `json:"state"`

	Email     omitnull.Val[string] `json:"email,omitempty"`
	CreatedAt omitnull.Val[Time]   `json:"created_at,omitempty"`
	AvatarURL omitnull.Val[string] `json:"avatar_url,omitempty"`
	WebURL    omitnull.Val[string] `json:"web_url,omitempty"`
}

func (a *AuthorModel) String() string {
	if a == nil {
		return ""
	}
	return utils.Join(
		fmt.Sprint(a.ID),
		a.Name,
		a.Username,
		a.State,
	)
	// omit
	// a.AvatarURL
	// a.WebURL

}

/*
this trick is not successfully round tripping
go back to pointers

// MarshalText fixes marshaling to yaml.
// If we use  "github.com/aarondl/opt/null"
//
//	Author null.Val[*AuthorModel] `json:"author,omitempty"`
//
// and marshal it with "gopkg.in/yaml.v2"
// we need to implement encoding.TextMarshaler() or we get this error
//
//	panic: unsupported Scan, storing driver.Value type *gitlab.AuthorModel into type *string
var _ encoding.TextMarshaler = (*AuthorModel)(nil)

func (a *AuthorModel) MarshalText() (text []byte, err error) {
	// calling yaml.Marshal(a) causes an infinite loop in the stack
	// FIXME this is not encoded in the expected format. Do that we need to know who is calling this..
	return []byte(a.String()), nil
}
*/
