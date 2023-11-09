package unbuffered

import "testing"

var c1 chan string
var c2 chan string

// 无缓冲1收1发 测试性能
func init() {

	c1 = make(chan string)
	go func() {
		for {
			<-c1
		}
	}()

	c2 = make(chan string)
	go func() {
		for {
			c2 <- "hello"
		}
	}()
}

func send(s string) {
	c1 <- s
}

func receive() {
	<-c2
}

func BenchmarkUnbuffered1to1Send(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send("aaa")
	}
}

func BenchmarkUnbuffered1toReceive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		receive()
	}
}
