package main

import (
	"sync"
	"time"
)

type combine struct {
	sync.Mutex

	parameters []string
	excutor    func([]string) error

	ticker <-chan time.Time
}

func new(f func([]string) error) *combine {
	return &combine{excutor: f}
}

func (c *combine) push(p []string) {
	c.Lock()
	defer c.Unlock()
	if c.ticker == nil {
		c.parameters = p
		c.ticker = time.Tick(time.Second * 5)
		go func() {
			<-c.ticker
			c.ticker = nil
			c.Lock()
			if c.excutor(c.parameters) == nil {
				c.parameters = nil
			}
			c.Unlock()
		}()
	} else {
		c.parameters = append(c.parameters, p...)
	}
}
