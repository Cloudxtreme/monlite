// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// license that can be found in the LICENSE file.

package ping

import (
	"net/url"
	"time"

	"github.com/fcavani/e"

	"gopkg.in/mgo.v2"
)

var MongoDbTimeout time.Duration = 60 * time.Second
var Tryies int = 10

func PingMongoDb(u *url.URL) error {
	session, err := mgo.DialWithTimeout(u.String(), MongoDbTimeout)
	if err != nil {
		return e.New(err)
	}
	defer session.Close()
	for i := 0; i < Tryies; i++ {
		err := session.Ping()
		if err != nil {
			return e.New(err)
		}
	}
	return nil
}

func init() {
	Add("mongodb", PingMongoDb)
}
