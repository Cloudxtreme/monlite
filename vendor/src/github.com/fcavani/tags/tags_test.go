// Copyright 2009 Felipe Alves Cavani. All rights reserved.
// Start date:        2011-04-18
// Last modification: 2011-

package tags

import (
	"testing"

	"github.com/fcavani/e"
)

type tests struct {
	rawtags string
	result  string
}

var Tests []tests = []tests{
	{"", ""},
	{"tag1", "tag1"},
	{"tag1, tag2", "tag1, tag2"},
	{"tag3, tag1, tag2", "tag1, tag2, tag3"},
	{"tag1, tag2, tag2", "tag1, tag2"},
	{"tag1,", "tag1"},
	{",tag1", "tag1"},
	{",tag1,", "tag1"},
	{",tag1, tag2", "tag1, tag2"},
	{",", ""},
}

var tags []Tags

func TestNewTags(t *testing.T) {
	const test = "a, b, c, d"
	tags, err := NewTags(test)
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	if tags.Get() != test {
		t.Fatal("tags doesn't match")
	}
}

func TestAdd(t *testing.T) {
	tags := &Tags{}
	tags.Add("zz z")
	tags.Add("cc c")
	tags.Add("ii i")
	tags.Add("aa a")
	str := tags.Get()
	if str != "aa a, cc c, ii i, zz z" {
		t.Fatal("Tags don't match.", str)
	}
}

func TestReplace(t *testing.T) {
	tags = make([]Tags, len(Tests))
	for i, test := range Tests {
		err := tags[i].Replace(test.rawtags)
		if err != nil {
			t.Fatal("Replace failed:", err)
		}
	}
	//fmt.Printf("%v\n", tags)
	for i, tag := range tags {
		//fmt.Printf("%v\n", tag)
		if tag.String() != Tests[i].result {
			t.Fatalf("raw string not match: '%v' != '%v'", tag.String(), Tests[i].result)
		}
	}
}

func TestCopy(t *testing.T) {
	tags2 := make([]*Tags, len(Tests))
	for i, tag := range tags {
		tags2[i] = tag.Copy()
		//println(tag.String(), "-", tags2[i].String())
		if tag.String() != tags2[i].String() {
			t.Fatal("copy failed:", tag.String(), "!=", tags2[i].String())
		}
	}
}

func TestFind(t *testing.T) {
	s := "users, admin, printer, foo, bar"
	tags := &Tags{}
	err := tags.Replace(s)
	if err != nil {
		t.Fatalf("Replace failed:", err)
	}
	if !tags.Exist("bar") {
		t.Fatal("Can't find tag.")
	}
	if !tags.Exist("foo") {
		t.Fatal("Can't find tag.")
	}
	if !tags.Exist("printer") {
		t.Fatal("Can't find tag.")
	}
	if tags.Exist("help") {
		t.Fatal("Find something strange.")
	}
}

func newtags(s string, t *testing.T) *Tags {
	tags := &Tags{}
	err := tags.Replace(s)
	if err != nil {
		t.Fatalf("Replace failed:", err)
	}
	return tags
}

func TestRemove(t *testing.T) {
	str := "tag1, tag2, tag3, tag4"
	tags := newtags(str, t)
	tags.Remove("tag1")
	if tags.Get() != "tag2, tag3, tag4" {
		t.Fatal("Remove failed")
	}
	tags = newtags(str, t)
	tags.Remove("tag2")
	if tags.Get() != "tag1, tag3, tag4" {
		t.Fatal("Remove failed")
	}
	tags = newtags(str, t)
	tags.Remove("tag3")
	if tags.Get() != "tag1, tag2, tag4" {
		t.Fatal("Remove failed")
	}
	tags = newtags(str, t)
	tags.Remove("tag4")
	if tags.Get() != "tag1, tag2, tag3" {
		t.Fatal("Remove failed")
	}
	tags.Remove("tag3")
	if tags.Get() != "tag1, tag2" {
		t.Fatal("Remove failed")
	}
	tags.Remove("tag2")
	if tags.Get() != "tag1" {
		t.Fatal("Remove failed")
	}
	tags.Remove("tag1")
	if tags.Get() != "" {
		t.Fatal("Remove failed")
	}
}

func TestEmpty(t *testing.T) {
	tags, err := NewTags("")
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	if tags.Get() != "" {
		t.Fatal("Get not empty.", err)
	}
	tagsCopy := tags.Copy()
	if tagsCopy.Get() != "" {
		t.Fatal("Get not empty.", err)
	}
	if err := tags.Remove("blá"); err == nil {
		t.Fatal("Tag removed from empty tags.")
	}
	if tags.Exist("blá") {
		t.Fatal("Tag blá exist in a empty tags.")
	}

	if err := tags.Add("blá"); err != nil {
		t.Fatal("Add failed:", err)
	}
	if !tags.Exist("blá") {
		t.Fatal("Tag blá don't exist.")
	}
	if err := tags.Remove("blá"); err != nil {
		t.Fatal("Remove failed:", err)
	}
}

func TestCompareAndSame(t *testing.T) {
	tags0, err := NewTags("")
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	tags1, err := NewTags("")
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	if count := tags0.Compare(tags1); count != 0 {
		t.Fatal("Compare failed.", count)
	}
	tags2, err := NewTags("aaa")
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	if count := tags2.Compare(tags1); count != 0 {
		t.Fatal("Compare failed.", count)
	}
	tags3, err := NewTags("aaa, 0")
	if err != nil {
		t.Fatal("NewTags failed:", err)
	}
	if count := tags2.Compare(tags3); count != 1 {
		t.Fatal("Compare failed.", count)
	}
	if !tags0.Same(tags1) {
		t.Fatal("Same failed.")
	}
	tags4 := tags3.Copy()
	if !tags3.Same(tags4) {
		t.Fatal("Same failed.")
	}
}

func TestLastComma(t *testing.T) {
	str := "tag1, tag2, tag3, "
	tags, err := NewTags(str)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	if tags.Get() != "tag1, tag2, tag3" {
		t.Fatal("failed")
	}
}

func TestMerge(t *testing.T) {
	str := "tag1, tag2, tag3"
	str2 := "tag8, tag9, tag4"
	tags, err := NewTags(str)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	tags2, err := NewTags(str2)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	err = tags.Merge(tags2)
	if tags.Get() != "tag1, tag2, tag3, tag4, tag8, tag9" {
		t.Fatal("failed")
	}
}

func TestMergeFromStringSlice(t *testing.T) {
	str := "tag1, tag2, tag3"
	tags, err := NewTags(str)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	err = tags.MergeFromStringSlice([]string{"tag8", "tag9", "tag4"})
	if tags.Get() != "tag1, tag2, tag3, tag4, tag8, tag9" {
		t.Fatal("failed")
	}
}
