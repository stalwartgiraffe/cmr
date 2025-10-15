package localhost

import (
	"fmt"
	"testing"
	"time"

	"github.com/aarondl/opt/omitnull"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func TestOmit(t *testing.T) {
	halfCount := 0
	halfEmpty := func() bool {
		halfCount++
		return halfCount%2 == 0
	}

	tests := []struct {
		name    string
		val     any
		isEmpty isEmptyFn
		wantErr bool
	}{
		{
			name: "test_case_name",
			val:  setupPerson(),
			//isEmpty: alwaysEmpty,
			isEmpty: halfEmpty,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := newOmit()
			o.isEmpty = tt.isEmpty
			err := structFunc(o, tt.val)
			fmt.Println(o.sb.String())

			spew.Dump(tt.val)

			require.True(t, tt.wantErr == (err != nil))
		})
	}
}

func alwaysEmpty() bool {
	return true
}

func setupPerson() *Person {
	tt := utils.Scantime("2020-04-22 17:54")
	return &Person{
		Attachment: omitnull.From("front"),
		ProjectID:  omitnull.From(66),

		Name:  "name",
		Email: "email",
		Age:   42,

		F64:  42.0,
		PF64: utils.Ptr(42.0),

		F32:  42.0,
		PF32: utils.Ptr[float32](42.0),

		I64:  42.0,
		PI64: utils.Ptr[int64](42.0),
		I32:  42.0,
		PI32: utils.Ptr[int32](42.0),
		I16:  42.0,
		PI16: utils.Ptr[int16](42.0),
		I8:   42.0,
		PI8:  utils.Ptr[int8](42.0),

		U64:  42.0,
		PU64: utils.Ptr[uint64](42.0),
		U32:  42.0,
		PU32: utils.Ptr[uint32](42.0),
		U16:  42.0,
		PU16: utils.Ptr[uint16](42.0),
		U8:   42.0,
		PU8:  utils.Ptr[uint8](42.0),

		B:  true,
		PB: utils.Ptr[bool](true),

		UpdatedAt:  utils.Scantime("2024-07-04 17:54"),
		PUpdatedAt: utils.Ptr[time.Time](tt),
		Timeout:    10 * time.Second,
		PTimeout:   utils.Ptr[time.Duration](10 * time.Second),

		Address: &Address{
			Street:  "street",
			City:    "city",
			ZipCode: "zipcode",
			Skip:    "skp",
			Comments: &Comments{
				Title: "title",
				Body:  "body",
			},
		},
		AllAddress: []Address{
			{
				Street:  "0street",
				City:    "0city",
				ZipCode: "0zipcode",
				Skip:    "0skp",
				Comments: &Comments{
					Title: "0title",
					Body:  "0body",
				},
			},
			{
				Street:  "1street",
				City:    "1city",
				ZipCode: "1zipcode",
				Skip:    "1skp",
				Comments: &Comments{
					Title: "1title",
					Body:  "1body",
				},
			},
			{
				Street:  "2street",
				City:    "2city",
				ZipCode: "2zipcode",
				Skip:    "2skp",
				Comments: &Comments{
					Title: "2title",
					Body:  "2body",
				},
			},
		},
		Tags: []string{"apple", "banana", "cantalope"},
		KeySet: map[string]int{
			"umbrella": 1,
			"clock":    2,
			"tea":      3,
		},
	}
}

type Person struct {
	Attachment omitnull.Val[string] `json:"attachment,omitempty"`
	ProjectID  omitnull.Val[int]    `json:"project_id,omitempty"`

	Name   string  `json:"name"`
	PName  *string `json:"pname"`
	Email  string  `json:"email,omitempty"`
	PEmail *string `json:"pemail,omitempty"`

	F64  float64  `json:"f64,omitempty"`
	PF64 *float64 `json:"pf64,omitempty"`
	F32  float32  `json:"f32,omitempty"`
	PF32 *float32 `json:"pf32,omitempty"`

	Age  int    `json:"age,omitempty"`
	PAge *int   `json:"age,omitempty"`
	I64  int64  `json:"i64,omitempty"`
	PI64 *int64 `json:"pi64,omitempty"`
	I32  int32  `json:"i32,omitempty"`
	PI32 *int32 `json:"pi32,omitempty"`
	I16  int16  `json:"i16,omitempty"`
	PI16 *int16 `json:"pi16,omitempty"`
	I8   int8   `json:"i8,omitempty"`
	PI8  *int8  `json:"pi8,omitempty"`

	U64  uint64  `json:"u64,omitempty"`
	PU64 *uint64 `json:"pu64,omitempty"`
	U32  uint32  `json:"pu32,omitempty"`
	PU32 *uint32 `json:"u32,omitempty"`
	U16  uint16  `json:"u16,omitempty"`
	PU16 *uint16 `json:"pu16,omitempty"`
	U8   uint8   `json:"u8,omitempty"`
	PU8  *uint8  `json:"pu8,omitempty"`

	B  bool  `json:"b,omitempty"`
	PB *bool `json:"pb,omitempty"`

	UpdatedAt  time.Time      `json:"updatedat,omitempty"`
	PUpdatedAt *time.Time     `json:"pupdatedat,omitempty"`
	Timeout    time.Duration  `json:"timeout,omitempty"`
	PTimeout   *time.Duration `json:"ptimeout,omitempty"`

	Address    *Address  `json:"address,omitempty"`
	AllAddress []Address `json:"alladdress"`
	Tags       []string  `json:"tags,omitempty"`

	KeySet map[string]int `json:"keyset,omitempty"`
}

type Address struct {
	Skip     string
	Comments *Comments
	Street   string `json:"street,omitempty"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code,omitempty"`
}

type Comments struct {
	Title string
	Body  string `json:"body,omitempty"`
}
