// Copyright 2014 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bytes"
	"encoding/json"
	"github.com/fsouza/go-dockerclient/engine"
	"io"
)

// Version returns version information about the docker server.
//
// See http://goo.gl/IqKNRE for more details.
func (c *Client) Version() (*engine.Env, error) {
	body, _, err := c.do("GET", "/version", nil)
	if err != nil {
		return nil, err
	}
	out := engine.NewOutput()
	remoteVersion, err := out.AddEnv()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	return remoteVersion, nil
}

type EventStreamMessage struct {
	Status string
	Id     string
	From   string
	Time   uint64
}

func (c *Client) Events(outchan chan EventStreamMessage) error {
	var outbuf bytes.Buffer
	go func(outchan chan EventStreamMessage) {
		c.stream("GET", "/events", nil, nil, &outbuf)
		dec := json.NewDecoder(&outbuf)
		for {
			var m EventStreamMessage
			if err := dec.Decode(&m); err == io.EOF {
				close(outchan)
				break
			} else if err != nil {
				break
			}
			outchan <- m
		}
	}(outchan)

	return nil
}

// Info returns system-wide information, like the number of running containers.
//
// See http://goo.gl/LOmySw for more details.
func (c *Client) Info() (*engine.Env, error) {
	body, _, err := c.do("GET", "/info", nil)
	if err != nil {
		return nil, err
	}
	var info engine.Env
	err = info.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return &info, nil
}
