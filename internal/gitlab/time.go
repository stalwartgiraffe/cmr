package gitlab

import (
	"encoding/json"
	"strings"
	"time"
)

// Custom time that supports time stamps from gitlab.
// Explicit null and emtpy string are the zero time value.
type Time struct {
	time.Time
}

var _ json.Unmarshaler = (*Time)(nil)

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	txt := strings.Trim(string(b), `'"`)
	var value time.Time
	var err error

	if len(txt) == 0 ||
		strings.EqualFold(txt, "null") {
		return nil
	}

	if value, err = time.Parse(time.RFC3339, txt); err == nil {
		t.Time = value.In(time.UTC)
	} else if value, err = time.Parse(`2006-01-02T15:04:05.999-07:00`, txt); err == nil {
		t.Time = value.In(time.UTC)
	} else if value, err = time.Parse(time.DateOnly, txt); err == nil {
		t.Time = value.In(time.UTC)
	}

	return err
}
