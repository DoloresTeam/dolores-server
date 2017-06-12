package main

import (
	"testing"
	"time"
)

func TestC(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	c := new(func(arg1 []string) error {
		t.Log(arg1)
		return nil
	})

	c.push([]string{`1`})
	time.Sleep(time.Second * 1)
	c.push([]string{`2`})
	time.Sleep(time.Second * 2)
	c.push([]string{`3`})
	time.Sleep(time.Second * 1)
	c.push([]string{`4`})
	time.Sleep(time.Second * 3)
	c.push([]string{`5`})
	time.Sleep(time.Second * 1)
	c.push([]string{`6`})
	time.Sleep(time.Second * 4)
	c.push([]string{`7`})
	time.Sleep(time.Second * 1)
	c.push([]string{`8`})
	time.Sleep(time.Second * 2)
	c.push([]string{`9`})
	time.Sleep(time.Second * 1)
	c.push([]string{`10`})
	time.Sleep(time.Second * 3)
	c.push([]string{`11`})

	time.Sleep(time.Second * 5)
}
