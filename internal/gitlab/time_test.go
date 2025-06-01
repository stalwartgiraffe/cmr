package gitlab

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeUnmarshalJSON(t *testing.T) {
	const expectErr = true
	const noErr = false
	testCases := []struct {
		input   string
		output  time.Time
		wantErr bool
	}{
		{"2023-11-23T12:34:56Z", time.Date(2023, 11, 23, 12, 34, 56, 0, time.UTC), noErr},
		{"2023-11-23T12:34:56.999-07:00", time.Date(2023, 11, 23, 12, 34, 56, 999*1000*1000, time.FixedZone("-0700", -7*60*60)), noErr},
		{"2023-11-23", time.Date(2023, 11, 23, 0, 0, 0, 0, time.UTC), noErr},
		{"", time.Time{}, noErr},
		{"null", time.Time{}, noErr},
		{"invalid_format", time.Time{}, expectErr},
		{"2023-13-32", time.Time{}, expectErr},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var timeVal Time
			err := json.Unmarshal([]byte(`"`+tc.input+`"`), &timeVal)
			if (err != nil) != tc.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && !timeVal.Time.Equal(tc.output) {
				t.Errorf("UnmarshalJSON() = %v, want %v", timeVal.Time, tc.output)
			}
		})
	}
}
