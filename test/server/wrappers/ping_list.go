package wrappers

import (
	"github.com/johanbrandhorst/protobuf/test/shared"
	"google.golang.org/grpc/metadata"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
)

type pingListWrapper struct {
	c testproto.TestService_PingListClient
}

func (pcsw pingListWrapper) Recv() (*shared.Response, error) {
	resp, err := pcsw.c.Recv()
	if err != nil {
		return nil, err
	}

	return (*shared.Response)(resp), nil
}

func (pcsw pingListWrapper) Header() (metadata.MD, error) {
	return pcsw.c.Header()
}

func (pcsw pingListWrapper) Trailer() metadata.MD {
	return pcsw.c.Trailer()
}
