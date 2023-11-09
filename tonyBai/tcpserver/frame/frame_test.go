package frame

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncode(t *testing.T) {
	myFrameCodec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)
	err := myFrameCodec.Encode(rw, []byte("hello"))
	if err != nil {
		t.Errorf("encode err:%s", err.Error())
	}

	var totalLen int32
	err = binary.Read(rw, binary.BigEndian, &totalLen)
	if err != nil {
		t.Errorf("binary read err:%s", err.Error())
	}

	if totalLen != 9 {
		t.Errorf("write len wrong, want 9 , real %d", totalLen)
	}

	writen := rw.Bytes()
	if string(writen) != "hello" {
		t.Errorf("write string wrong, real is %s", writen)
	}
}

func TestDecode(t *testing.T) {

}
