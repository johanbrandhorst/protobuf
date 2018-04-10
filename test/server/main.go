package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
	"github.com/johanbrandhorst/protobuf/test/server/proto/types"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

func main() {
	grpcServer := grpc.NewServer()
	testproto.RegisterTestServiceServer(grpcServer, &testSrv{})
	types.RegisterEchoServiceServer(grpcServer, &testSrv{})
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	}
	grpclog.SetLogger(log)
	wrappedServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			origin := req.Header.Get("Origin")
			parsedURL, err := url.ParseRequestURI(origin)
			if err != nil {
				grpclog.Warningf("Unable to parse url for grpc-websocket origin check: %s. error: %v", origin, err)
				return false
			}
			// Allow connections from any port
			return stripPort(parsedURL.Host) == stripPort(req.Host)
		}),
	)

	emptyGrpcServer := grpc.NewServer()
	emptyWrappedServer := grpcweb.WrapServer(emptyGrpcServer, grpcweb.WithCorsForRegisteredEndpointsOnly(false))
	emptyHandler := func(resp http.ResponseWriter, req *http.Request) {
		emptyWrappedServer.ServeHTTP(resp, req)
	}

	http1Server := http.Server{
		Addr:    shared.HTTP1Server,
		Handler: wrappedServer,
	}
	http1Server.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){} // Disable HTTP2
	http1EmptyServer := http.Server{
		Addr:    shared.EmptyHTTP1Server,
		Handler: http.HandlerFunc(emptyHandler),
	}
	http1EmptyServer.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){} // Disable HTTP2

	http2Server := http.Server{
		Addr:    shared.HTTP2Server,
		Handler: wrappedServer,
	}
	http2EmptyServer := http.Server{
		Addr:    shared.EmptyHTTP2Server,
		Handler: http.HandlerFunc(emptyHandler),
	}

	grpclog.Printf("Starting servers")

	// Start the empty Http1.1 server
	go func() {
		if err := http1EmptyServer.ListenAndServeTLS("./insecure/localhost.crt", "./insecure/localhost.key"); err != nil {
			grpclog.Fatalf("failed starting http1.1 empty server: %v", err)
		}
	}()

	// Start the Http1.1 server
	go func() {
		if err := http1Server.ListenAndServeTLS("./insecure/localhost.crt", "./insecure/localhost.key"); err != nil {
			grpclog.Fatalf("failed starting http1.1 server: %v", err)
		}
	}()

	// Start the empty Http2 server
	go func() {
		if err := http2EmptyServer.ListenAndServeTLS("./insecure/localhost.crt", "./insecure/localhost.key"); err != nil {
			grpclog.Fatalf("failed starting http2 empty server: %v", err)
		}
	}()

	// Start the Http2 server
	go func() {
		if err := http2Server.ListenAndServeTLS("./insecure/localhost.crt", "./insecure/localhost.key"); err != nil {
			grpclog.Fatalf("failed starting http2 server: %v", err)
		}
	}()

	// Host the GopherJS code
	httpsSrv := &http.Server{
		Addr:    shared.GopherJSServer,
		Handler: http.FileServer(http.Dir("./client/html")),
	}
	grpclog.Printf("Serving tests on https://localhost:8080")
	if err := httpsSrv.ListenAndServeTLS("./insecure/localhost.crt", "./insecure/localhost.key"); err != nil {
		grpclog.Fatalf("failed to start JS server: %v", err)
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

type testSrv struct {
}

func (s *testSrv) PingEmpty(ctx context.Context, _ *empty.Empty) (*testproto.PingResponse, error) {
	grpc.SendHeader(ctx, metadata.Pairs("HeaderTestKey1", "ServerValue1", "HeaderTestKey2", "ServerValue2"))
	grpc.SetTrailer(ctx, metadata.Pairs("TrailerTestKey1", "ServerValue1", "TrailerTestKey2", "ServerValue2"))
	return &testproto.PingResponse{Value: "foobar"}, nil
}

func (s *testSrv) Ping(ctx context.Context, ping *testproto.PingRequest) (*testproto.PingResponse, error) {
	if ping.GetCheckMetadata() {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md[strings.ToLower(shared.ClientMDTestKey)]) == 0 ||
			md[strings.ToLower(shared.ClientMDTestKey)][0] != shared.ClientMDTestValue {
			return nil, status.Errorf(codes.InvalidArgument, "Metadata was invalid")
		}
	}
	if ping.GetSendHeaders() {
		grpc.SendHeader(
			ctx,
			metadata.Pairs(
				shared.ServerMDTestKey1, shared.ServerMDTestValue1,
				shared.ServerMDTestKey2, shared.ServerMDTestValue2))
	}
	if ping.GetSendTrailers() {
		grpc.SetTrailer(
			ctx,
			metadata.Pairs(
				shared.ServerTrailerTestKey1, shared.ServerMDTestValue1,
				shared.ServerTrailerTestKey2, shared.ServerMDTestValue2))
	}
	return &testproto.PingResponse{Value: ping.GetValue(), Counter: ping.GetResponseCount()}, nil
}

