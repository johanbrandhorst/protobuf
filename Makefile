regenerate:
	cd ptypes && make regenerate
	cd protoc-gen-gopherjs && make regenerate
	cd test && make regenerate

install:
	cd protoc-gen-gopherjs && go install ./

tests:
	cd test && make regenerate && make test
