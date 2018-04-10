package wrappers

import (
	"github.com/johanbrandhorst/protobuf/test/shared"
	"google.golang.org/grpc/metadata"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
)

type pingClientStreamWrapper struct {
	c testproto.TestService_PingClientStreamClient
}

func (pcsw pingClientStreamWrapper) Send(req *shared.Request) error {
	return pcsw.c.Send(sharedToProtoReq(req))
}

func (pcsw pingClientStreamWrapper) CloseAndRecv() (*shared.Response, error) {
	resp, err := pcsw.c.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return (*shared.Response)(resp), nil
}

func (pcsw pingClientStreamWrapper) Header() (metadata.MD, error) {
	return pcsw.c.Header()
}

func (pcsw pingClientStreamWrapper) Trailer() metadata.MD {
	return pcsw.c.Trailer()
}

type pingClientStreamErrorWrapper struct {
	c testproto.TestService_PingClientStreamErrorClient
}

func (pcsew pingClientStreamErrorWrapper) Send(req *shared.Request) error {
	return pcsew.c.Send(sharedToProtoReq(req))
}

func (pcsew pingClientStreamErrorWrapper) CloseAndRecv() (*shared.Response, error) {
	_, err := pcsew.c.CloseAndRecv()
	return nil, err
}