func (s *testSrv) PingError(ctx context.Context, ping *testproto.PingRequest) (*empty.Empty, error) {
	if ping.FailureType == testproto.PingRequest_DROP {
		t, _ := transport.StreamFromContext(ctx)
		_ = t.ServerTransport().Close()
		return nil, status.Errorf(codes.Unavailable, "You got closed. You probably won't see this error")

	}
	if ping.GetSendHeaders() {
		grpc.SendHeader(
			ctx,
			metadata.Pairs(
				shared.ServerMDTestKey1, shared.ServerMDTestValue1,
				shared.ServerMDTestKey2, shared.ServerMDTestValue2))
	}
	if ping.GetSendTrailers() {
		grpc.SetTrailer(
			ctx,
			metadata.Pairs(
				shared.ServerTrailerTestKey1, shared.ServerMDTestValue1,
				shared.ServerTrailerTestKey2, shared.ServerMDTestValue2))
	}
	return nil, status.Errorf(codes.Code(ping.ErrorCodeReturned), ping.Value)
}

func (s *testSrv) PingList(ping *testproto.PingRequest, stream testproto.TestService_PingListServer) error {
	if ping.GetCheckMetadata() {
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok || len(md[strings.ToLower(shared.ClientMDTestKey)]) == 0 ||
			md[strings.ToLower(shared.ClientMDTestKey)][0] != shared.ClientMDTestValue {
			return status.Errorf(codes.InvalidArgument, "Metadata was invalid")
		}
	}
	if ping.GetSendHeaders() {
		stream.SendHeader(
			metadata.Pairs(
				shared.ServerMDTestKey1, shared.ServerMDTestValue1,
				shared.ServerMDTestKey2, shared.ServerMDTestValue2))
	}
	if ping.GetSendTrailers() {
		stream.SetTrailer(
			metadata.Pairs(
				shared.ServerTrailerTestKey1, shared.ServerMDTestValue1,
				shared.ServerTrailerTestKey2, shared.ServerMDTestValue2))
	}
	if ping.FailureType == testproto.PingRequest_DROP {
		t, _ := transport.StreamFromContext(stream.Context())
		_ = t.ServerTransport().Close()
		return status.Errorf(codes.Unavailable, "You got closed. You probably won't see this error")
	}
	if ping.GetFailureType() == testproto.PingRequest_CODE {
		return status.Errorf(codes.Code(ping.ErrorCodeReturned), ping.GetValue())
	}
	for i := int32(0); i < ping.ResponseCount; i++ {
		sleepDuration := ping.GetMessageLatencyMs()
		time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
		stream.Send(&testproto.PingResponse{Value: fmt.Sprintf("%s %d", ping.Value, i), Counter: i})
		if sleepDuration != 0 {
			// Flush the stream
			lowLevelServerStream, ok := transport.StreamFromContext(stream.Context())
			if !ok {
				return status.Errorf(codes.Internal, "lowLevelServerStream does not exist in context")
			}
			lowLevelServerStream.ServerTransport().Write(lowLevelServerStream, nil, nil, &transport.Options{
				Delay: false,
			})
		}
	}
	return nil
}

func (s *testSrv) EchoAllTypes(ctx context.Context, in *types.TestAllTypes) (*types.TestAllTypes, error) {
	return in, nil
}

func (s *testSrv) EchoMaps(ctx context.Context, in *types.TestMap) (*types.TestMap, error) {
	return in, nil
}

