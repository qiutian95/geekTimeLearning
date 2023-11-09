package main

import (
	"fmt"
	"sync"
	"time"
)

// 发送或接收的channel用在函数或方法参数列表或返回值列表中

func produce(c chan<- int) {
	for i := 0; i <= 9; i++ {
		c <- i
		time.Sleep(time.Second)
	}
	close(c)
}

func consume(c <-chan int) {
	for v := range c {
		fmt.Println("consume:", v)
	}
}

func main() {
	c := make(chan int, 5) // 容量是5， 在消费完之后，可以继续向里填充内容
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		produce(c)
		wg.Done()
	}()

	go func() {
		consume(c)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("all is over!")
}
