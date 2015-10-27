// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package monlite

import (
	"time"

	"github.com/fcavani/e"
	"github.com/fcavani/log"
	"github.com/fcavani/ping"
)

type Monitor struct {
	Name    string
	Url     string
	Timeout time.Duration
	Periode time.Duration
	Sleep   time.Duration
	Fails   int
	OnFail  func(m *Monitor) error
	chclose chan chan struct{}
	count   int
}

func (m *Monitor) Start() error {
	if m.Name == "" {
		return e.New("empty name")
	}
	if m.Url == "" {
		return e.New("empty url")
	}
	if m.Periode == 0 {
		return e.New("periode must be greater than zero")
	}
	m.chclose = make(chan chan struct{})
	go func() {
		for {
			select {
			case ch := <-m.chclose:
				ch <- struct{}{}
				return
			case <-time.After(m.Periode):
				resp := make(chan error)
				go func() {
					log.DebugLevel().Printf("Pinging %v", m.Name)
					start := time.Now()
					err := ping.PingRawUrl(m.Url)
					if err != nil {
						log.Errorf("Ping failed for %v with error: %v", m.Name, err)
						resp <- e.Forward(err)
						return
					}
					log.DebugLevel().Printf("Ping ok for %v (%v)", m.Name, time.Since(start))
					resp <- nil
				}()
				select {
				case err := <-resp:
					if err == nil {
						continue
					}
					m.count++
					if m.Fails != 0 || m.Fails <= m.count {
						m.count = 0
						if m.OnFail == nil {
							continue
						}
						err := m.OnFail(m)
						if err != nil {
							log.Errorf("Onfail function on %v returned an error: %v", m.Name, err)
						}
						log.Printf("%v going to sleep for %v", m.Name, m.Sleep)
						select {
						case <-time.After(m.Sleep):
						case ch := <-m.chclose:
							ch <- struct{}{}
							return
						}
					}
					continue
				case <-time.After(m.Timeout):
					log.Errorf("Ping timeout for %v", m.Name)
					m.count++
					if m.Fails != 0 || m.Fails <= m.count {
						m.count = 0
						if m.OnFail == nil {
							continue
						}
						err := m.OnFail(m)
						if err != nil {
							log.Errorf("Onfail function on %v returned an error: %v", m.Name, err)
						}
						log.Printf("%v going to sleep for %v", m.Name, m.Sleep)
						select {
						case <-time.After(m.Sleep):
						case ch := <-m.chclose:
							ch <- struct{}{}
							return
						}
					}
				}
			}
		}
	}()
	return nil
}

func (m *Monitor) Stop() error {
	ch := make(chan struct{})
	m.chclose <- ch
	<-ch
	return nil
}
