package shared

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Addresses for the various servers
const (
	HTTP1Server           = ":9090"
	EmptyHTTP1Server      = ":9095"
	HTTP2Server           = ":9100"
	EmptyHTTP2Server      = ":9105"
	GopherJSServer        = ":8080"
	ClientMDTestKey       = "HeaderTestKey1"
	ClientMDTestValue     = "ClientValue1"
	ServerMDTestKey1      = "HeaderTestKey1"
	ServerTrailerTestKey1 = "TrailerTestKey1"
	ServerMDTestValue1    = "ServerValue1"
	ServerTrailerTestKey2 = "TrailerTestKey2"
	ServerMDTestKey2      = "HeaderTestKey2"
	ServerMDTestValue2    = "ServerValue2"
)

type TestClient interface {
	Ping(context.Context, *Request, *metadata.MD, *metadata.MD) (*Response, error)
	PingError(context.Context, *Request) error
	PingList(context.Context, *Request, *metadata.MD, *metadata.MD) (TestPingListClient, error)
	PingClientStream(context.Context, *metadata.MD, *metadata.MD) (TestPingClientStreamClient, error)
	PingClientStreamError(context.Context) (TestPingClientStreamErrorClient, error)
	PingBidiStream(context.Context, *metadata.MD, *metadata.MD) (TestPingBidiStreamClient, error)
	PingBidiStreamError(context.Context) (TestPingBidiStreamErrorClient, error)
}

type FailureType int

const (
	NONE FailureType = iota
	CODE
	DROP
)

type Request struct {
	Value             string
	ResponseCount     int32
	ErrorCodeReturned uint32
	FailureType       FailureType
	CheckMetadata     bool
	SendHeaders       bool
	SendTrailers      bool
	MessageLatencyMs  int32
}

type Response struct {
	Value   string
	Counter int32
}
