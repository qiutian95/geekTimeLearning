package packet

import (
	"bytes"
	"fmt"
)

const (
	CommandCon = iota + 0x1
	CommandSubmit
)

const (
	CommandConAck = iota + 0x81
	CommandSubmitAck
)

type Packet interface {
	Decode([]byte) error
	Encode() ([]byte, error)
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(pktBody []byte) error {
	s.ID = string(pktBody[:8])
	s.Payload = pktBody[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (sa *SubmitAck) Decode(pktBody []byte) error {
	sa.ID = string(pktBody[:8])
	sa.Result = uint8(pktBody[8])
	return nil
}

func (sa *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(sa.ID[:8]), []byte{sa.Result}}, nil), nil // TODO 这里sa.ID[:8]为的是防止超长吗？
}

func Decode(frameBody []byte) (Packet, error) {
	id := frameBody[0]
	packet := frameBody[1:]
	switch id {
	case CommandCon:
		return nil, nil
	case CommandConAck:
		return nil, nil
	case CommandSubmit:
		var submit Submit
		submit.ID = string(id)
		err := submit.Decode(packet)
		if err != nil {
			return nil, err
		}
		return &submit, nil
	case CommandSubmitAck:
		var submitack SubmitAck
		submitack.ID = string(id)
		err := submitack.Decode(packet)
		if err != nil {
			return nil, err
		}
		return &submitack, nil
	default:
		return nil, fmt.Errorf("commandID : [%d] not fount", id)

	}
}

func Encode(p Packet) ([]byte, error) {
	var id uint8
	var body []byte
	var err error
	switch t := p.(type) {
	case *Submit:
		id = CommandSubmit
		body, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		id = CommandSubmitAck
		body, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("未找到类型：[%s]", t)
	}
	return bytes.Join([][]byte{[]byte{id}, body}, nil), nil
}
