package main

import (
	"fmt"
	"time"
)

// 无缓冲的1:1的channel

type signal struct {
}

func worker() {
	fmt.Println("now is working...")
	time.Sleep(time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		fmt.Println("start work!")
		f()
		c <- signal{}
	}()
	return c
}

func main() {
	c := make(<-chan signal)

	c = spawn(worker)
	<-c
	fmt.Println("end work!")
}
