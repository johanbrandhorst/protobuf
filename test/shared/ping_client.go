package shared

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type TestPingClientStreamClient interface {
	Send(*Request) error
	CloseAndRecv() (*Response, error)
	Trailer() metadata.MD
	Header() (metadata.MD, error)
}

func testPingClientStream(client TestClient, req *Request) error {
	headers, trailers := metadata.MD{}, metadata.MD{}
	ctx := context.Background()
	if req.CheckMetadata {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(ClientMDTestKey, ClientMDTestValue))
	}
	srv, err := client.PingClientStream(ctx, &headers, &trailers)
	if err != nil {
		return unexpectedError("PingClientStream", err)
	}

	iterations := 10
	req.Value = "test"
	req.ResponseCount = 1
	for i := 0; i < iterations; i++ {
		err := srv.Send(req)
		if err != nil {
			return unexpectedError("Send", err)
		}
	}

	resp, err := srv.CloseAndRecv()
	if err != nil {
		return unexpectedError("CloseAndRecv", err)
	}
	if resp.Value != "Closed" {
		return reportError("response value", resp.Value, "Closed")
	}
	if resp.Counter != 0 {
		return reportError("response counter", resp.Counter, 0)
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
				return reportError("trailer", h, trailer)
			} else if len(values) != iterations || values[0] != value {
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

type TestPingClientStreamErrorClient interface {
	Send(*Request) error
	CloseAndRecv() (*Response, error)
}

func testPingClientStreamError(client TestClient, req *Request, getStatus func(error) (codes.Code, string)) error {
	srv, err := client.PingClientStreamError(context.Background())
	if err != nil {
		return unexpectedError("PingClientStream", err)
	}

	if req.FailureType == CODE {
		// Send OK first
		err = srv.Send(&Request{
			Value:         req.Value,
			ResponseCount: req.ResponseCount,
		})
		if err != nil {
			return unexpectedError("Send", err)
		}

		// Trigger error
		err = srv.Send(req)
		if err != nil {
			return unexpectedError("Send", err)
		}

		// Shouldn't error
		err = srv.Send(&Request{
			Value:         req.Value,
			ResponseCount: req.ResponseCount,
		})
		if err != nil {
			return unexpectedError("Send", err)
		}
	}

	_, err = srv.CloseAndRecv()
	if err == nil {
		return errors.Errorf("Got nil error, expected non-nil")
	}

	code, message := getStatus(err)
	if code != codes.Code(req.ErrorCodeReturned) {
		return reportError("code", code, codes.Code(req.ErrorCodeReturned))
	}
	if message != req.Value {
		return reportError("message", message, req.Value)
	}

	return nil
}

func TestPingClientStream(client TestClient, getStatus func(error) (codes.Code, string)) error {
	err := testPingClientStream(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers and trailers")
	}

	err = testPingClientStream(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers only")
	}

	err = testPingClientStream(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send trailer only")
	}

	err = testPingClientStream(client, &Request{
		Value:         "test",
		ResponseCount: 1,
	})
	if err != nil {
		return errors.WithMessage(err, "send nethier header or trailer")
	}

	req := &Request{
		Value:             "error",
		ResponseCount:     1,
		ErrorCodeReturned: uint32(codes.Internal),
	}
	err = testPingClientStreamError(client, req, getStatus)
	if err != nil {
		return errors.WithMessage(err, "error after close")
	}

	req.FailureType = CODE
	req.ErrorCodeReturned = uint32(codes.DataLoss)
	req.Value = "test"
	err = testPingClientStreamError(client, req, getStatus)
	if err != nil {
		return errors.WithMessage(err, "trigger return code")
	}

	req.FailureType = DROP
	req.ErrorCodeReturned = uint32(codes.Internal)
	req.Value = "error"
	return errors.WithMessage(testPingClientStreamError(client, req, getStatus), "trigger network error")
}
