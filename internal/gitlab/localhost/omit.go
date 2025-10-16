package localhost

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"time"
)

var RecursiveDepth = 10

type omit struct {
	sb strings.Builder

	isSetEmpty isEmptyFn
}
type isEmptyFn func() bool

const emptyRate = 0.33

func newOmit() *omit {
	flt := func() bool {
		return rand.Float64() <= emptyRate
	}
	return &omit{
		isSetEmpty: flt,
	}
}

func Omit(v any) error {
	o := newOmit()
	return structFunc(o, v)
}

func structFunc(o *omit, v any) error {
	return r(o, reflect.TypeOf(v), reflect.ValueOf(v), "", 0)
}

// r handles well known types
func r(o *omit, t reflect.Type, v reflect.Value, tag string, depth int) error {
	typeStr := t.String()
	// handle types that are not built into the reflect package
	switch {
	case typeStr == "time.Time", typeStr == "*time.Time":
		return rTime(o, t, v, tag)
	case strings.HasPrefix(typeStr, "omitnull.Val["):
		return rOmitnull(o, t, v, tag)
	}

	switch t.Kind() {
	case reflect.Ptr:
		return rPointer(o, t, v, tag, depth)
	case reflect.Struct:
		return rStruct(o, t, v, tag, depth)
	case reflect.String:
		return rString(o, t, v, tag)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rUint(o, t, v, tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rInt(o, t, v, tag)
	case reflect.Float32, reflect.Float64:
		return rFloat(o, t, v, tag)
	case reflect.Bool:
		return rBool(o, t, v, tag)
	case reflect.Array, reflect.Slice:
		return rSlice(o, t, v, tag, depth)
	case reflect.Map:
		return rMap(o, t, v, tag, depth)
	}
	return nil
}

func isOmitempty(s string) bool {
	return strings.Contains(s, "omitempty")
}

func rPointer(o *omit, t reflect.Type, v reflect.Value, tag string, depth int) error {
	elemT := t.Elem()
	// Prevent recursing deeper than configured levels
	if depth >= RecursiveDepth {
		return nil
	}

	if v.IsNil() {
		return nil
	}
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}
	return r(o, elemT, v.Elem(), tag, depth+1)
}

// setOmitEmpty sets the value to zero
func setOmitEmpty(o *omit, t reflect.Type, v reflect.Value, tag string) bool {
	if !isOmitempty(tag) {
		return false
	}
	if !v.CanSet() {
		return false
	}

	if !o.isSetEmpty() {
		return false
	}
	v.Set(reflect.Zero(t))
	return true
}

func rStruct(o *omit, t reflect.Type, v reflect.Value, tag string, depth int) error {
	if depth >= RecursiveDepth {
		return nil
	}
	for i := range t.NumField() {
		fieldT := t.Field(i)
		if fieldT.Anonymous {
			continue
		}
		jsonTag, ok := fieldT.Tag.Lookup("json")
		if !ok { // the struct field is not tagged with json
			continue
		}
		fieldV := v.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		if err := r(o, fieldT.Type, fieldV, jsonTag, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func isJsonOmitempty(jsonTag string) bool {
	parts := strings.Split(jsonTag, ",")
	return slices.Contains(parts, "omitempty")
}

func rSlice(o *omit, t reflect.Type, v reflect.Value, tag string, depth int) error {
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}

	elemT := t.Elem()
	for i := range v.Len() {
		if err := r(o, elemT, v.Index(i), "", depth+1); err != nil {
			return err
		}
	}
	return nil
}

func rMap(o *omit, t reflect.Type, v reflect.Value, tag string, depth int) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rString(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rInt(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rUint(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rFloat(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rBool(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

func rOmitnull(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	setOmitEmpty(o, t, v, tag)
	return nil
}

// rTime will set a time.Time field the best it can from either the default date tag or from the generate tag
func rTime(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	if !isJsonOmitempty(tag) {
		return nil
	}
	if !o.isSetEmpty() {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		v.Set(reflect.Zero(t))
	} else {
		v.Set(reflect.ValueOf(time.Time{}))
	}
	return nil
}
