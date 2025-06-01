package withstack

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"unsafe"

	"github.com/gin-gonic/gin"
)

// DumpContextInternals writes the contents of the context w.
// This will recurse on wrapped context.
func DumpContextInternals(ctx interface{}, inner bool, w io.Writer) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Fprintf(w, "\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				DumpContextInternals(reflectValue.Interface(), true, w)
			} else {
				fmt.Fprintf(w, "field name: %+v\n", reflectField.Name)
				fmt.Fprintf(w, "value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Fprintf(w, "context is empty (int)\n")
	}
}

// DumpRoutes dumps the routes of a gin engine to a slice of strings.
func DumpRoutes(engine *gin.Engine) []string {
	r := []string{}
	for _, route := range engine.Routes() {
		r = append(r, (fmt.Sprintf("%s %s", route.Path, route.Method)))

	}
	sort.Strings(r)
	return r
}
