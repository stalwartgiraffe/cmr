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

	t.Run("pointer to string", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Str *string `json:"str,omitempty"`
		}

		str := "test"
		s := &StructType{
			Str: &str,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Nil(t, s.Str)
	})

	t.Run("string without omitempty", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Str string `json:"str"`
		}

		s := &StructType{
			Str: "test",
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Equal(t, "test", s.Str) // should not be cleared
	})

	t.Run("string with omitempty", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Str string `json:"str,omitempty"`
		}

		s := &StructType{
			Str: "test",
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.Str)
	})

	t.Run("uint types", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			U8  uint8  `json:"u8,omitempty"`
			U16 uint16 `json:"u16,omitempty"`
			U32 uint32 `json:"u32,omitempty"`
			U64 uint64 `json:"u64,omitempty"`
		}

		s := &StructType{
			U8:  8,
			U16: 16,
			U32: 32,
			U64: 64,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.U8)
		require.Zero(t, s.U16)
		require.Zero(t, s.U32)
		require.Zero(t, s.U64)
	})

	t.Run("int types", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			I8  int8  `json:"i8,omitempty"`
			I16 int16 `json:"i16,omitempty"`
			I32 int32 `json:"i32,omitempty"`
			I64 int64 `json:"i64,omitempty"`
		}

		s := &StructType{
			I8:  8,
			I16: 16,
			I32: 32,
			I64: 64,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.I8)
		require.Zero(t, s.I16)
		require.Zero(t, s.I32)
		require.Zero(t, s.I64)
	})

	t.Run("float types", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			F32 float32 `json:"f32,omitempty"`
			F64 float64 `json:"f64,omitempty"`
		}

		s := &StructType{
			F32: 32.5,
			F64: 64.5,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.F32)
		require.Zero(t, s.F64)
	})

	t.Run("bool type", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			B bool `json:"b,omitempty"`
		}

		s := &StructType{
			B: true,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.B)
	})

	t.Run("slice with omitempty", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Slice []string `json:"slice,omitempty"`
		}

		s := &StructType{
			Slice: []string{"a", "b", "c"},
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Nil(t, s.Slice)
	})

	t.Run("slice of structs", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty

		type Item struct {
			Name string `json:"name,omitempty"`
		}

		type StructType struct {
			Items []Item `json:"items"`
		}

		s := &StructType{
			Items: []Item{
				{Name: "item1"},
				{Name: "item2"},
			},
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Len(t, s.Items, 2)
		require.Zero(t, s.Items[0].Name)
		require.Zero(t, s.Items[1].Name)
	})

	t.Run("nested struct", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty

		type Inner struct {
			Value string `json:"value,omitempty"`
		}

		type StructType struct {
			Inner Inner `json:"inner"`
		}

		s := &StructType{
			Inner: Inner{Value: "test"},
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Zero(t, s.Inner.Value)
	})

	t.Run("pointer to nested struct", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty

		type Inner struct {
			Value string `json:"value,omitempty"`
		}

		type StructType struct {
			Inner *Inner `json:"inner,omitempty"`
		}

		s := &StructType{
			Inner: &Inner{Value: "test"},
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Nil(t, s.Inner)
	})

	t.Run("map with int values and omitempty", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Mset map[string]int `json:"mset,omitempty"`
		}

		s := &StructType{
			Mset: map[string]int{
				"first":  1,
				"second": 2,
			},
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Nil(t, s.Mset)
	})

	t.Run("nil pointer", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			Ptr *string `json:"ptr,omitempty"`
		}

		s := &StructType{
			Ptr: nil,
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Nil(t, s.Ptr)
	})

	t.Run("field without json tag", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty
		type StructType struct {
			NoTag  string
			Tagged string `json:"tagged,omitempty"`
		}

		s := &StructType{
			NoTag:  "no-tag",
			Tagged: "tagged",
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Equal(t, "no-tag", s.NoTag) // should not be touched
		require.Zero(t, s.Tagged)
	})

	t.Run("anonymous field", func(t *testing.T) {
		t.Parallel()

		o := newOmit()
		o.isSetEmpty = alwaysEmpty

		type Embedded struct {
			Value string `json:"value,omitempty"`
		}

		type StructType struct {
			Embedded
			Other string `json:"other,omitempty"`
		}

		s := &StructType{
			Embedded: Embedded{Value: "embedded"},
			Other:    "other",
		}
		err := structFunc(o, s)

		require.NoError(t, err)
		require.Equal(t, "embedded", s.Value) // anonymous fields are skipped
		require.Zero(t, s.Other)
	})
}

func alwaysEmpty() bool {
	return true
}

