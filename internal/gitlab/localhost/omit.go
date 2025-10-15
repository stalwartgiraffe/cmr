package localhost

import (
	"fmt"
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
	return r(o, reflect.TypeOf(v), reflect.ValueOf(v), "", 0, 0)
}

// r handles well known types
func r(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
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
		return rPointer(o, t, v, tag, size, depth)
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
		return rSlice(o, t, v, tag, size, depth)
	case reflect.Map:
		return rMap(o, t, v, tag, size, depth)
	}
	return nil
}

func isOmitempty(s string) bool {
	return strings.Contains(s, "omitempty")
}

func rPointer(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
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
	return r(o, elemT, v.Elem(), tag, size, depth+1)
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
	// Prevent recursing deeper than configured levels
	if depth >= RecursiveDepth {
		return nil
	}

	// Loop through all the fields of the struct
	n := t.NumField()
	for i := range n {

		elementT := t.Field(i)
		fieldName := elementT.Name
		if fieldName == "AllAddress" {
			fmt.Println("found it")
		}

		fmt.Println(fieldName)
		elementV := v.Field(i)
		jsonTag, ok := elementT.Tag.Lookup("json")
		if !ok {
			continue
		}

		// Check to make sure you can set it or that it's an embedded(anonymous) field
		if !elementV.CanSet() && !elementT.Anonymous {
			continue
		}

		// Check if fakesize is set
		size := -1 // Set to -1 to indicate fakesize was not set

		// Recursively call r() to fill in the struct
		err := r(o, elementT.Type, elementV, jsonTag, size, depth+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func isJsonOmitempty(jsonTag string) bool {
	parts := strings.Split(jsonTag, ",")
	return slices.Contains(parts, "omitempty")
}

func rSlice(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
	elemStr := t.String()
	typeName := t.Name()
	fmt.Println(elemStr, typeName)
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}

	// Get the element type

	elemT := t.Elem()

	// Loop through the elements length and set based upon the index
	ogSize := size
	for i := 0; i < size; i++ {
		nv := v.Index(i)
		err := r(o, elemT, nv.Elem(), tag, ogSize, depth+1)
		if err != nil {
			return err
		}

		// If values are already set fill them up, otherwise append
		//if elemLen != 0 {
		//	v.Index(i).Set(reflect.Indirect(nv))
		//} else {
		//	v.Set(reflect.Append(reflect.Indirect(v), reflect.Indirect(nv)))
		//}
	}
	return nil
}

func rMap(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
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
