package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var gourtineSpace = []byte("goroutine ")

func goroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, gourtineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "    "
	}
	fmt.Printf("[%05d]%s%s%s\n", id, indents, arrow, name)

}

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	fn := runtime.FuncForPC(pc)
	st := fn.Name()
	id := goroutineId()

	mu.Lock()
	indents := m[id]
	m[id] = indents + 1
	mu.Unlock()

	printTrace(id, st, "->", indents+1)
	return func() {
		mu.Lock()
		indents = m[id]
		m[id] = indents - 1
		mu.Unlock()
		printTrace(id, st, "<-", indents)
	}
}
