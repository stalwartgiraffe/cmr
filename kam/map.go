package kam

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"gopkg.in/yaml.v2"
)

type Map map[string]any

func AsMap(a any) (Map, error) {
	// in type assertion and type conversions are not casting
	// https://stackoverflow.com/questions/19577423/how-to-cast-to-a-type-alias-in-go
	if a == nil {
		return nil, nil
	}

	m, ok := a.(map[string]any) // type assert
	if !ok {
		return nil, fmt.Errorf("Value can not be asserted to map[string]any: %v", a)
	}
	return Map(m), nil // type convert
}

func NewMap(s string) (Map, error) {
	return NewMapWithByte([]byte(s))
}

func NewMapWithByte(b []byte) (Map, error) {
	m := Map{}

	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("Could not unmarshal %w", err)
	}
	return m, nil
}

// ToQueryParameters return the map in URL query parameter format.
func (m Map) ToQueryParameters() string {
	if m == nil {
		return ""
	}

	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb := strings.Builder{}
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%v", k, m[k]))
		if i+1 < len(m) {
			sb.WriteString("&")
		}
	}
	return sb.String()
}

func (m Map) ToYaml() string {
	b, _ := yaml.Marshal(m)
	return string(b)
}

// Keys returns the map keys as a slice.
func (m Map) Keys() []string {
	return maps.Keys(m) // order will be randomized by the language
}

// SortedKeys returns the map keys as a sorted slice.
func (m Map) SortedKeys() []string {
	k := m.Keys()
	slices.Sort(k)
	return k
}

// Values returns the values as a slice.
func (m Map) Values() []any {
	return maps.Values(m) // order will randomized by the language
}

// anyClone makes a deep clone of any on a supported type.
func anyClone(anyv any) any {
	switch v := anyv.(type) {
	case complex64,
		complex128,
		float32,
		float64,
		uint8,
		uint16,
		uint32,
		uint64,
		int8,
		int16,
		int32,
		int64,
		int,
		uintptr,
		error,
		string,
		bool: // should be all primitives
		return v
	case map[string]any:
		return Map(v).Clone()
	case Map:
		return v.Clone()
	case []any:
		a := make([]any, 0, len(v))
		for _, anyElement := range v {
			a = append(a, anyClone(anyElement))
		}
		return a
	default:
		if v == nil {
			return nil
		}
		panic(fmt.Sprintf("unexpected %T %+v", v, v))
	}
}

// Clone returns a deep clone of m
func (m Map) Clone() Map {
	c := Map{}
	for k, anyv := range m {
		c[k] = anyClone(anyv)
	}
	return c
}

// anyEquals returns true iff va == vb on a supported type.
func anyEquals(va, vb any) bool {
	switch ca := va.(type) {
	case bool:
		if cb, ok := vb.(bool); !ok {
			return false
		} else {
			return ca == cb
		}
	case float64:
		if cb, ok := vb.(float64); !ok {
			return false
		} else {
			return ca == cb
		}
	case string:
		if cb, ok := vb.(string); !ok {
			return false
		} else {
			return ca == cb
		}
	case map[string]any:
		ma := Map(ca)
		var mb Map
		var ok bool
		var mm map[string]any
		mm, ok = vb.(map[string]any)
		if !ok {
			mb, ok = vb.(Map)
			if !ok {
				return false
			}
		} else {
			mb = Map(mm)
		}
		return ma.Equals(mb)
	case Map:
		cb, ok := vb.(Map)
		if !ok {
			return false
		}
		return ca.Equals(cb)
	case []any:
		cb, ok := vb.([]any)
		if !ok {
			return false
		}
		if len(ca) != len(cb) {
			return false
		}
		for i := range ca {
			if !anyEquals(ca[i], cb[i]) {
				return false
			}
		}
		return true
	default:
		if ca == nil {
			return vb == nil
		}
		panic(fmt.Sprintf("unexpected %T %+v", ca, ca))
	}
	return false
}

// Equals - deep equals
// Return true if a contains the same keys and values as b.
//
// Currently only string or Map values are supported.
func (a Map) Equals(b Map) bool {
	if len(a) != len(b) {
		return false
	}
	for k, anyv := range a {
		anyb, ok := b[k]
		if !ok {
			return false
		}
		if !anyEquals(anyv, anyb) {
			return false
		}
	}
	return true
}

// Try and get the value for key k.
func TryGet[T any](m Map, k string) (T, bool) {
	anyv, ok := m[k]
	if ok {
		v, ok := anyv.(T)
		if ok {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Must get value or panic.
func (m Map) Bool(k string) bool {
	v, ok := TryGet[bool](m, k)
	if !ok {
		panic(fmt.Sprintf("key %s of type bool not found", k))
	}
	return v
}

// Must get value or panic.
func (m Map) Float64(k string) float64 {
	v, ok := TryGet[float64](m, k)
	if !ok {
		panic(fmt.Sprintf("key %s of type bool not found", k))
	}
	return v
}

// Must get value or panic.
func (m Map) String(k string) string {
	v, ok := TryGet[string](m, k)
	if !ok {
		panic(fmt.Sprintf("key %s of type bool not found", k))
	}
	return v
}

// Must get value of map or panic. May be nil but key must be present.
func (m Map) Map(k string) Map {
	v, ok := TryGet[map[string]any](m, k)
	if !ok {
		panic(fmt.Sprintf("key %s of type bool not found", k))
	}
	return Map(v)
}

// toIntVal will return the val as a int if the conversion is lossless.
func toIntVal(f float64) any {
	i := int(f)
	a := float64(i)
	if a == f {
		return i
	} else {
		return f
	}
}

func FloatsToIntsWithArray(ar []any) {
	for i, a := range ar {
		switch v := a.(type) {
		case float64:
			ar[i] = toIntVal(v)
		case []any:
			FloatsToIntsWithArray(v)
		case map[string]any:
			FloatToInts(v)
		}
	}

}

// FloatToInts replace floats with in ints if its lossless
func FloatToInts(m map[string]any) {
	for k, a := range m {
		switch v := a.(type) {
		case float64:
			m[k] = toIntVal(v)
		case []any:
			FloatsToIntsWithArray(v)
		case map[string]any:
			FloatToInts(v)
		}
	}
}

type Array []any

type JSONValue struct {
	AnyVal any
}

func (j *JSONValue) UnmarshalJSON(b []byte) error {
	var m Map
	err := json.Unmarshal(b, &m)
	if err == nil {
		FloatToInts(m)
		j.AnyVal = m
		return nil
	}

	var ar Array
	err = json.Unmarshal(b, &ar)
	if err == nil {
		FloatsToIntsWithArray(ar)
		j.AnyVal = ar
		return nil
	}

	return err
}

// Marshal to struct that holds the raw body string
type TextValue struct {
	Val string
}

func (j *TextValue) UnmarshalJSON(b []byte) error {
	j.Val = string(b)
	return nil
}
