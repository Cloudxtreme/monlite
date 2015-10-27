// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Couch package provides a way to interact with Couch database.
package couch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/fcavani/e"
	utilUrl "github.com/fcavani/net/url"
)

const Err404 = "404 Object Not Found"

const ErrDbExist = "db exist"
const ErrDbNotFound = "database not found"
const ErrCantGetDoc = "can't get document"
const ErrDocConflict = "document conflict"
const ErrDocDontExist = "document don't exist"

var HttpClient *http.Client

func init() {
	if HttpClient == nil {
		HttpClient = &http.Client{}
	}
}

func authorizationHeader(userinfo string) string {
	enc := base64.URLEncoding
	encoded := make([]byte, enc.EncodedLen(len(userinfo)))
	enc.Encode(encoded, []byte(userinfo))
	return "Basic " + string(encoded)
}

type client struct {
	url *url.URL
}

func (c *client) do(method, path, params string, body []byte) (int, []byte, error) {
	url := utilUrl.Copy(c.url)
	url.Path += path
	url.RawQuery = params

	var buf io.Reader
	if len(body) > 0 {
		buf = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return 0, nil, e.New("can't create request")
	}

	if c.url.User != nil {
		req.Header.Set("Authorization", authorizationHeader(url.User.String()))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpClient.Do(req)
	if err != nil {
		return 0, nil, e.Push(err, "can't put")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return resp.StatusCode, nil, e.New("wrong status code - %v: %v", resp.StatusCode, resp.Status)
	}

	var data []byte
	data, _ = ioutil.ReadAll(resp.Body)

	return resp.StatusCode, data, nil
}

func (c *client) get(path, params string, body []byte) (int, []byte, error) {
	code, data, err := c.do("GET", path, params, body)
	return code, data, e.Forward(err)
}

func (c *client) put(path, params string, body []byte) (int, []byte, error) {
	code, data, err := c.do("PUT", path, params, body)
	return code, data, e.Forward(err)
}

func (c *client) delete(path, params string, body []byte) (int, []byte, error) {
	code, data, err := c.do("DELETE", path, params, body)
	return code, data, e.Forward(err)
}

// Couch object support the basic operations to the Couch data base.

type Ok struct {
	Ok bool `json:"ok"`
}

// CreateDB creates a new database. url is the url for the database and dbname is the name of
// the new database.
func CreateDB(url *url.URL, dbname string) error {
	c := &client{url}
	code, resp, err := c.put(dbname, "", nil)
	if err != nil && code == 412 {
		return e.New(ErrDbExist)
	} else if err != nil {
		return e.Push(err, "can't create the database")
	}
	if code != http.StatusCreated {
		return e.New("can't create the database, wrong response code: %v", code)
	}
	ok := new(Ok)
	err = json.Unmarshal(resp, ok)
	if err != nil {
		return e.Push(err, e.New("can't unserialize the response: %v", string(resp)))
	}
	if !ok.Ok {
		return e.New("can't create the database: wrong response")
	}
	return nil
}

// DeleteDB removes the database. url is the url for the database and dbname is the name of
// the new database.
func DeleteDB(url *url.URL, dbname string) error {
	c := &client{url}
	code, resp, err := c.delete(dbname, "", nil)
	if e.Contains(err, Err404) {
		return e.Push(err, ErrDbNotFound)
	} else if err != nil {
		return e.Push(err, "can't delete the database")
	}
	if code != http.StatusOK {
		return e.New("can't delete the database, wrong response code: %v", code)
	}
	ok := new(Ok)
	err = json.Unmarshal(resp, ok)
	if err != nil {
		return e.Push(err, e.New("can't unserialize the response: %v", string(resp)))
	}
	if !ok.Ok {
		return e.New("can't delete the database: wrong response")
	}
	return nil
}

// Struct representing a Couch database info.
type DatabaseInfo struct {
	Db_name              string `json:"db_name"`
	Doc_count            int    `json:"doc_count"`
	Doc_del_count        int64  `json:"doc_del_count"`
	Update_seq           int64  `json:"update_seq"`
	Purge_seq            int64  `json:"purge_seq"`
	Compact_running      bool   `json:"compact_running"`
	Disk_size            int64  `json:"disk_size"`
	Instance_start_time  string `json:"instance_start_time"`
	Disk_format_version  int    `json:"disk_format_version"`
	Committed_update_seq int64  `json:"committed_update_seq"`
}

// InfoDB reports the status of the database.
// url is the url for the database and dbname is the name of the new database.
func InfoDB(url *url.URL, dbname string) (*DatabaseInfo, error) {
	c := &client{url}
	code, data, err := c.get(dbname, "", nil)
	if e.Contains(err, Err404) {
		return nil, e.Push(err, ErrDbNotFound)
	} else if err != nil {
		return nil, e.Push(err, "can't get database information")
	}
	if code != http.StatusOK {
		return nil, e.New("can't get database information, wrong response code: %v", code)
	}
	di := new(DatabaseInfo)
	err = json.Unmarshal(data, di)
	if err != nil {
		return nil, e.Push(err, "can't get database information")
	}
	return di, nil
}

// Couch object support the basic operations to the Couch data base.
type Couch struct {
	*client
}

// NewCouch creates and initializes a new instance of the Couch struct.
func NewCouch(url *url.URL, dbname string) *Couch {
	url = utilUrl.Copy(url)
	url.Path += "/" + dbname + "/"
	return &Couch{
		client: &client{url},
	}
}

func grepId(data []byte) string {
	re, err := regexp.Compile("\"_id\":\"([^\"]*)")
	if err != nil {
		return ""
	}
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

type response struct {
	Ok  bool   `json:"ok"`
	Id  string `json:"id"`
	Rev string `json:"rev"`
}

// Put inserts in the database the struct in v, v must have the fild _id and
// the field _rev both strings. Like the TestStruct struct.
func (c *Couch) Put(v interface{}) (id, rev string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", "", e.Push(err, "can't serialize the data")
	}
	id = grepId(data)
	if id == "" {
		return "", "", e.New("document have no id")
	}
	code, resp, err := c.put(id, "", data)
	if err != nil && code == 409 {
		return "", "", e.Push(err, ErrDocConflict)
	} else if err != nil {
		return "", "", e.Push(err, "can't put the document")
	}
	if code == http.StatusCreated {
		ir := new(response)
		err := json.Unmarshal(resp, ir)
		if err != nil {
			return "", "", e.Push(e.Push(err, "can't desserialize returned data"), "can't put the document")
		}
		if !ir.Ok {
			return "", "", e.New("put failed")
		}
		return ir.Id, ir.Rev, nil
	}
	return "", "", e.Push(e.Push(err, e.New("received wrong response code, %v", code)), "can't put the document")
}

// Get the document id from the database and put it in
// the struct point by v. rev is the revision and is optional.
func (c *Couch) Get(id, rev string, v interface{}) error {
	if rev != "" {
		rev = "rev=" + rev
	}
	code, data, err := c.get(id, rev, nil)
	if err != nil {
		return e.Push(err, ErrCantGetDoc)
	}
	if code != http.StatusOK {
		return e.New("can't get the document, wrong code: %v", code)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return e.Push(err, "can't unserialize the document")
	}
	return nil
}

type idRev struct {
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
}

// Delete deletes from database the id id with rev revision.
// If rev is empty Delete deletes the last revision.
func (c *Couch) Delete(id, rev string) (string, error) {
	if rev == "" {
		revision := new(idRev)
		err := c.Get(id, "", revision)
		if e.Equal(err, ErrCantGetDoc) {
			return "", e.Push(err, ErrDocDontExist)
		} else if err != nil {
			return "", e.Push(err, "not deleted")
		}
		rev = revision.Rev
	}
	rev = "rev=" + rev
	code, resp, err := c.delete(id, rev, nil)
	if err != nil && code == 404 {
		return "", e.Push(err, ErrDocDontExist)
	} else if err != nil {
		return "", e.Push(err, "not deleted")
	}
	if code != http.StatusOK && code != http.StatusAccepted {
		return "", e.New("can't delete the document, wrong code: %v", code)
	}
	dr := new(response)
	err = json.Unmarshal(resp, dr)
	if err != nil {
		return "", e.Push(e.Push(err, "can't desserialize returned data"), "can't put the document")
	}
	if !dr.Ok {
		return "", e.New("delete failed")
	}
	if dr.Id != id {
		return "", e.New("response with the wrong id")
	}
	return dr.Rev, nil
}

// TestStruct is a sample of what the structs
// look like in this Couch package.
type TestStruct struct {
	Id   string `json:"_id"`
	Rev  string `json:"_rev,omitempty"`
	Data int
}
