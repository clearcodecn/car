package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Second)

		close(ch)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		<-ch
		fmt.Println(1)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		<-ch
		fmt.Println(2)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		select {
		case <-ch:
			fmt.Println(3)
		}
	}()

	time.Sleep(1 * time.Hour)
}
