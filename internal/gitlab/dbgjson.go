package gitlab

import (
	"reflect"
	"strings"
)

// getJsonKeys returns the set of expected json keys
func getJsonKeys(s any) map[string]struct{} {
	fieldNames := make(map[string]struct{})
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			tags := strings.Split(jsonTag, `,`)
			if 0 < len(tags) {
				fieldNames[tags[0]] = struct{}{}
			} else {
				fieldNames[jsonTag] = struct{}{}
			}
		} else {
			fieldNames[field.Name] = struct{}{}
		}
	}
	return fieldNames
}
