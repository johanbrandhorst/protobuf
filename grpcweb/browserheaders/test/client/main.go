package main

import (
	"fmt"

	"github.com/rusco/qunit"

	"github.com/johanbrandhorst/protobuf/grpcweb/browserheaders"
	"google.golang.org/grpc/metadata"
)

//go:generate gopherjs build main.go -m -o html/index.js

func recoverer() {
	e := recover()
	if e == nil {
		return
	}

	qunit.Ok(false, fmt.Sprintf("Saw panic: %v", e))
}

func main() {
	defer recoverer() // recovers any panics and fails tests

	qunit.Module("grpcweb")

	qunit.Test("Creating a new browserheaders type", func(assert qunit.QUnitAssert) {
		qunit.Expect(1)

		h := browserheaders.New(nil)
		assert.Equal(h.Len(), 0, "Len of an empty browserheader is 0")
	})

	qunit.Test("Creating a new browserheaders type with metadata", func(assert qunit.QUnitAssert) {
		qunit.Expect(6)

		h := browserheaders.New(metadata.Pairs("one", "1", "two", "2", "one", "11"))
		assert.Equal(h.Len(), 2, "Len is 2")
		assert.Equal(len(h.Get("one")), 2, `Size of value of key "one" is 2`)
		assert.Equal(h.Get("one")[0], "1", `First value of "one" is "1"`)
		assert.Equal(h.Get("one")[1], "11", `Second value of "one" is "11"`)
		assert.Equal(len(h.Get("two")), 1, `Size of value of key "two" is 1`)
		assert.Equal(h.Get("two")[0], "2", `Value of "two" is "2"`)
	})

	qunit.Test("Set, Get and Delete", func(assert qunit.QUnitAssert) {
		qunit.Expect(5)

		h := browserheaders.New(nil)
		assert.Equal(h.Len(), 0, "Len of an empty browserheader is 0")
		assert.Equal(len(h.Get("one")), 0, "Size of an unused key is 0")

		h.Set("one", []string{"1"})
		assert.Equal(len(h.Get("one")), 1, `Size of value of key "one" is 1`)
		assert.Equal(h.Get("one")[0], "1", `Value of "one" is "1"`)

		h.Delete("one")
		assert.Equal(len(h.Get("one")), 0, "Size of an unused key is 0")
	})
}
