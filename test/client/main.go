package main

import (
	"context"
	"fmt"
	"math"
	"net/url"

	"github.com/go-test/deep"
	"github.com/gopherjs/gopherjs/js"
	"github.com/rusco/qunit"
	"google.golang.org/grpc/codes"
	"honnef.co/go/js/dom"

	"github.com/johanbrandhorst/protobuf/grpcweb/status"
	grpctest "github.com/johanbrandhorst/protobuf/grpcweb/test"
	gentest "github.com/johanbrandhorst/protobuf/protoc-gen-gopherjs/test"
	"github.com/johanbrandhorst/protobuf/protoc-gen-gopherjs/test/multi"
	"github.com/johanbrandhorst/protobuf/protoc-gen-gopherjs/test/types"
	"github.com/johanbrandhorst/protobuf/ptypes/empty"
	"github.com/johanbrandhorst/protobuf/test/client/proto/test"
	"github.com/johanbrandhorst/protobuf/test/client/wrappers"
	"github.com/johanbrandhorst/protobuf/test/recoverer"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

//go:generate gopherjs build main.go -o html/index.js

var uri string

func init() {
	u, err := url.Parse(dom.GetWindow().Document().BaseURI())
	if err != nil {
		panic(err)
	}
	uri = u.Scheme + "://" + u.Hostname()
}

func typeTests() {
	qunit.Module("Integration Types tests")

	qunit.Test("PingRequest Marshal and Unmarshal", func(assert qunit.QUnitAssert) {
		req := &test.PingRequest{
			Value:             "1234",
			ResponseCount:     10,
			ErrorCodeReturned: 1,
			FailureType:       test.PingRequest_CODE,
			CheckMetadata:     true,
			SendHeaders:       true,
			SendTrailers:      true,
			MessageLatencyMs:  100,
		}

		marshalled := req.Marshal()
		newReq, err := new(test.PingRequest).Unmarshal(marshalled)
		if err != nil {
			assert.Ok(false, "Unexpected error returned: "+err.Error()+"\n"+err.(*js.Error).Stack())
		}
		assert.DeepEqual(req, newReq, "Marshalling and unmarshalling results in the same struct")
	})

	qunit.Test("ExtraStuff Marshal and Unmarshal", func(assert qunit.QUnitAssert) {
		req := &test.ExtraStuff{
			Addresses: map[int32]string{
				1234: "The White House",
				5678: "The Empire State Building",
			},
			Title: &test.ExtraStuff_FirstName{
				FirstName: "Allison",
			},
			CardNumbers: []uint32{
				1234, 5678,
			},
		}

		marshalled := req.Marshal()
		newReq, err := new(test.ExtraStuff).Unmarshal(marshalled)
		if err != nil {
			assert.Ok(false, "Unexpected error returned: "+err.Error()+"\n"+err.(*js.Error).Stack())
		}
		assert.DeepEqual(req, newReq, "Marshalling and unmarshalling results in the same struct")
	})
}

func serverTests(label, serverAddr, emptyServerAddr string) {
	qunit.Module(fmt.Sprintf("%s Integration tests", label))

	c := test.NewTestServiceClient(uri + serverAddr)
	w := wrappers.ClientWrapper{C: c}
	getStatus := func(err error) (codes.Code, string) {
		st := status.FromError(err)
		return st.Code, st.Message
	}

	qunit.AsyncTest("Unary call to empty server", func() interface{} {
		c := test.NewTestServiceClient(uri + emptyServerAddr)

		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			_, err := c.PingEmpty(context.Background(), &empty.Empty{})
			if err == nil {
				qunit.Ok(false, "Expected error, returned nil")
				return
			}

			st := status.FromError(err)
			if st.Message != "unknown service test.TestService" {
				qunit.Ok(false, "Unexpected error, saw "+st.Message)
			}

			qunit.Ok(true, "Error was as expected")
		}()

		return nil
	})

	qunit.AsyncTest("Unary call to echo server with many types", func() interface{} {
		c := types.NewEchoServiceClient(uri + serverAddr)
		req := &types.TestAllTypes{
			SingleInt32:       -math.MaxInt32,
			SingleInt64:       -math.MaxInt64,
			SingleUint32:      math.MaxUint32,
			SingleUint64:      math.MaxUint64,
			SingleSint32:      -math.MaxInt32,
			SingleSint64:      -math.MaxInt64,
			SingleFixed32:     math.MaxUint32,
			SingleFixed64:     math.MaxUint64,
			SingleSfixed32:    -math.MaxInt32,
			SingleSfixed64:    -math.MaxInt64,
			SingleFloat:       math.MaxFloat32,
			SingleDouble:      math.MaxFloat64,
			SingleBool:        true,
			SingleString:      "Alfred",
			SingleBytes:       []byte("Megan"),
			SingleNestedEnum:  types.TestAllTypes_BAR,
			SingleForeignEnum: types.ForeignEnum_FOREIGN_BAR,
			SingleImportedMessage: &multi.Multi1{
				Color:   multi.Multi2_GREEN,
				HatType: multi.Multi3_FEDORA,
			},
			SingleNestedMessage: &types.TestAllTypes_NestedMessage{
				B: math.MaxInt32,
			},
			SingleForeignMessage: &types.ForeignMessage{
				C: math.MaxInt32,
			},
			RepeatedInt32:       []int32{-math.MaxInt32, math.MaxInt32},
			RepeatedInt64:       []int64{-math.MaxInt64, math.MaxInt64},
			RepeatedUint32:      []uint32{0, math.MaxUint32},
			RepeatedUint64:      []uint64{0, math.MaxUint64},
			RepeatedSint32:      []int32{-math.MaxInt32, math.MaxInt32},
			RepeatedSint64:      []int64{-math.MaxInt64, math.MaxInt64},
			RepeatedFixed32:     []uint32{0, math.MaxUint32},
			RepeatedFixed64:     []uint64{0, math.MaxUint64},
			RepeatedSfixed32:    []int32{-math.MaxInt32, math.MaxInt32},
			RepeatedSfixed64:    []int64{-math.MaxInt64, math.MaxInt64},
			RepeatedFloat:       []float32{-math.MaxFloat32, math.MaxFloat32},
			RepeatedDouble:      []float64{-math.MaxFloat64, math.MaxFloat64},
			RepeatedBool:        []bool{true, false, true},
			RepeatedString:      []string{"Alfred", "Robin", "Simon"},
			RepeatedBytes:       [][]byte{[]byte("David"), []byte("Henrik")},
			RepeatedNestedEnum:  []types.TestAllTypes_NestedEnum{types.TestAllTypes_BAR, types.TestAllTypes_BAZ},
			RepeatedForeignEnum: []types.ForeignEnum{types.ForeignEnum_FOREIGN_BAR, types.ForeignEnum_FOREIGN_BAZ},
			RepeatedImportedMessage: []*multi.Multi1{
				{
					Color:   multi.Multi2_RED,
					HatType: multi.Multi3_FEZ,
				},
				{
					Color:   multi.Multi2_GREEN,
					HatType: multi.Multi3_FEDORA,
				},
			},
			RepeatedNestedMessage: []*types.TestAllTypes_NestedMessage{
				{
					B: -math.MaxInt32,
				},
				{
					B: math.MaxInt32,
				},
			},
			RepeatedForeignMessage: []*types.ForeignMessage{
				{
					C: -math.MaxInt32,
				},
				{
					C: math.MaxInt32,
				},
			},
			OneofField: &types.TestAllTypes_OneofImportedMessage{
				OneofImportedMessage: &multi.Multi1{
					Multi2: &multi.Multi2{
						RequiredValue: math.MaxInt32,
						Color:         multi.Multi2_BLUE,
					},
					Color:   multi.Multi2_RED,
					HatType: multi.Multi3_FEDORA,
				},
			},
		}

		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			resp, err := c.EchoAllTypes(context.Background(), req)
			if err != nil {
				st := status.FromError(err)
				qunit.Ok(false, "Unexpected error:"+st.Error())
				return
			}
			if diff := deep.Equal(req, resp); diff != nil {
				var s string
				for _, v := range diff {
					s += "\n" + v
				}
				qunit.Ok(false, s)
				return
			}

			qunit.Ok(true, "Request and Response matched")
		}()

		return nil
	})

	qunit.AsyncTest("Unary call to echo server with many maps", func() interface{} {
		c := types.NewEchoServiceClient(uri + serverAddr)
		req := &types.TestMap{
			MapInt32Int32: map[int32]int32{
				1: 2,
				3: 4,
			},
			MapInt64Int64: map[int64]int64{
				5: 6,
				7: 8,
			},
			MapUint32Uint32: map[uint32]uint32{
				9:  10,
				11: 12,
			},
			MapUint64Uint64: map[uint64]uint64{
				13: 14,
				15: 16,
			},
			MapSint32Sint32: map[int32]int32{
				17: 18,
				19: 20,
			},
			MapSint64Sint64: map[int64]int64{
				21: 22,
				23: 24,
			},
			MapFixed32Fixed32: map[uint32]uint32{
				25: 26,
				27: 28,
			},
			MapFixed64Fixed64: map[uint64]uint64{
				29: 30,
				31: 32,
			},
			MapSfixed32Sfixed32: map[int32]int32{
				33: 34,
				35: 36,
			},
			MapSfixed64Sfixed64: map[int64]int64{
				37: 38,
				39: 40,
			},
			MapInt32Float: map[int32]float32{
				41:  42.41,
				432: 44.43,
			},
			MapInt32Double: map[int32]float64{
				45: 46.45,
				47: 48.47,
			},
			MapBoolBool: map[bool]bool{
				true:  false,
				false: false,
			},
			MapStringString: map[string]string{
				"Henrik": "David",
				"Simon":  "Robin",
			},
			MapInt32Bytes: map[int32][]byte{
				49: []byte("Astrid"),
				50: []byte("Ebba"),
			},
			MapInt32Enum: map[int32]types.MapEnum{
				51: types.MapEnum_MAP_ENUM_BAR,
				52: types.MapEnum_MAP_ENUM_BAZ,
			},
			MapInt32ForeignMessage: map[int32]*types.ForeignMessage{
				53: {C: 54},
				55: {C: 56},
			},
			MapInt32ImportedMessage: map[int32]*multi.Multi1{
				57: {
					Multi2: &multi.Multi2{
						RequiredValue: 58,
						Color:         multi.Multi2_RED,
					},
					Color:   multi.Multi2_GREEN,
					HatType: multi.Multi3_FEZ,
				},
				59: {
					Color:   multi.Multi2_BLUE,
					HatType: multi.Multi3_FEDORA,
				},
			},
		}

		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			resp, err := c.EchoMaps(context.Background(), req)
			if err != nil {
				st := status.FromError(err)
				qunit.Ok(false, "Unexpected error:"+st.Error())
				return
			}
			if diff := deep.Equal(req, resp); diff != nil {
				var s string
				for _, v := range diff {
					s += "\n" + v
				}
				qunit.Ok(false, s)
				return
			}

			qunit.Ok(true, "Request and Response matched")
		}()

		return nil
	})

	qunit.AsyncTest("Unary server call", func() interface{} {
		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			err := shared.TestPing(w, getStatus)
			if err != nil {
				qunit.Ok(false, err.Error())
				return
			}

			qunit.Ok(true, "Request succeeded")
		}()

		return nil
	})

	qunit.AsyncTest("Server Streaming call", func() interface{} {
		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			err := shared.TestPingList(w, getStatus)
			if err != nil {
				qunit.Ok(false, err.Error())
				return
			}

			qunit.Ok(true, "Request succeeded")
		}()

		return nil
	})

	qunit.AsyncTest("Client Streaming call", func() interface{} {
		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			err := shared.TestPingClientStream(w, getStatus)
			if err != nil {
				qunit.Ok(false, err.Error())
				return
			}

			qunit.Ok(true, "Request succeeded")
		}()

		return nil
	})

	qunit.AsyncTest("Bi-directional streaming call", func() interface{} {
		go func() {
			defer recoverer.Recover() // recovers any panics and fails tests
			defer qunit.Start()

			err := shared.TestPingBidiStream(w, getStatus)
			if err != nil {
				qunit.Ok(false, err.Error())
				return
			}

			qunit.Ok(true, "Request succeeded")
		}()

		return nil
	})
}

func main() {
	defer recoverer.Recover() // recovers any panics and fails tests

	typeTests()
	serverTests("HTTP2", shared.HTTP2Server, shared.EmptyHTTP2Server)
	serverTests("HTTP1", shared.HTTP1Server, shared.EmptyHTTP1Server)

	// protoc-gen-gopherjs tests
	gentest.GenTypesTest()

	// grpcweb tests
	grpctest.GRPCWebTest()
}
