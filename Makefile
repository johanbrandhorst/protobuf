regenerate:
	cd ptypes && make regenerate
	cd protoc-gen-gopherjs && make regenerate
	cd test && make regenerate

install:
	cd protoc-gen-gopherjs && go install ./

tests:
	(cd protoc-gen-gopherjs && make tests) && \
	(cd test && make test)

docker:
	bash -c "\
		trap '\
			docker-compose logs selenium && \
			docker-compose logs chromedriver && \
			docker-compose down' EXIT; \
		docker-compose up -d && \
		docker-compose exec -T testrunner bash -c '\
            mkdir -p /go/src/github.com/johanbrandhorst/protobuf/' && \
		docker cp ./ testrunner:/go/src/github.com/johanbrandhorst/protobuf/ && \
		docker-compose exec -T testrunner bash -c '\
			cd /go/src/github.com/johanbrandhorst/protobuf && \
			go install ./vendor/github.com/onsi/ginkgo/ginkgo &&\
			cd test && make test'\
		"

rebuild:
	cd grpcweb/grpcwebjs && make build
