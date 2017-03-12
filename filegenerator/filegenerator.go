package filegenerator

import (
	"fmt"
	"io"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type FileGenerator struct {
	w      io.Writer
	indent uint
}

func New(w io.Writer) *FileGenerator {
	return &FileGenerator{
		w: w,
	}
}

func (fg *FileGenerator) In() {
	fg.indent++
}

func (fg *FileGenerator) Out() {
	fg.indent--
}

func (fg *FileGenerator) P(format string, a ...interface{}) error {
	var err error

	// If format is empty, avoid printing just whitespaces.
	if format != "" {
		_, err = fmt.Fprintf(fg.w, strings.Repeat("    ", int(fg.indent)))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(fg.w, format, a...)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(fg.w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func (fg *FileGenerator) Generate(file *descriptor.FileDescriptorProto) {
	fg.P(`package %s`, file.GetPackage())
	fg.P("")

	fg.P(`import "github.com/gopherjs/gopherjs/js"`)
	fg.P("")

	for _, msg := range file.GetMessageType() {
		fg.generateProtoMessage(file, msg)
	}
}

func (fg *FileGenerator) generateProtoMessage(file *descriptor.FileDescriptorProto, message *descriptor.DescriptorProto) {
	ccTypeName := CamelCase(message.GetName())

	fg.P(`type %s struct {`, ccTypeName)
	fg.In()
	fg.P(`*js.Object`)
	for _, field := range message.GetField() {
		fg.P(`%s %s `+"`js:"+`"%s"`+"`", CamelCase(field.GetName()), GoType(message, field), field.GetJsonName())
	}
	fg.Out()
	fg.P(`}`)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// CamelCaseSlice is like CamelCase, but the argument is a slice of strings to
// be joined with "_".
func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

// GoType returns a string representing the type name
func GoType(message *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) (typ string) {
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typ = "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typ = "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typ = "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		typ = "string"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		typ = "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM, descriptor.FieldDescriptorProto_TYPE_GROUP, descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		typ = field.GetTypeName()
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		typ = "int64"
	default:
		panic("unknown type for " + field.GetName())
	}
	if needsStar(field, field.Extendee == nil, message != nil) {
		typ = "*" + typ
	}
	if isRepeated(field) {
		typ = "[]" + typ
	}
	return
}

func needsStar(field *descriptor.FieldDescriptorProto, proto3 bool, allowOneOf bool) bool {
	if isRepeated(field) &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE) &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_GROUP) {
		return false
	}
	if *field.Type == descriptor.FieldDescriptorProto_TYPE_BYTES {
		return false
	}
	if field.OneofIndex != nil && allowOneOf &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE) &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_GROUP) {
		return false
	}
	if proto3 &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE) &&
		(*field.Type != descriptor.FieldDescriptorProto_TYPE_GROUP) {
		return false
	}
	return true
}

// Is this field repeated?
func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}
