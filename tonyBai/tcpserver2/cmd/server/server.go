package main

import (
	"fmt"
	"net"
	"tcpserver2/frame"
	"tcpserver2/metrics"
	"tcpserver2/packet"
)

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen err : " + err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err : " + err.Error())
			break
		}
		go handConn(conn)
	}
}

func handConn(conn net.Conn) {
	metrics.ClientConnected.Inc() // 建连时增加1
	defer func() {
		metrics.ClientConnected.Dec() // 断连时减1
		conn.Close()
	}()
	for {
		myFrameCodec := frame.NewMyFrameCodec()
		framePayload, err := myFrameCodec.Decode(conn)
		if err != nil {
			fmt.Println("myFrameCodec decode err：" + err.Error())
			return
		}
		metrics.ReqRecvTotal.Inc() // 收到请求，接收统计加1

		ackFramePayload, err := handPacket(framePayload)
		if err != nil {
			fmt.Println("handPacket err: " + err.Error())
			return
		}

		err = myFrameCodec.Encode(conn, ackFramePayload)
		if err != nil {
			fmt.Println("myFrameCodec encode err:" + err.Error())
			return
		}
		metrics.RspSendTotal.Inc() // 响应完后，发送统计加1
	}

}

func handPacket(framePayload []byte) ([]byte, error) {
	var p packet.Packet
	p, err := packet.Decode(framePayload)
	if err != nil {
		return nil, err
	}
	switch t := p.(type) {
	case *packet.Submit:
		var s *packet.Submit = p.(*packet.Submit)
		fmt.Printf("receive submitId[%s], submitPayload is [%s]\n", s.ID, string(s.Payload))
		submitAck := &packet.SubmitAck{
			ID:     s.ID,
			Result: 0,
		}
		ackFramePayload, err := packet.Encode(submitAck)
		if err != nil {
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("type [%s] not found", t)

	}
}
