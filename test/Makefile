regenerate:
	# Generate protofiles
	protoc proto/test/test.proto \
		--gopherjs_out=plugins=grpc,Mgoogle/protobuf/empty.proto=github.com/johanbrandhorst/protobuf/ptypes/empty:./client/ \
		--go_out=plugins=grpc:./server/

	# Regenerate dependencies
	(cd ../protoc-gen-gopherjs/test && make regenerate)

	# Generate GopherJS files
	go generate ./client/...

test:
	ginkgo -v .
