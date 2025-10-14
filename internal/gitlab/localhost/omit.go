package localhost

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"time"
)

var RecursiveDepth = 10

type omit struct {
	sb strings.Builder

	pOmitEmpty float64
	randFloat  randFloatFn
}

func newOmit() *omit {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)
	flt := func() float64 {
		return rng.Float64()
	}
	return &omit{
		randFloat: flt,
	}
}

type randFloatFn func() float64

func Omit(v any) error {
	o := newOmit()
	return structFunc(o, v)
}

func structFunc(o *omit, v any) error {
	return r(o, reflect.TypeOf(v), reflect.ValueOf(v), "", 0, 0)
}

func r(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {

	/*
		if t.PkgPath() == "encoding/json" {
			// encoding/json has two special types:
			// - RawMessage
			// - Number

			switch t.Name() {
			case "RawMessage":
				return rJsonRawMessage(o, v, tag)
			case "Number":
				return rJsonNumber(o, v, tag)
			default:
				return errors.New("unknown encoding/json type: " + t.Name())
			}
		}
	*/
	name := t.Name()
	kk := t.Kind().String()
	fmt.Println(name, kk, tag)

	// Handle generic types
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

func setOmitEmpty(o *omit, t reflect.Type, v reflect.Value, tag string) bool {
	if !isOmitempty(tag) {
		return false
	}
	if !v.CanSet() {
		return false
	}

	if o.pOmitEmpty < o.randFloat() {
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
	// Check if tag exists, if so run custom function
	//if t.Name() != "" && tag != "" {
	//	return rCustom(o, v, tag)
	//}

	/*


		// Check if struct is fakeable
		if isFakeable(t) {
			value, err := callFake(f, v, reflect.Struct)
			if err != nil {
				return err
			}

			v.Set(reflect.ValueOf(value))
			return nil
		}
	*/

	// Loop through all the fields of the struct
	n := t.NumField()
	for i := range n {
		elementT := t.Field(i)
		elementV := v.Field(i)
		jsonTag, ok := elementT.Tag.Lookup("json")
		if !ok {
			continue
		}
		parts := strings.Split(jsonTag, ",")
		if !slices.Contains(parts, "omitempty") {
			continue
		}

		//o.sb.WriteString(jsonTag)
		//o.sb.WriteString("\n")

		// Check whether or not to skip this field
		// if ok && jsonTag == "skip" || jsonTag == "-" {
		// 	// Do nothing, skip it
		// 	continue
		// }

		// Check to make sure you can set it or that it's an embedded(anonymous) field
		if !elementV.CanSet() && !elementT.Anonymous {
			continue
		}

		// Check if reflect type is of values we can specifically set
		p2 := elementT.Type.PkgPath()
		p3 := elementT.Type.Name()
		elemStr := elementT.Type.String()
		fmt.Println(p2, p3, elemStr)
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
			err := rTime(o, elementT, elemV, jsonTag)
			if err != nil {
				return err
			}

			if elemStr == "*time.Time" {
				elementV.Set(elemV.Addr())
			}

			continue
		}

		// Check if fakesize is set
		size := -1 // Set to -1 to indicate fakesize was not set
		// fs, ok := elementT.Tag.Lookup("fakesize")
		// if ok {
		// 	var err error
		//
		// 	// Check if size has params separated by ,
		// 	if strings.Contains(fs, ",") {
		// 		sizeSplit := strings.SplitN(fs, ",", 2)
		// 		if len(sizeSplit) == 2 {
		// 			var sizeMin int
		// 			var sizeMax int
		//
		// 			sizeMin, err = strconv.Atoi(sizeSplit[0])
		// 			if err != nil {
		// 				return err
		// 			}
		// 			sizeMax, err = strconv.Atoi(sizeSplit[1])
		// 			if err != nil {
		// 				return err
		// 			}
		//
		// 			size = f.IntN(sizeMax-sizeMin+1) + sizeMin
		// 		}
		// 	} else {
		// 		size, err = strconv.Atoi(fs)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }

		// Recursively call r() to fill in the struct
		err := r(o, elementT.Type, elementV, jsonTag, size, depth+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func rSlice(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}
	/*
		// If you cant even set it dont even try
		if !v.CanSet() {
			return errors.New("cannot set slice")
		}

		// Prevent recursing deeper than configured levels
		if depth >= RecursiveDepth {
			return nil
		}


		// Grab original size to use if needed for sub arrays
		ogSize := size

		// If the value has a len and is less than the size
		// use that instead of the requested size
		elemLen := v.Len()
		if elemLen == 0 && size == -1 {
			size = number(o, 1, 10)
		} else if elemLen != 0 && (size == -1 || elemLen < size) {
			size = elemLen
		}

		// Get the element type
		elemT := t.Elem()

		// Loop through the elements length and set based upon the index
		for i := 0; i < size; i++ {
			nv := reflect.New(elemT)
			err := r(o, elemT, nv.Elem(), tag, ogSize, depth+1)
			if err != nil {
				return err
			}

			// If values are already set fill them up, otherwise append
			if elemLen != 0 {
				v.Index(i).Set(reflect.Indirect(ov))
			} else {
				v.Set(reflect.Append(reflect.Indirect(v), reflect.Indirect(ov)))
			}
		}
	*/

	return nil
}

func rMap(o *omit, t reflect.Type, v reflect.Value, tag string, size int, depth int) error {
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}
	/*

		// If you cant even set it dont even try
		if !v.CanSet() {
			return errors.New("cannot set slice")
		}

		// Prevent recursing deeper than configured levels
		if depth >= RecursiveDepth {
			return nil
		}

		// Check if tag exists, if so run custom function
		if tag != "" {
			return rCustom(o, v, tag)
		} else if isFakeable(t) && size <= 0 {
			// Only call custom function if no fakesize is specified (size <= 0)
			value, err := callFake(o, v, reflect.Map)
			if err != nil {
				return err
			}

			v.Set(reflect.ValueOf(value))
			return nil
		}

		// Set a size
		newSize := size
		if newSize == -1 {
			newSize = number(o, 1, 10)
		}

		// Create new map based upon map key value type
		mapType := reflect.MapOf(t.Key(), t.Elem())
		newMap := reflect.MakeMap(mapType)

		for i := 0; i < newSize; i++ {
			// Create new key
			mapIndex := reflect.New(t.Key())
			err := r(o, t.Key(), mapIndex.Elem(), "", -1, depth+1)
			if err != nil {
				return err
			}

			// Create new value
			mapValue := reflect.New(t.Elem())
			err = r(o, t.Elem(), mapValue.Elem(), "", -1, depth+1)
			if err != nil {
				return err
			}

			newMap.SetMapIndex(mapIndex.Elem(), mapValue.Elem())
		}

		// Set newMap into struct field
		if t.Kind() == reflect.Ptr {
			v.Set(oewMap.Elem())
		} else {
			v.Set(oewMap)
		}
	*/

	return nil
}

func rString(o *omit, t reflect.Type, v reflect.Value, tag string) error {
	if setOmitEmpty(o, t, v, tag) {
		return nil
	}
	/*
		if tag != "" {
			genStr, err := generate(o, tag)
			if err != nil {
				return err
			}

			v.SetString(genStr)
		} else if isFakeable(t) {
			value, err := callFake(o, v, reflect.String)
			if err != nil {
				return err
			}

			valueStr, ok := value.(string)
			if !ok {
				return errors.New("call to Fake method did not return a string")
			}
			v.SetString(valueStr)
		} else {
			genStr, err := generate(o, strings.Repeat("?", number(o, 4, 10)))
			if err != nil {
				return err
			}

			v.SetString(genStr)
		}
	*/

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
	/*
		if tag != "" {
			// Generate time
			timeOutput, err := generate(o, tag)
			if err != nil {
				return err
			}

			// Check to see if timeOutput has monotonic clock reading
			// if so, remove it. This is because time.Parse() does not
			// support parsing the monotonic clock reading
			if strings.Contains(timeOutput, " m=") {
				timeOutput = strings.Split(timeOutput, " m=")[0]
			}

			// Check to see if they are passing in a format	to parse the time
			timeFormat, timeFormatOK := t.Tag.Lookup("format")
			if timeFormatOK {
				timeFormat = javaDateTimeFormatToGolangFormat(timeFormat)
			} else {
				// If tag == "{date}" use time.RFC3339
				// They are attempting to use the default date lookup
				if tag == "{date}" {
					timeFormat = time.RFC3339
				} else {
					// Default format of time.Now().String()
					timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
				}
			}

			// If output is larger than format cut the output
			// This helps us avoid errors from time.Parse
			if len(timeOutput) > len(timeFormat) {
				timeOutput = timeOutput[:len(timeFormat)]
			}

			// Attempt to parse the time
			timeStruct, err := time.Parse(timeFormat, timeOutput)
			if err != nil {
				return err
			}

			v.Set(reflect.ValueOf(timeStruct))
			return nil
		}

		v.Set(reflect.ValueOf(date(o)))
	*/
}

func rCustom(o *omit, v reflect.Value, tag string) error {
	/*
		// If tag is empty return error
		if tag == "" {
			return errors.New("tag is empty")
		}

		fName, fParams := parseNameAndParamsFromTag(tag)
		info := GetFuncLookup(fName)

		// Check to see if it's a replaceable lookup function
		if info == nil {
			return fmt.Errorf("function %q not found", tag)
		}

		// Parse map params
		mapParams, err := parseMapParams(info, fParams)
		if err != nil {
			return err
		}

		// Call function
		fValue, err := info.Generate(f, mapParams, info)
		if err != nil {
			return err
		}

		// Create new element of expected type
		field := reflect.New(reflect.TypeOf(fValue))
		field.Elem().Set(reflect.ValueOf(fValue))

		// Check if element is pointer if so
		// grab the underlying value
		fieldElem := field.Elem()
		if fieldElem.Kind() == reflect.Ptr {
			fieldElem = fieldElem.Elem()
		}

		// Check if field kind is the same as the expected type
		if fieldElem.Kind() != v.Kind() {
			// return error saying the field and kinds that do not match
			return errors.New("field kind " + fieldElem.Kind().String() + " does not match expected kind " + v.Kind().String())
		}

		// Set the value
		v.Set(fieldElem.Convert(v.Type()))
	*/

	// If a function is called to set the struct
	// stop from going through sub fields
	return nil
}
