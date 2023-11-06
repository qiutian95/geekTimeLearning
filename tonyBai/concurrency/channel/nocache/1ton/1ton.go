package main

import (
	"fmt"
	"sync"
	"time"
)

//1对n的信号通知

func worker(i int) {
	fmt.Printf("worker[%d] is working \n", i)
	time.Sleep(time.Second)
	fmt.Printf("worker[%d] finish work \n", i)
}

type signal struct {
}

func spawnGroup(f func(int), groupNum int, groupChan <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup
	for i := 0; i < groupNum; i++ {
		wg.Add(1)
		go func(a int) {
			<-groupChan // 所有goroutine中的接收channel都会收到信号，然后同时往下执行
			fmt.Printf("worker[%d] start to work \n", a)
			f(a)
			wg.Done()
		}(i)
	}

	go func() { // 这里为什么还要用一个goroutine包起来？ 如果这里不放在goroutine中那么，会导致main的goroutine阻塞，然后执行不了close，然后这边要等主goroutine执行，导致了死锁
		wg.Wait()
		c <- signal{}
	}()

	return c
}

func main() {
	println("program start")
	gc := make(chan signal)

	c := spawnGroup(worker, 5, gc)
	time.Sleep(5 * time.Second)
	close(gc)
	<-c
	fmt.Println("all finish")

}
