package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func testPing(client TestClient, req *Request) error {
	headers, trailers := metadata.MD{}, metadata.MD{}
	ctx := context.Background()
	if req.CheckMetadata {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(ClientMDTestKey, ClientMDTestValue))
	}
	resp, err := client.Ping(ctx, req, &headers, &trailers)
	if err != nil {
		return unexpectedError("Ping", err)
	}
	if resp.Value != req.Value {
		return reportError("response value", resp.Value, req.Value)
	}
	if resp.Counter != req.ResponseCount {
		return reportError("response counter", resp.Counter, req.ResponseCount)
	}

	if req.SendHeaders {
		for header, value := range map[string]string{
			ServerMDTestKey1: ServerMDTestValue1,
			ServerMDTestKey2: ServerMDTestValue2,
		} {
			if values, ok := headers[strings.ToLower(header)]; !ok {
				return reportError("header", headers, header)
			} else if len(values) != 1 || values[0] != value {
				return reportError("header value", values, value)
			}
		}
	} else {
		for _, header := range []string{ServerMDTestKey1, ServerMDTestKey2} {
			_, ok := headers[header]
			if ok {
				return reportError("unexpected header", headers[ServerMDTestKey1], "")
			}
		}
	}

	if req.SendTrailers {
		for trailer, value := range map[string]string{
			ServerTrailerTestKey1: ServerMDTestValue1,
			ServerTrailerTestKey2: ServerMDTestValue2,
		} {
			if values, ok := trailers[strings.ToLower(trailer)]; !ok {
				return reportError("trailer", trailers, trailer)
			} else if len(values) != 1 || values[0] != value {
				return reportError("trailer value", values, value)
			}
		}
	} else {
		for _, trailer := range []string{ServerMDTestKey1, ServerMDTestKey2} {
			_, ok := trailers[trailer]
			if ok {
				return reportError("unexpected trailer", trailers[ServerMDTestKey1], "")
			}
		}
	}

	return nil
}

func testPingError(client TestClient, req *Request, getStatus func(error) (codes.Code, string)) error {
	err := client.PingError(context.Background(), req)
	if err == nil {
		return errors.Errorf("Got nil error, expected non-nil")
	}

	code, message := getStatus(err)
	if code != codes.Code(req.ErrorCodeReturned) {
		fmt.Println(message)
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

func TestPing(client TestClient, getStatus func(error) (codes.Code, string)) error {
	err := testPing(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers and trailers")
	}

	err = testPing(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendHeaders:   true,
	})
	if err != nil {
		return errors.WithMessage(err, "send headers only")
	}

	err = testPing(client, &Request{
		Value:         "test",
		ResponseCount: 1,
		SendTrailers:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "send trailer only")
	}

	err = testPing(client, &Request{
		Value:         "test",
		ResponseCount: 1,
	})
	if err != nil {
		return errors.WithMessage(err, "send neither header or trailer")
	}

	req := &Request{
		Value:             "test",
		ResponseCount:     1,
		FailureType:       CODE,
		ErrorCodeReturned: uint32(codes.DataLoss),
	}
	err = testPingError(client, req, getStatus)
	if err != nil {
		return errors.WithMessage(err, "trigger return code")
	}

	req.FailureType = DROP
	req.ErrorCodeReturned = uint32(codes.Unknown)
	req.Value = ""
	return errors.WithMessage(testPingError(client, req, getStatus), "trigger network error")
}
