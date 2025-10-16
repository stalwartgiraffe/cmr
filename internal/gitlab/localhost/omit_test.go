package localhost

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonAlwaysEmpty(t *testing.T) {
	t.Run("one int member", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			I1 int `json:"f1,omitempty"`
		}

		s := &StructType{
			I1: 123,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.I1)
	})

	t.Run("map with struct value member", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty

		type ValueType struct {
			Body string `json:"body,omitempty"`
		}

		type StructType struct {
			Mset map[string]ValueType `json:"mset"`
		}

		s := &StructType{
			Mset: map[string]ValueType{
				"first": {
					Body: "body",
				},
			},
		}
		err := structFunc(o, s)

		require.NoError(t, err)

		require.True(t, s.Mset["first"].Body == "")
	})
}

func alwaysEmpty() bool {
	return true
}

func TestOmitBody(t *testing.T) {
}
