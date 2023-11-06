package main

import (
	"fmt"
	"sync"
)

// 使用无缓冲的channel实现锁机制

// 传统的锁实现
/*type counter struct {
	sync.Mutex
	c int
}

var cter counter

func increase() {
	cter.Lock()
	defer cter.Unlock()
	cter.c++
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			increase()
			fmt.Printf("cter.c is %d\n", cter.c)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("finish")
}*/

// 无缓冲channel实现

// 1.我的实现 在传统的基础上改的，有问题，会出现锁定锁不住
/*var count int

func increase(lock <-chan int) {
	go func() {
		<-lock
		count++
	}()
}

func main() {
	var wg sync.WaitGroup
	lock := make(chan int)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		increase(lock)
		go func() {
			lock <- i
			fmt.Printf("counter is %d\n", count)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("finish")
}*/

// 2.tony实现
type counter struct {
	c chan int
	i int
}

func newCounter() *counter {
	cter := &counter{
		c: make(chan int),
		i: 0,
	}
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	return cter
}

func (cter *counter) increase() {
	<-cter.c
}
func main() {
	cter := newCounter()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			cter.increase()
			fmt.Println("count is ", cter.i)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("finish")
}
