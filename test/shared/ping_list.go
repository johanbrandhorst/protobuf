package shared

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type TestPingListClient interface {
	Recv() (*Response, error)
	Trailer() metadata.MD
	Header() (metadata.MD, error)
}

func testPingList(client TestClient, req *Request) error {
	headers, trailers := metadata.MD{}, metadata.MD{}
	ctx := context.Background()
	if req.CheckMetadata {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(ClientMDTestKey, ClientMDTestValue))
	}
	srv, err := client.PingList(ctx, req, &headers, &trailers)
	if err != nil {
		return unexpectedError("PingList", err)
	}

	var i int32
	for ; ; i++ {
		resp, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return unexpectedError("Recv", err)
		}

		want := fmt.Sprintf("%s %d", req.Value, i)
		if resp.Value != want {
			return reportError("response value", resp.Value, want)
		}
		if resp.Counter != i {
			return reportError("response counter", resp.Counter, i)
		}
	}

	if i != req.ResponseCount {
		return reportError("number of replies", i, req.ResponseCount)
	}

	// Headers used as callOptions should not be populated
	if len(headers) > 0 {
		return reportError("header", headers, nil)
	}

	// Trailers used as callOptions should not be populated
	if len(trailers) > 0 {
		return reportError("trailer", trailers, nil)
	}

	h, err := srv.Header()
	if err != nil {
		return unexpectedError("header", err)
	}

	if req.SendHeaders {
		for header, value := range map[string]string{
			ServerMDTestKey1: ServerMDTestValue1,
			ServerMDTestKey2: ServerMDTestValue2,
		} {
			if values, ok := h[strings.ToLower(header)]; !ok {
				return reportError("header", h, header)
			} else if len(values) != 1 || values[0] != value {
				return reportError("header value", values, value)
			}
		}
	} else {
		for _, header := range []string{ServerMDTestKey1, ServerMDTestKey2} {
			_, ok := h[header]
			if ok {
				return reportError("unexpected header", h[ServerMDTestKey1], "")
			}
		}
	}

	t := srv.Trailer()
	if req.SendTrailers {
		for trailer, value := range map[string]string{
			ServerTrailerTestKey1: ServerMDTestValue1,
			ServerTrailerTestKey2: ServerMDTestValue2,
		} {
			if values, ok := t[strings.ToLower(trailer)]; !ok {
				return reportError("trailer", t, trailer)
			} else if len(values) != 1 || values[0] != value {
				return reportError("trailer value", values, value)
			}
		}
	} else {
		for _, trailer := range []string{ServerMDTestKey1, ServerMDTestKey2} {
			_, ok := h[trailer]
			if ok {
				return reportError("unexpected trailer", t[ServerMDTestKey1], "")
			}
		}
	}

	return nil
}

func testPingListError(client TestClient, req *Request, getStatus func(error) (codes.Code, string)) error {
	headers, trailers := metadata.MD{}, metadata.MD{}
	srv, err := client.PingList(context.Background(), req, &headers, &trailers)
	if err != nil {
		return unexpectedError("PingList", err)
	}

	_, err = srv.Recv()
	if err == nil {
		return errors.Errorf("Got nil error, expected non-nil")
	}

	code, message := getStatus(err)
	if code != codes.Code(req.ErrorCodeReturned) {
		return reportError("code", code, codes.Code(req.ErrorCodeReturned))
	}
	// Message differs when connection is severed - this is OK
	if req.FailureType != DROP {
		if message != req.Value {
			return reportError("message", message, req.Value)
		}
	}

	return nil
}

func TestPingList(client TestClient, getStatus func(error) (codes.Code, string)) error {
	err := testPingList(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers and trailers")
	}

	err = testPingList(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers only")
	}

	err = testPingList(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send trailer only")
	}

	err = testPingList(client, &Request{
		Value:         "test",
		ResponseCount: 1,
	})
	if err != nil {
		return errors.WithMessage(err, "send nethier header or trailer")
	}

	req := &Request{
		Value:             "test",
		ResponseCount:     1,
		FailureType:       CODE,
		ErrorCodeReturned: uint32(codes.DataLoss),
	}
	err = testPingListError(client, req, getStatus)
	if err != nil {
		return errors.WithMessage(err, "trigger return code")
	}

	req.FailureType = DROP
	req.ErrorCodeReturned = uint32(codes.Unknown)
	req.Value = ""
	return errors.WithMessage(testPingListError(client, req, getStatus), "trigger network error")
}
