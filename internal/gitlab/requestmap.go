package gitlab

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/TwiN/go-color"
	easyjson "github.com/mailru/easyjson"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/withstack"
	"golang.org/x/exp/maps"
)

type MergeRequestMap map[int]MergeRequestModel

func NewMergeRequestMapFromYaml(filepath string) (MergeRequestMap, error) {
	var requests []MergeRequestModel
	if err := utils.ReadFromYamlFile(filepath, &requests); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println(color.Ize(color.Red, err.Error()))
			return nil, err
		}
	}

	return NewMergeRequestMapFromSlice(requests), nil
}

func NewMergeRequestMapFromSlice(requests []MergeRequestModel) MergeRequestMap {
	m := make(MergeRequestMap)
	for _, mr := range requests {
		m[mr.ID] = mr
	}
	return m
}

func (m MergeRequestMap) Insert(requests MergeRequestMap) MergeRequestMap {
	for k, v := range requests {
		m[k] = v
	}
	return m
}

func (m MergeRequestMap) LastCreatedDate() string {
	if m == nil || len(m) < 1 {
		return ""
	}
	last := Time{}
	for _, v := range m {

		if v.CreatedAt.Time.After(last.Time) {
			last = v.CreatedAt
		}
	}
	return last.Format("2006-01-02")
}

// easyjson:json
type MergeRequestSlice []MergeRequestModel

func (m MergeRequestMap) WriteToYamlFile(filepath string) error {
	requests := maps.Values(m)
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].ID > requests[j].ID
	})

	file, err := os.Create(filepath)
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	defer file.Close()

	ss := MergeRequestSlice(requests)
	if bb, err := easyjson.Marshal(ss); err != nil {
		return err
	} else {
		const permissions os.FileMode = 0644
		return os.WriteFile(filepath, bb, permissions)
	}
}