func (s *testSrv) PingClientStream(stream testproto.TestService_PingClientStreamServer) error {
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			time.Sleep(time.Duration(ping.GetMessageLatencyMs()) * time.Millisecond)
			return stream.SendAndClose(&testproto.PingResponse{Value: "Closed"})
		}
		if ping.GetCheckMetadata() {
			md, ok := metadata.FromIncomingContext(stream.Context())
			if !ok || len(md[strings.ToLower(shared.ClientMDTestKey)]) == 0 ||
				md[strings.ToLower(shared.ClientMDTestKey)][0] != shared.ClientMDTestValue {
				return status.Errorf(codes.InvalidArgument, "Metadata was invalid")
			}
		}
		if ping.GetSendHeaders() {
			stream.SendHeader(
				metadata.Pairs(
					shared.ServerMDTestKey1, shared.ServerMDTestValue1,
					shared.ServerMDTestKey2, shared.ServerMDTestValue2))
		}
		if ping.GetSendTrailers() {
			stream.SetTrailer(
				metadata.Pairs(
					shared.ServerTrailerTestKey1, shared.ServerMDTestValue1,
					shared.ServerTrailerTestKey2, shared.ServerMDTestValue2))
		}
		if err != nil {
			return err
		}
	}
}

func (s *testSrv) PingClientStreamError(stream testproto.TestService_PingClientStreamErrorServer) error {
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			time.Sleep(time.Duration(ping.GetMessageLatencyMs()) * time.Millisecond)
			_ = stream.SendAndClose(&testproto.PingResponse{Value: "Closed"})
			return status.Errorf(codes.Internal, "error")
		}
		if err != nil {
			return err
		}
		if ping.FailureType == testproto.PingRequest_DROP {
			t, _ := transport.StreamFromContext(stream.Context())
			_ = t.ServerTransport().Close()
			return status.Errorf(codes.Unavailable, "You got closed. You probably won't see this error")
		}
		if ping.GetFailureType() == testproto.PingRequest_CODE {
			return status.Errorf(codes.Code(ping.ErrorCodeReturned), ping.GetValue())
		}
	}
}
func (s *testSrv) PingBidiStream(stream testproto.TestService_PingBidiStreamServer) error {
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if ping.GetCheckMetadata() {
			md, ok := metadata.FromIncomingContext(stream.Context())
			if !ok || len(md[strings.ToLower(shared.ClientMDTestKey)]) == 0 ||
				md[strings.ToLower(shared.ClientMDTestKey)][0] != shared.ClientMDTestValue {
				return status.Errorf(codes.InvalidArgument, "Metadata was invalid")
			}
		}
		if ping.GetSendHeaders() {
			stream.SendHeader(
				metadata.Pairs(
					shared.ServerMDTestKey1, shared.ServerMDTestValue1,
					shared.ServerMDTestKey2, shared.ServerMDTestValue2))
		}
		if ping.GetSendTrailers() {
			stream.SetTrailer(
				metadata.Pairs(
					shared.ServerTrailerTestKey1, shared.ServerMDTestValue1,
					shared.ServerTrailerTestKey2, shared.ServerMDTestValue2))
		}
		time.Sleep(time.Duration(ping.GetMessageLatencyMs()) * time.Millisecond)
		err = stream.Send(&testproto.PingResponse{
			Value:   ping.GetValue(),
			Counter: ping.GetResponseCount(),
		})
		if err != nil {
			return err
		}
	}
}

func (s *testSrv) PingBidiStreamError(stream testproto.TestService_PingBidiStreamErrorServer) error {
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			return status.Errorf(codes.Internal, "error")
		}
		if err != nil {
			return err
		}
		if ping.GetFailureType() == testproto.PingRequest_CODE {
			return status.Errorf(codes.Code(ping.ErrorCodeReturned), ping.GetValue())
		}
		if ping.FailureType == testproto.PingRequest_DROP {
			t, _ := transport.StreamFromContext(stream.Context())
			_ = t.ServerTransport().Close()
			return status.Errorf(codes.Unavailable, "You got closed. You probably won't see this error")
		}
		time.Sleep(time.Duration(ping.GetMessageLatencyMs()) * time.Millisecond)
		err = stream.Send(&testproto.PingResponse{
			Value:   ping.GetValue(),
			Counter: ping.GetResponseCount(),
		})
		if err != nil {
			return err
		}
	}
}
