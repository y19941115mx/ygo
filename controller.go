package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"ygo/framework"
)

func FooController(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panic := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))

	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panic <- p
			}
		}()
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")
		finish <- struct{}{}
	}()

	select {
	case p := <-panic:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}
