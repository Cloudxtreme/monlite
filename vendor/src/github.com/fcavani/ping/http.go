// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"io"
	"net/http"
	"net/url"

	"github.com/fcavani/e"
	"github.com/fcavani/log"
)

// PingHttp connect a http or https server and try to
// receive something. If the server return a code different
// of 2xx, it will fail. Ignores insecure certificates.
func PingHttp(url *url.URL) error {
	resp, err := httpClient.Get(url.String())
	if e.Contains(err, "connection refused") {
		return e.Push(e.New(err), "get failed: connection refused")
	} else if err != nil {
		return e.Push(e.New(err), "get failed")
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return e.Forward(err)
		}
		buf = buf[:n]
		//status.Log(status.Protocol, "PingHttp status code is %v and received it from server: %v", resp.StatusCode, string(buf))
		log.ProtoLevel().Printf("PingHttp status code is %v and received it from server: %v", resp.StatusCode, string(buf))
		return e.New("returned status code %v, expected 2xx", resp.StatusCode)
	}
	_, err = resp.Body.Read(buf)
	if err != nil && err != io.EOF {
		return e.Forward(err)
	}
	return nil
}

func init() {
	Add("http", PingHttp)
	Add("https", PingHttp)
}
