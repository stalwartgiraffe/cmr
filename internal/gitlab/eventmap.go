package gitlab

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/TwiN/go-color"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/withstack"
	"golang.org/x/exp/maps"
)

type EventMap map[int]EventModel

func NewEventMapFromYaml(filepath string) (EventMap, error) {
	var oldEvents []EventModel

	if err := utils.ReadFromYamlFile(filepath, &oldEvents); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println(color.Ize(color.Red, err.Error()))
			return nil, err
		}
	}

	return NewEventMapFromSlice(oldEvents), nil
}

func NewEventMapFromSlice(events []EventModel) EventMap {
	m := make(EventMap)
	if events != nil {
		for _, e := range events {
			m[e.ID] = e
		}
	}
	return m
}

func (m EventMap) Insert(events EventMap) EventMap {
	for k, v := range events {
		m[k] = v
	}
	return m
}

func (m EventMap) WriteToJsonFile(filepath string) error {
	events := maps.Values(m)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID > events[j].ID
	})

	file, err := os.Create(filepath)
	if err != nil {
		return withstack.Errorf("%w", err)
	}
	defer file.Close()

	ss := EventModelSlice(events)
	if bb, err := ss.MarshalJSON(); err != nil {
		return err
	} else {
		const permissions os.FileMode = 0644
		return os.WriteFile(filepath, bb, permissions)
	}
}
func (m EventMap) WriteToYamlFile(filepath string) error {
	events := maps.Values(m)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID > events[j].ID
	})
	return utils.WriteToYamlFile(filepath, events)
}

func (m EventMap) LastDate() string {
	if m == nil || len(m) < 1 {
		return ""
	}
	last := Time{}
	for _, v := range m {
		if !v.CreatedAt.IsZero() {
			if v.CreatedAt.Time.After(last.Time) {
				last = v.CreatedAt
			}
		}
	}
	return last.Format("2006-01-02")
}

func (m EventMap) ProjectIDs() []int {
	d := make(map[int]struct{})
	for _, v := range m {
		d[v.ProjectID] = struct{}{}
	}
	ids := sort.IntSlice{}
	for k := range d {
		ids = append(ids, k)
	}
	ids.Sort()
	return ids
}

// consider this but seems a lot of work
// https://github.com/abice/go-enum
/*
  DiffNote
  DiscussionNote
  MergeRequest
  Note
*/
/*

  accepted
  approved
  closed
  commented on
  created
  deleted
  joined
  opened
  pushed new
  pushed to
*/
func (m EventMap) GetTargetAction(actionName string, targetType string) []EventModel {
	events := []EventModel{}
	for _, v := range m {
		fmt.Printf("#%s#%s#\n", v.TargetType, v.ActionName)
		if v.TargetType == targetType && v.ActionName == actionName {
			events = append(events, v)
		}
	}
	return events
}
