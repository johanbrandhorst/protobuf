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
		qunit.Expect(32)

		h := browserheaders.New(nil)
		assert.Equal(h.Len(), 0, "Len of an empty browserheader is 0")
		assert.Equal(len(h.Get("one")), 0, "Size of an unused key is 0")
		assert.Equal(h.HasKey("one"), false, `An empty browserheader does not have the key "one"`)
		assert.Equal(h.HasKeyWithValue("one", "1"), false, `An empty browserheader does not have a key "one" with a value of "1"`)

		h.Set("one", []string{"1"})
		assert.Equal(len(h.Get("one")), 1, `Set sets the size of the key "one" to 1`)
		assert.Equal(h.Get("one")[0], "1", `Set sets the value of "one" to "1"`)
		assert.Equal(h.HasKey("one"), true, `Set adds the key "one"`)
		assert.Equal(h.HasKeyWithValue("one", "1"), true, `Set sets the value of the key "one" to "1"`)

		h.Delete("one")
		assert.Equal(h.Len(), 0, "Delete reduced size of headers to 0")
		assert.Equal(len(h.Get("one")), 0, `Delete removes the key "one"`)
		assert.Equal(h.HasKey("one"), false, `Delete removes the key "one"`)
		assert.Equal(h.HasKeyWithValue("one", "1"), false, `Delete removes the key "one"`)

		h.Append("one", "1")
		assert.Equal(len(h.Get("one")), 1, `Append increases the size of key "one" by one`)
		assert.Equal(h.Get("one")[0], "1", `Append adds "1" to the key "one"`)
		assert.Equal(h.HasKey("one"), true, `Append adds the key "one"`)
		assert.Equal(h.HasKeyWithValue("one", "1"), true, `Append adds "1" the key "one"`)

		h.Append("one", "11")
		assert.Equal(len(h.Get("one")), 2, `Append increases the size of key "one" by one`)
		assert.Equal(h.Get("one")[1], "11", `Append adds "11" to the key "one"`)
		assert.Equal(h.HasKeyWithValue("one", "11"), true, `Append adds "11" the key "one"`)

		h.DeleteValueFromKey("one", "11")
		assert.Equal(len(h.Get("one")), 1, `DeleteValueFromKey removes one value from the key "one"`)
		assert.Equal(h.Get("one")[0], "1", `DeleteValueFromKey removes the correct value`)

		h.Append("two", "2")
		assert.Equal(len(h.Get("two")), 1, `Append increases the size of key "two" by one`)
		assert.Equal(h.Get("two")[0], "2", `Append adds "2" to the key "two"`)
		assert.Equal(h.HasKey("two"), true, `Append adds the key "two"`)
		assert.Equal(h.HasKeyWithValue("two", "2"), true, `Append adds "2" the key "two"`)

		items := []struct {
			Key    string
			Values []string
		}{}
		h.ForEach(func(key string, values []string) {
			items = append(items, struct {
				Key    string
				Values []string
			}{
				Key:    key,
				Values: values,
			})
		})
		assert.Equal(len(items), 2, "ForEach iterated over 2 keys")

		assert.Equal(items[0].Key, "one", `First key is "one"`)
		assert.Equal(len(items[0].Values), 1, "First key has one value")
		assert.Equal(items[0].Values[0], "1", `First keys value is "1"`)

		assert.Equal(items[1].Key, "two", `Second key is "two"`)
		assert.Equal(len(items[1].Values), 1, "Second key has one value")
		assert.Equal(items[1].Values[0], "2", `Second keys value is "2"`)
	})
}
