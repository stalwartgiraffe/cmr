package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v2"

	"github.com/stalwartgiraffe/cmr/withstack"
)

func YamlString(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func WriteToYamlFile[T any](path string, t T) error {
	file, err := os.Create(path)
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(t)
}

func WriteStringToFile(filepath string, content string) error {
    return os.WriteFile(filepath, []byte(content), 0644)
}

func PrettyJSON(data []byte) (string, error) {
	var buf bytes.Buffer
	err := json.Indent(&buf, data, "", "  ")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ReadFromYamlFile[T any](path string, t *T) error {
	file, err := os.Open(path)
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(t); err != nil {
		return withstack.Errorf("%w", err)
	}
	return nil
}

func ToSortedSlice[T any](m map[int]T) []T {
	keys := maps.Keys(m)
	sort.Ints(keys)
	s := make([]T, len(m))
	for i, k := range keys {
		s[i] = m[k]
	}
	return s
}

func Join(strs ...string) string {
	return strings.Join(strs, " ")
}
