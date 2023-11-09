package buffered

import "testing"

var c1 chan string
var c2 chan string

func init() {
	c1 = make(chan string, 100)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				<-c1
			}
		}()
		go func() {
			for {
				c1 <- "mult"
			}
		}()
	}

	c2 = make(chan string, 100)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				c2 <- "mult"
			}
		}()

		go func() {
			for {
				<-c2
			}
		}()
	}
}

func send(s string) {
	c1 <- s
}

func receive() {
	<-c2
}

func BenchmarkNtoNSend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send("multi")
	}
}

func BenchmarkNtoNReceive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		receive()
	}
}
