package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

// 测试正常分支
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
	myFrameCodec := NewMyFrameCodec()
	valueArray := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'} // 这里前4个字节代表长度，也就是totalLen， 然后以binary来读的时候，读到的数据是9，也是就是长度为9
	buf := bytes.NewBuffer(valueArray)
	frame, err := myFrameCodec.Decode(buf)
	if err != nil {
		t.Errorf("myFrameCodec decode err:%s", err.Error())
	}
	if string(frame) != "hello" {
		t.Errorf("decode wrong,want hello,real %s", string(frame))
	}
}

// 测试错误分支
type ErrorWriter struct {
	W  io.Writer
	Wn int // 第几次调用返回错误
	wc int // 写操作次数
}

func (w *ErrorWriter) Write(p []byte) (n int, err error) {
	w.wc++
	if w.wc >= w.Wn { // 写操作次数 大于等于 第n次调用返回错误
		return 0, errors.New("发生write错误")
	}
	return w.W.Write(p)
}

type ErrorReader struct {
	R  io.Reader
	Rn int // 第n次读返回错误
	rc int // 读操作次数
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	r.rc++
	if r.rc >= r.Rn { // 好巧妙
		return 0, errors.New("发生read错误")
	}
	return r.R.Read(p)
}

func TestEncodeErrorWrite(t *testing.T) {
	myCodec := NewMyFrameCodec()
	buf := make([]byte, 128)
	eW := &ErrorWriter{
		W:  bytes.NewBuffer(buf),
		Wn: 1,
	}
	// 验证binary write错误
	err := myCodec.Encode(eW, []byte("hello"))
	if err == nil {
		t.Error("binary.write没报错")
	}

	err = myCodec.Encode(&ErrorWriter{
		W:  bytes.NewBuffer(buf),
		Wn: 2,
	}, []byte("hello"))
	if err == nil {
		t.Error("w.write没出错")
	}
}

func TestDecodeErrorRead(t *testing.T) {
	myCodec := NewMyFrameCodec()
	array := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	_, err := myCodec.Decode(&ErrorReader{
		R:  bytes.NewReader(array),
		Rn: 1,
	})
	if err == nil {
		t.Errorf("binary.Read未报错")
	}

	_, err = myCodec.Decode(&ErrorReader{
		R:  bytes.NewReader(array),
		Rn: 2,
	})
	if err == nil {
		t.Errorf("r.R.Read未报错")
	}
}
