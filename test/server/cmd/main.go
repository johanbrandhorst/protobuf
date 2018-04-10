package main

import (
	"context"
	"io/ioutil"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
	"github.com/johanbrandhorst/protobuf/test/server/wrappers"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard))
	tc, err := credentials.NewClientTLSFromFile("../../insecure/localhost.crt", "")
	if err != nil {
		grpclog.Fatalln(err)
	}
	cc, err := grpc.Dial("localhost"+shared.HTTP2Server,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(tc),
		grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			grpclog.Infoln("Calling", method)
			defer func() { grpclog.Infoln("Finished", method) }()
			return streamer(ctx, desc, cc, method, opts...)
		}),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			grpclog.Infoln("Calling", method)
			defer func() { grpclog.Infoln("Finished", method) }()
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		grpclog.Fatalln(err)
	}

	client := wrappers.ClientWrapper{C: testproto.NewTestServiceClient(cc)}
	getStatus := func(err error) (codes.Code, string) {
		st, _ := status.FromError(err)
		return st.Code(), st.Message()
	}

	err = shared.TestPing(client, getStatus)
	if err != nil {
		grpclog.Fatalf("%+v\n", err)
	}

	err = shared.TestPingList(client, getStatus)
	if err != nil {
		grpclog.Fatalf("%+v\n", err)
	}

	err = shared.TestPingClientStream(client, getStatus)
	if err != nil {
		grpclog.Fatalf("%+v\n", err)
	}

	err = shared.TestPingBidiStream(client, getStatus)
	if err != nil {
		grpclog.Fatalf("%+v\n", err)
	}

	grpclog.Infoln("All tests successful")
}
