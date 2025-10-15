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

	isEmpty isEmptyFn
}
type isEmptyFn func() bool

const emptyRate = 0.33

func newOmit() *omit {
	flt := func() bool {
		return rand.Float64() <= emptyRate
	}
	return &omit{
		isEmpty: flt,
	}
}

func Omit(v any) error {
	o := newOmit()
	return structFunc(o, v)
}

func structFunc(o *omit, v any) error {
	return r(o, reflect.TypeOf(v), reflect.ValueOf(v), "", 0, 0)
}

// r handles well known native types
func r(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
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

	if !o.isEmpty() {
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
		parts := strings.Split(jsonTag, ",")
		if !slices.Contains(parts, "omitempty") {
			continue
		}

		// Check to make sure you can set it or that it's an embedded(anonymous) field
		if !elementV.CanSet() && !elementT.Anonymous {
			continue
		}

		// Check if reflect type is of values we can specifically set
		elemStr := elementT.Type.String()
		switch {
		case strings.HasPrefix(elemStr, "omitnull.Val["):
			if err := rOmitnull(o, elementT.Type, elementV, jsonTag); err != nil {
				return err
			}
			continue
		case elemStr == "time.Time", elemStr == "*time.Time":
			// Check if element is a pointer
			elemV := elementV
			if elemStr == "*time.Time" {
				elemV = reflect.New(elementT.Type.Elem()).Elem()
			}

			// Run rTime on the element
			if err := rTime(o, elementT, elemV, jsonTag); err != nil {
				return err
			}

			if elemStr == "*time.Time" {
				elementV.Set(elemV.Addr())
			}

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
func rTime(o *omit, t reflect.StructField, v reflect.Value, tag string) error {
	timeStruct := time.Time{}
	v.Set(reflect.ValueOf(timeStruct))
	return nil
}
