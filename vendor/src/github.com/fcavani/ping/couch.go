// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net/url"
	"strconv"

	"github.com/fcavani/couch"
	"github.com/fcavani/e"
	utilUrl "github.com/fcavani/net/url"
)

// DbName is the name of the data used by PingCouch
var DbName = "pingmonitortest"

// PingCouch tests if the database is online and operational.
func PingCouch(url *url.URL) error {
	url = utilUrl.Copy(url)
	url.Scheme = "http"

	err := couch.CreateDB(url, DbName)
	if err != nil && !e.Equal(err, couch.ErrDbExist) {
		return e.Forward(err)
	}
	defer couch.DeleteDB(url, DbName)

	c := couch.NewCouch(url, DbName)

	for i := 0; i < 10; i++ {
		t := couch.TestStruct{
			Id:   strconv.FormatInt(int64(i), 10),
			Data: i,
		}
		id, _, err := c.Put(t)
		if err != nil {
			return e.Forward(err)
		}
		if id != t.Id {
			return e.New("wrong id (%v)", id)
		}
	}
	di, err := couch.InfoDB(url, DbName)
	if err != nil {
		return e.Forward(err)
	}
	if di.Db_name != DbName {
		return e.New("wrong db name")
	}
	if di.Doc_count != 10 {
		return e.New("wrong document count (%v)", di.Doc_count)
	}
	for i := 0; i < 10; i++ {
		t := new(couch.TestStruct)
		err := c.Get(strconv.FormatInt(int64(i), 10), "", t)
		if err != nil {
			return e.Forward(err)
		}
		if t.Data != i {
			return e.New("retrieved wrong document (%v)", t.Data)
		}
	}
	_, err = c.Delete("9", "")
	if err != nil {
		return e.Forward(err)
	}
	t := new(couch.TestStruct)
	err = c.Get("9", "", t)
	if err != nil && !e.Equal(err, couch.ErrCantGetDoc) {
		return e.Forward(err)
	}
	err = couch.DeleteDB(url, DbName)
	if err != nil {
		return e.Forward(err)
	}
	_, err = couch.InfoDB(url, DbName)
	if err != nil && !e.Equal(err, couch.ErrDbNotFound) {
		return e.Push(err, "database not deleted")
	}
	return nil
}

func init() {
	Add("couch", PingCouch)
}
