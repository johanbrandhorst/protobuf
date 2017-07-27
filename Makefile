regenerate:
	cd ptypes && make regenerate
	cd protoc-gen-gopherjs && make regenerate
	cd test && make regenerate

install:
	cd protoc-gen-gopherjs && go install ./

tests:
	(cd protoc-gen-gopherjs/test && make test) && \
	(cd test && make regenerate && make test)

rebuild:
	cd grpcweb/grpcwebjs && make build
