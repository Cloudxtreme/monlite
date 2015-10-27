// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// license that can be found in the LICENSE file.

// Tags packages have simple data structures and functions to
// store and manipulate a list of tags.
package tags

import (
	"sort"
	"strings"
	"unicode"

	"github.com/fcavani/e"
	u "github.com/fcavani/unicode"
	"github.com/fcavani/types"
)

var TagStringMax = 3000
var TagStringMin = 1

const ErrTagExist = "tag exist"
const ErrTagDoesntExit = "tag doesn't exist"
const ErrTagNotFound = "tag not found"
const ErrNoTags = "no tags"
const ErrInvalidTagsLength = "list of tags is big"

const ErrInvalidTagChar = "invalid tag character"

type Tags []string

func init() {
	types.Insert(&Tags{})
}

// NewTags creates a new tag list from a comma separeted list of tags.
func NewTags(tags string) (*Tags, error) {
	if err := CheckTags(tags); err != nil {
		return nil, e.Forward(err)
	}
	t := &Tags{}
	if tags != "" {
		splitted := strings.Split(tags, ",")
		for _, tag := range splitted {
			trim := strings.Trim(tag, " ")
			if trim == "" {
				continue
			}
			err := t.Add(tag)
			if err != nil {
				return t, e.Forward(err)
			}
		}
	}
	return t, nil
}

// Get all tags
func (t *Tags) Get() string {
	if t == nil {
		return ""
	}
	tags := *t
	s := ""
	lentags := len(tags)
	// TODO: Usar range
	for i := 0; i < lentags; i++ {
		s = s + tags[i]
		if i < lentags-1 {
			s = s + ", "
		}
	}
	return s
}

// String givers tags string formated
func (t *Tags) String() string {
	return t.Get()
}

// Copy makes a copy
func (t *Tags) Copy() *Tags {
	if t == nil {
		return nil
	}
	dst := make(Tags, len(*t))
	copy(dst, *t)
	return &dst
}

// Add one tag to the list
func (t *Tags) Add(tag string) error {
	ts := *t
	trim := strings.TrimSpace(tag)
	if trim == "" {
		return e.New(e.ErrEmptyString)
	}
	if err := CheckTag(trim); err != nil {
		return e.Forward(err)
	}
	i := sort.SearchStrings(ts, trim)
	if i >= ts.Len() || ts[i] != trim { // Not found
		//ts = append(ts, trim)
		//sort.Strings(ts)
		ts = append(ts[:i], append(Tags{trim}, ts[i:]...)...)
	} else {
		return e.New(ErrTagExist)
	}
	*t = ts
	return nil
}

// Replace all tags with new tags in the string, that is formated in a comma separated list.
func (t *Tags) Replace(s string) error {
	if err := CheckTags(s); err != nil {
		return e.Forward(err)
	}
	*t = Tags{}
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		trim := strings.TrimSpace(tag)
		if trim == "" {
			continue
		}
		err := t.Add(tag)
		if err != nil {
			continue
		}
	}
	return nil
}

// Merge merges two tags sets.
func (t *Tags) Merge(tags *Tags) error {
	a := *tags
	for _, tag := range a {
		err := t.Add(tag)
		if err != nil && !e.Equal(err, ErrTagExist) {
			return e.Forward(err)
		}
	}
	return nil
}

// MergeFromStringSlice merges a slice, unordered and unchecked, into a tags.
// Check all slice itens in the process.
func (t *Tags) MergeFromStringSlice(slice []string) error {
	for _, s := range slice {
		err := t.Add(s)
		if err != nil && !e.Equal(err, ErrTagExist) {
			return e.Forward(err)
		}
	}
	return nil
}

// Remove a tag from the list.
func (t *Tags) Remove(tag string) error {
	ts := *t
	i := sort.SearchStrings(ts, tag)
	if i < len(ts) && ts[i] == tag {
		if i+1 == len(ts) {
			ts = ts[:i]
		} else {
			ts = append(ts[:i], ts[i+1:]...)
		}
		*t = ts
		return nil
	}
	return e.New(ErrTagNotFound)
}

// Exist verifies if a tag exist.
func (t *Tags) Exist(target string) bool {
	ts := *t
	i := sort.SearchStrings(ts, target)
	if i < len(ts) && ts[i] == target {
		return true
	}
	return false
}

// Compare compares two sets of tags and return the number of matchs.
func (t *Tags) Compare(tr *Tags) int {
	tags := *t
	count := 0
	for _, tag := range tags {
		if tr.Exist(tag) {
			count++
		}
	}
	return count
}

// Same verifies if the tags sets are the same.
func (t *Tags) Same(tr *Tags) bool {
	tags1 := *t
	tags2 := *tr
	if len(tags1) != len(tags2) {
		return false
	}
	for i, tag := range tags1 {
		if tag != tags2[i] {
			return false
		}
	}
	return true
}

// Return the number of tags.
func (t *Tags) Len() int {
	if t == nil {
		return 0
	}
	return len(*t)
}

// Check if one that is less than other.
func (t *Tags) Less(i, j int) bool {
	a := *t
	return a[i] < a[j]
}

// Change the place of two tags.
func (t *Tags) Swap(i, j int) {
	a := *t
	a[i], a[j] = a[j], a[i]
	*t = a
}

// Check if a set of tags are correctly formated.
func CheckTags(tags string) error {
	if len(tags) <= TagStringMin && len(tags) >= TagStringMax {
		return e.New(ErrInvalidTagsLength)
	}
	for _, v := range tags {
		if !u.IsLetter(v) && !unicode.IsDigit(v) && v != ',' && v != ' ' && v != '-' && v != '_' && v != '.' && v != '/' && v != ':' {
			return e.Push(e.New(ErrInvalidTagChar), e.New("the character '%v' is invalid", string([]byte{byte(v)})))
		}
	}
	return nil
}

// Check if a tag is formed with the corrects characters
func CheckTag(tag string) error {
	for _, v := range tag {
		if !u.IsLetter(v) && !unicode.IsDigit(v) && v != ' ' && v != '-' && v != '_' && v != '.' && v != '/' && v != ':' {
			return e.Push(e.New(ErrInvalidTagChar), e.New("the character '%v' is invalid", string([]byte{byte(v)})))
		}
	}
	return nil
}
