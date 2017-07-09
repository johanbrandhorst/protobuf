package main

import (
	"github.com/johanbrandhorst/protobuf/test/client/proto/test"
	"github.com/rusco/qunit"
)

//go:generate gopherjs build main.go -m -o html/index.js

func main() {
	qunit.Module("grpcweb")
	qunit.Test("Type constructors", func(assert qunit.QUnitAssert) {
		qunit.Expect(15)

		req := new(test.PingRequest).New("1234", 10, 1, test.PingRequest_CODE, true, true, true, 100)
		assert.Equal(req.GetValue(), "1234", "Value is set as expected")
		assert.Equal(req.GetResponseCount(), 10, "ResponseCount is set as expected")
		assert.Equal(req.GetErrorCodeReturned(), 1, "ErrorCodeReturned is set as expected")
		assert.Equal(req.GetFailureType(), test.PingRequest_CODE, "ErrorCodeReturned is set as expected")
		assert.Equal(req.GetCheckMetadata(), true, "CheckMetadata is set as expected")
		assert.Equal(req.GetSendHeaders(), true, "SendHeaders is set as expected")
		assert.Equal(req.GetSendTrailers(), true, "SendTrailers is set as expected")
		assert.Equal(req.GetMessageLatencyMs(), 100, "MessageLatencyMs is set as expected")

		es := new(test.ExtraStuff).New(
			map[int32]string{1234: "The White House", 5678: "The Empire State Building"},
			&test.ExtraStuff_FirstName{FirstName: "Allison"},
			[]uint32{1234, 5678})
		addrs := es.GetAddresses()
		assert.Equal(addrs[1234], "The White House", "Address 1234 is set as expected")
		assert.Equal(addrs[5678], "The Empire State Building", "Address 5678 is set as expected")
		crdnrs := es.GetCardNumbers()
		assert.Equal(crdnrs[0], 1234, "CardNumber #1 is set as expected")
		assert.Equal(crdnrs[1], 5678, "CardNumber #2 is set as expected")
		assert.Equal(es.GetFirstName(), "Allison", "FirstName is set as expected")
		assert.Equal(es.GetIdNumber(), 0, "IdNumber is not set, as expected")
		assert.Equal(
			es.GetTitle().(*test.ExtraStuff_FirstName).FirstName,
			"Allison",
			"Title is set as expected")
	})
}
