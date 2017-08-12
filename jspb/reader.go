package jspb

import "github.com/gopherjs/gopherjs/js"

// Reader encapsulates the jspb.BinaryReader
type Reader interface {
	Next() bool
	Err() error

	GetFieldNumber() int
	SkipField()

	// Scalars
	ReadInt32() int32
	ReadUint32() uint32
	ReadString() string
	ReadBool() bool

	// Slices
	ReadUint32Slice() []uint32

	// Specials
	ReadMessage(func())
	ReadEnum() int
}

// NewReader returns a new Reader ready for writing
func NewReader(data []byte) Reader {
	return &reader{
		Object: js.Global.Get("BinaryReader").New(data),
	}
}

// reader implements the Reader interface
type reader struct {
	*js.Object
	err error
}

// Reads the next field header in the stream if there is one, returns true if
// we saw a valid field header or false if we've read the whole stream.
// Sets err if we encountered a deprecated START_GROUP/END_GROUP field.
func (r *reader) Next() bool {
	defer catchException(&r.err)
	return r.err == nil && r.Call("nextField").Bool() && !r.Call("isEndGroup").Bool()
}

// Err returns the error state of the Reader.
func (r reader) Err() error {
	return r.err
}

// The field number of the next field in the buffer, or
// InvalidFieldNumber if there is no next field.
func (r reader) GetFieldNumber() int {
	return r.Call("getFieldNumber").Int()
}

// Skips over the next field in the binary stream.
func (r reader) SkipField() {
	r.Call("skipField")
}

// ReadInt32 reads a signed 32-bit integer field from the binary
// stream, sets err if the next field in the
// stream is not of the correct wire type.
func (r *reader) ReadInt32() int32 {
	defer catchException(&r.err)
	return int32(r.Call("readInt32").Int())
}

// ReadUit32 reads an unsigned 32-bit integer field from the binary
// stream, sets err if the next field in the
// stream is not of the correct wire type.
func (r *reader) ReadUint32() uint32 {
	defer catchException(&r.err)
	return uint32(r.Call("readUint32").Int())
}

// ReadString reads a string field from the binary stream, sets err
// if the next field in the stream is not of the correct wire type.
func (r *reader) ReadString() string {
	defer catchException(&r.err)
	return r.Call("readString").String()
}

// ReadBool reads a bool field from the binary stream, sets err
// if the next field in the stream is not of the correct wire type.
func (r *reader) ReadBool() bool {
	defer catchException(&r.err)
	return r.Call("readBool").Bool()
}

// ReadUint32Slice reads a packed 32-bit unsigned integer field
// from the binary stream, sets err if the next field
// in the stream is not of the correct wire type.
func (r *reader) ReadUint32Slice() (ret []uint32) {
	defer catchException(&r.err)
	values := r.Call("readPackedUint32").Interface().([]interface{})
	for _, value := range values {
		ret = append(ret, uint32(value.(float64)))
	}

	return ret
}

// ReadMessage deserializes a proto using
// the provided reader function.
func (r *reader) ReadMessage(readFunc func()) {
	defer catchException(&r.err)
	r.Call("readMessage", js.Undefined /* Unused */, readFunc)
}

// ReadEnum reads an enum field from the binary stream,
// sets err if the next field in the stream
// is not of the correct wire type.
func (r *reader) ReadEnum() int {
	defer catchException(&r.err)
	return r.Call("readEnum").Int()
}
