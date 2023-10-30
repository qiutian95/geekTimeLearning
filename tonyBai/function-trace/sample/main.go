package main

import "fmt"

func Trace(st string) func() {
	fmt.Println("enter:", st)
	return func() {
		fmt.Println("exit:", st)
	}
}

func foo() {
	defer Trace("foo")()
	Bar()
}

func Bar() {
	defer Trace("bar")()
}

func main() {
	defer Trace("main")()
	foo()
}
