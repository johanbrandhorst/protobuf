package wrappers

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/johanbrandhorst/protobuf/grpcweb"
	testproto "github.com/johanbrandhorst/protobuf/test/client/proto/test"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

type ClientWrapper struct {
	C testproto.TestServiceClient
}

func (cw ClientWrapper) Ping(ctx context.Context, req *shared.Request, headers, trailers *metadata.MD) (*shared.Response, error) {
	resp, err := cw.C.Ping(ctx, sharedToProtoReq(req), grpcweb.Header(headers), grpcweb.Trailer(trailers))
	if err != nil {
		return nil, err
	}

	return (*shared.Response)(resp), nil
}

func (cw ClientWrapper) PingError(ctx context.Context, req *shared.Request) error {
	_, err := cw.C.PingError(ctx, sharedToProtoReq(req))

	return err
}

func (cw ClientWrapper) PingList(ctx context.Context, req *shared.Request, headers, trailers *metadata.MD) (shared.TestPingListClient, error) {
	pcs, err := cw.C.PingList(ctx, sharedToProtoReq(req), grpcweb.Header(headers), grpcweb.Trailer(trailers))
	if err != nil {
		return nil, err
	}
	return pingListWrapper{c: pcs}, nil
}

func (cw ClientWrapper) PingClientStream(ctx context.Context, headers, trailers *metadata.MD) (shared.TestPingClientStreamClient, error) {
	pcs, err := cw.C.PingClientStream(ctx, grpcweb.Header(headers), grpcweb.Trailer(trailers))
	if err != nil {
		return nil, err
	}
	return pingClientStreamWrapper{c: pcs}, nil
}

func (cw ClientWrapper) PingClientStreamError(ctx context.Context) (shared.TestPingClientStreamErrorClient, error) {
	pcse, err := cw.C.PingClientStreamError(ctx)
	if err != nil {
		return nil, err
	}
	return pingClientStreamErrorWrapper{c: pcse}, nil
}

func (cw ClientWrapper) PingBidiStream(ctx context.Context, headers, trailers *metadata.MD) (shared.TestPingBidiStreamClient, error) {
	pbs, err := cw.C.PingBidiStream(ctx, grpcweb.Header(headers), grpcweb.Trailer(trailers))
	if err != nil {
		return nil, err
	}
	return pingBidiStreamWrapper{c: pbs}, nil
}

func (cw ClientWrapper) PingBidiStreamError(ctx context.Context) (shared.TestPingBidiStreamErrorClient, error) {
	pbse, err := cw.C.PingBidiStreamError(ctx)
	if err != nil {
		return nil, err
	}
	return pingBidiStreamErrorWrapper{c: pbse}, nil
}

func sharedToProtoReq(req *shared.Request) *testproto.PingRequest {
	return &testproto.PingRequest{
		Value:             req.Value,
		ResponseCount:     req.ResponseCount,
		ErrorCodeReturned: req.ErrorCodeReturned,
		FailureType:       testproto.PingRequest_FailureType(req.FailureType),
		CheckMetadata:     req.CheckMetadata,
		SendHeaders:       req.SendHeaders,
		SendTrailers:      req.SendTrailers,
		MessageLatencyMs:  req.MessageLatencyMs,
	}
}
