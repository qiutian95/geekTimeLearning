package main

import (
	"fmt"
	"runtime"
	"time"
)

func deadloop() {
	for {
		//fmt.Println("I am deadloop")
	}
}

func main() {
	runtime.GOMAXPROCS(1) // 加上这个还是能运行，因为在go1.14之后增加了非协作式抢占，
	go deadloop()

	for {
		time.Sleep(time.Second + 1)
		fmt.Println("I get scheduled!")
	}
}
