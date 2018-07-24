all: checkout build

checkout:
	npm install protobufjs webpack webpack-cli long

build:
	./node_modules/.bin/webpack-cli

clean:
	rm decodeio.inc.js
	rm -rf node_modules
