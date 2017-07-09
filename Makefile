regenerate:
	cd ptypes && make regenerate
	cd protoc-gen-gopherjs && make regenerate

install:
	cd protoc-gen-gopherjs && go install ./
