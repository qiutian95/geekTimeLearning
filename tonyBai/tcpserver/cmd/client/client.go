package main

import (
	"fmt"
	"github.com/lucasepe/codename"
	"net"
	"sync"
	"tcpserver/frame"
	"tcpserver/packet"
	"time"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i)
	}

	wg.Wait()
}

func startClient(i int) {
	quit := make(chan struct{})
	done := make(chan struct{})

	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("net dial err: " + err.Error())
	}
	defer conn.Close()
	fmt.Printf("dial tcp success, client [%d]\n", i)

	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	frameCodec := frame.NewMyFrameCodec()
	var count int // 控制send的数量

	go func() {
		// handle ack
		for {
			select {
			case <-quit:
				done <- struct{}{}
				return
			default:
			}

			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			framePayload, err := frameCodec.Decode(conn)
			if err != nil {
				fmt.Println("client frameCode decode err: " + err.Error())
				return
			}
			p, err := packet.Decode(framePayload)
			if err != nil {
				fmt.Println("client packet decode err: " + err.Error())
				return
			}
			subAck, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("not submitAck")
			}
			fmt.Printf("client receive submitAck id: [%s], result:[%d]\n", subAck.ID, subAck.Result)
		}
	}()

	for {
		// send submit
		count++
		payload := codename.Generate(rng, 4)
		id := fmt.Sprintf("%08d", count) // 8个字节宽度
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}
		framePayload, err := packet.Encode(s)
		if err != nil {
			fmt.Println("client packet encode err:" + err.Error())
			return
		}
		fmt.Printf("client[%d] send submit id:%s, payload:%s, len:%d\n", count, s.ID, s.Payload, len(s.Payload)+4)
		err = frameCodec.Encode(conn, framePayload)
		if err != nil {
			fmt.Println("client frameCodec encode err: " + err.Error())
			return
		}

		time.Sleep(time.Second)

		if count >= 10 {
			quit <- struct{}{}
			<-done
			fmt.Printf("client[%d] exit succ\n", count)
			return
		}
	}
}
