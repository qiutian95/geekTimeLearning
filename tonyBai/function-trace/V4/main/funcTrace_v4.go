package main

import (
	"funcTraceV4/trace"
	"sync"
	"time"
)

// 增加自动注入Trace

func A1() {
	defer trace.Trace()()
	B1()
}

func B1() {
	defer trace.Trace()()
	C1()
}

func C1() {
	defer trace.Trace()()
}

func A2() {
	defer trace.Trace()()
	B2()
}
func B2() {
	defer trace.Trace()()
	C2()
}
func C2() {
	defer trace.Trace()()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A2()
		wg.Done()
	}()

	time.Sleep(time.Millisecond) // 让两个goroutine分开打印
	A1()
	wg.Wait()
}
