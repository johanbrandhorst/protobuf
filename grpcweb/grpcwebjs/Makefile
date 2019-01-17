all: checkout build

checkout:
	npm install grpc-web-client@0.6.3 webpack@4.28.4 webpack-cli@3.2.1

build:
	./node_modules/.bin/webpack

clean:
	rm -f grpc.inc.js
	rm -rf node_modules
