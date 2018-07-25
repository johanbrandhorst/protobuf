package wrappers

import (
	"google.golang.org/grpc/metadata"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

type pingBidiStreamWrapper struct {
	c testproto.TestService_PingBidiStreamClient
}

func (pcsw pingBidiStreamWrapper) Send(req *shared.Request) error {
	return pcsw.c.Send(sharedToProtoReq(req))
}

func (pcsw pingBidiStreamWrapper) Recv() (*shared.Response, error) {
	resp, err := pcsw.c.Recv()
	if err != nil {
		return nil, err
	}

	return &shared.Response{
		Value:   resp.GetValue(),
		Counter: resp.GetCounter(),
	}, nil
}

func (pcsw pingBidiStreamWrapper) Header() (metadata.MD, error) {
	return pcsw.c.Header()
}

func (pcsw pingBidiStreamWrapper) Trailer() metadata.MD {
	return pcsw.c.Trailer()
}

func (pcsw pingBidiStreamWrapper) CloseSend() error {
	return pcsw.c.CloseSend()
}

type pingBidiStreamErrorWrapper struct {
	c testproto.TestService_PingBidiStreamErrorClient
}

func (pbsew pingBidiStreamErrorWrapper) Send(req *shared.Request) error {
	return pbsew.c.Send(sharedToProtoReq(req))
}

func (pbsew pingBidiStreamErrorWrapper) Recv() (*shared.Response, error) {
	resp, err := pbsew.c.Recv()
	if err != nil {
		return nil, err
	}

	return &shared.Response{
		Value:   resp.GetValue(),
		Counter: resp.GetCounter(),
	}, nil
}

func (pbsew pingBidiStreamErrorWrapper) CloseSend() error {
	return pbsew.c.CloseSend()
}
