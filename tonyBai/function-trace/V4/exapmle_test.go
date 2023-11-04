// instrument_trace/example_test.go
package trace_test

import (
	"funcTraceV4/trace"
)

func a() {
	defer trace.Trace()()
	b()
}

func b() {
	defer trace.Trace()()
	c()
}

func c() {
	defer trace.Trace()()
	d()
}

func d() {
	defer trace.Trace()()
}

func ExampleTrace() {
	a()
	//Output:
	//[00001]    ->funcTraceV4_test.a
	//[00001]        ->funcTraceV4_test.b
	//[00001]            ->funcTraceV4_test.c
	//[00001]                ->funcTraceV4_test.d
	//[00001]                <-funcTraceV4_test.d
	//[00001]            <-funcTraceV4_test.c
	//[00001]        <-funcTraceV4_test.b
	//[00001]    <-funcTraceV4_test.a
}
