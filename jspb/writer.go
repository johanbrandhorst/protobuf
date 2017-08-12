package jspb

import "github.com/gopherjs/gopherjs/js"

// Writer encapsulates the jspb.BinaryWriter.
type Writer interface {
	GetResult() []byte

	// Scalars
	WriteInt32(int, int32)
	WriteUint32(int, uint32)
	WriteString(int, string)
	WriteBool(int, bool)

	// Slices
	WriteUint32Slice(int, []uint32)

	// Specials
	WriteMessage(int, func())
	WriteEnum(int, int)
}

// NewWriter returns a new Writer ready for writing.
func NewWriter() Writer {
	return &writer{
		Object: js.Global.Get("BinaryWriter").New(),
	}
}

// writer implements the Writer interface.
type writer struct {
	*js.Object
}

// GetResult returns the contents of the buffer as a byte slice.
func (w writer) GetResult() []byte {
	return w.Call("getResultBuffer").Interface().([]byte)
}

// WriteInt32 writes an int32 field to the buffer.
func (w writer) WriteInt32(field int, value int32) {
	w.Call("writeInt32", field, value)
}

// WriteInt32 writes a uint32 field to the buffer.
func (w writer) WriteUint32(field int, value uint32) {
	w.Call("writeUint32", field, value)
}

// WriteString writes a string field to the buffer
func (w writer) WriteString(field int, value string) {
	w.Call("writeString", field, value)
}

// WriteBool writes a string field to the buffer
func (w writer) WriteBool(field int, value bool) {
	w.Call("writeBool", field, value)
}

// WriteUint32Slice writes a uint32 slice field to the buffer
func (w writer) WriteUint32Slice(field int, value []uint32) {
	w.Call("writePackedUint32", field, value)
}

// WriteMessage writes a message to the buffer using writeFunc
func (w writer) WriteMessage(field int, writeFunc func()) {
	w.Call("writeMessage", field, 0 /* Unused */, writeFunc)
}

// WriteEnum writes an enum (as an int) to the buffer
func (w writer) WriteEnum(field int, value int) {
	w.Call("writeEnum", field, value)
}
