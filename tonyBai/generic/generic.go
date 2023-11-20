package main

import (
	"cmp"
	"fmt"
	"sync"
)

/*func Add[T constraint.AndExpr](a, b T) {
	return a + b
}*/

/*type maxableSlice[T any] struct {
	elems []T
}

func main() {

	var slice = maxableSlice[int]{elems: []int{1, 2, 3, 4}}
	fmt.Println(slice)
}*/

// 泛型类型

// 泛型类型嵌入
type Slice[T any] []T
type a[T cmp.Ordered] struct {
}

func (s Slice[T]) String() string { // 此处就是泛型类型方法
	if len(s) == 0 {
		return ""
	}
	result := fmt.Sprintf("%v", s[0])
	for _, v := range s[1:] {
		result += fmt.Sprintf("%v", v)
	}
	return result
}

type Lockable[T any] struct {
	t          T
	Slice[int] // 泛型类型实例化后可作为类型嵌入参数，嵌入泛型类型
	sync.Mutex
}

type Foo struct {
	Slice[int] // 泛型类型实例，嵌入普通结构体
}

func main() {
	l := Lockable[string]{
		t:     "hello",
		Slice: []int{1, 2, 3},
		Mutex: sync.Mutex{},
	}

	fmt.Println(l.String())

	f := Foo{
		Slice: []int{1, 3, 3, 4, 5},
	}
	fmt.Println("f:" + f.String())
}
