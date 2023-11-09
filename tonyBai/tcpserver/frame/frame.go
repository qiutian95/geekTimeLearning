package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(io.Writer, FramePayload) error          // 编码
	Decode(reader io.Reader) (FramePayload, error) // 解码
}

var ErrShortWrite = errors.New("write short")
var ErrShortRead = errors.New("read short")

type MyFrameCodec struct{}

func NewMyFrameCodec() *MyFrameCodec {
	return &MyFrameCodec{}
}

func (m *MyFrameCodec) Encode(w io.Writer, payload FramePayload) error {
	var totalLength int32 = int32(len(payload)) + 4
	err := binary.Write(w, binary.BigEndian, &totalLength) // 先将长度写入w中，读取的时候用int32类型来读，就可以读出长度
	if err != nil {
		return errors.New("write binary err: " + err.Error())
	}
	n, err := w.Write([]byte(payload))
	if err != nil {
		return errors.New("write from payload err:" + err.Error())
	}

	if n != len(payload) {
		return ErrShortWrite
	}

	return nil

}

func (m *MyFrameCodec) Decode(r io.Reader) (FramePayload, error) {
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil

}
