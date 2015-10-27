// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package couch

import (
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/fcavani/e"
)

const CouchUrlTravis = "http://localhost:5984"

var CouchUrl string = "http://site2:1234567site2@localhost:5984"

var couchUrl *url.URL
var dbName = "couchtestdatabase"

func TestMain(m *testing.M) {
	if os.Getenv("TRAVIS_BUILD_DIR") != "" {
		CouchUrl = CouchUrlTravis
	}
	os.Exit(m.Run())
}

func TestUrl(t *testing.T) {
	var err error
	couchUrl, err = url.Parse(CouchUrl)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDB(t *testing.T) {
	err := CreateDB(couchUrl, dbName)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	defer DeleteDB(couchUrl, dbName)

	err = CreateDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}

	info, err := InfoDB(couchUrl, dbName)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	if info.Db_name != dbName || info.Doc_count != 0 {
		t.Fatal("samething wrong with the db")
	}

	err = DeleteDB(couchUrl, dbName)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}

	_, err = InfoDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbNotFound) {
		t.Fatal(e.Trace(e.Forward(err)))
	}

	err = DeleteDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbNotFound) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

func TestPut(t *testing.T) {
	couch := NewCouch(couchUrl, dbName)
	err := CreateDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	defer DeleteDB(couchUrl, dbName)
	for i := 0; i < 10; i++ {
		test := TestStruct{
			Id:   strconv.FormatInt(int64(i), 10),
			Data: i,
		}
		id, _, err := couch.Put(test)
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		if id != test.Id {
			t.Fatal("wrong id", id)
		}
	}
	di, err := InfoDB(couchUrl, dbName)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	if di.Db_name != dbName {
		t.Fatal("wrong db name")
	}
	if di.Doc_count != 10 {
		t.Fatal("wrong document count", di.Doc_count)
	}
	test := TestStruct{
		Id:   "1",
		Data: 1,
	}
	_, _, err = couch.Put(test)
	if err != nil && !e.Equal(err, ErrDocConflict) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

func TestGet(t *testing.T) {
	couch := NewCouch(couchUrl, dbName)
	err := CreateDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	defer DeleteDB(couchUrl, dbName)
	for i := 0; i < 10; i++ {
		test := TestStruct{
			Id:   strconv.FormatInt(int64(i), 10),
			Data: i,
		}
		id, _, err := couch.Put(test)
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		if id != test.Id {
			t.Fatal("wrong id", id)
		}
	}
	for i := 0; i < 10; i++ {
		test := new(TestStruct)
		err := couch.Get(strconv.FormatInt(int64(i), 10), "", test)
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		if test.Data != i {
			t.Fatal("retrieved wrong document", test.Data)
		}
	}
	test := new(TestStruct)
	err = couch.Get("catoto", "", test)
	if err != nil && !e.Equal(err, ErrCantGetDoc) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

func TestEdit(t *testing.T) {
	couch := NewCouch(couchUrl, dbName)
	err := CreateDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	defer DeleteDB(couchUrl, dbName)

	test := TestStruct{
		Id:   "1",
		Data: 1,
	}
	id, rev, err := couch.Put(test)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	t.Log(rev)

	test2 := new(TestStruct)
	err = couch.Get(id, "", test2)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	if test2.Id != id || test2.Data != 1 {
		t.Fatal("wrong test data")
	}

	test2.Data = 42
	id, rev, err = couch.Put(test2)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	t.Log(rev)

	test3 := new(TestStruct)
	err = couch.Get(id, rev, test3)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	if test3.Id != id || test3.Data != 42 {
		t.Fatal("wrong test data")
	}
}

func TestDelete(t *testing.T) {
	couch := NewCouch(couchUrl, dbName)
	err := CreateDB(couchUrl, dbName)
	if err != nil && !e.Equal(err, ErrDbExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	defer DeleteDB(couchUrl, dbName)
	test := TestStruct{
		Id:   "1",
		Data: 1,
	}
	_, _, err = couch.Put(test)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	_, err = couch.Delete(test.Id, "")
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	_, err = couch.Delete(test.Id, "")
	if err != nil && !e.Equal(err, ErrDocDontExist) {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}
