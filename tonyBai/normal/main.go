package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

/*
// 测试Go中的内存对齐
type A struct {
	a uint16
	b int64
	c byte
}

type B struct {
	a uint16
	c byte
	b int64
}

func main() {

		 内存对齐：
			1.对齐大小至少为1
			2.每个字段需要对齐
			3.结构体整体需要对齐

	fmt.Println(unsafe.Sizeof(A{})) // 24
	fmt.Println(unsafe.Sizeof(B{})) // 16
}*/

/*
循环变量重用，使用闭包解决
*/
/*func main() {
	var m = []int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func(i, v int) {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}(i, v)
	}

	time.Sleep(time.Second * 10)
}*/

/*
数组是值拷贝，而切片是引用，操作后变化的是底层数组
*/
/*func main() {
	var a = []int{1, 2, 3, 4, 5}
	var r []int

	fmt.Println("original a =", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r = append(r, v)
		//r[i] = v
	}

	fmt.Println("after for range loop, r =", r)
	fmt.Println("after for range loop, a =", a)
}*/

/*// map遍历是无序的
func main() {
	var m = map[string]int{
		"tony": 21,
		"tom":  22,
		"jim":  23,
	}

	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "tony")
		}
		counter++
		fmt.Println(k, v)
	}
	fmt.Println("counter is ", counter)
}*/

/*func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello Gopher!")
}
func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(greeting))
	err1 := errors.New("aaa")

	err2 := fmt.Errorf("err2,warpErr1:%w", err1) // Errorf可以对错误进行包装，然后可以通过errors.Is来判断错误
	if err2 != nil {
		panic(err2)
	}


}*/

/*func main() {
	var array = []byte{0x0, 0x0, 0x0, 0x21, 'h', 'e', 'l'}
	var array1 = []byte{'h', 'e', 'l'}
	fmt.Println("array:" + string(array))
	fmt.Println("array:" + string(array1))
}*/

func main() {
	// 模拟一个实现了 io.Reader 的对象
	reader := bytes.NewBuffer([]byte{0x9, 0x0, 0x0, 0x0, 'h', 'e', 'l', 'l', 'o'}) // ASCII 编码的 "Hello"

	var totalLen int32
	err := binary.Read(reader, binary.LittleEndian, &totalLen)
	fmt.Printf("totalLen:%d\n", int(totalLen))
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// 读取 5 个字节的数据
	buf := make([]byte, 5)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// 打印读取到的字节数组
	fmt.Printf("Read data: %v\n", buf)
}
