package grpcweb

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// MD represents a version of a typical metadata
// This types primary purpose is to automatically lowercase
// keys, since header keys are lowercase when returned from gRPC-Web.
type MD struct {
	md metadata.MD
}

// Get gets all values assigned to this header key.
// Mutation of the returned slice mutates the underlying metadata.
// This functions primary purpose is to automatically lowercase
// keys, since header keys are lowercase when returned from gRPC-Web.
func (m MD) Get(key string) []string {
	return m.md[strings.ToLower(key)]
}

// Metadata returns the underlying metadata
func (m MD) Metadata() metadata.MD {
	return m.md
}

// Len returns the size of the metadata stored
func (m MD) Len() int {
	return len(m.md)
}
