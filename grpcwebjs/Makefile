all: checkout build

checkout:
	git submodule update --init ./grpc-web
	npm install google-closure-compiler google-protobuf @types/google-protobuf grpc-web-client

build:
	java -jar node_modules/google-closure-compiler/compiler.jar \
		--dependency_mode STRICT \
		--use_types_for_optimization \
		--module_resolution NODE \
		--formatting PRETTY_PRINT \
		--process_common_js_modules \
		--js_output_file grpc-web.inc.js \
		--common_js_entry_module="./node_modules/grpc-web-client/dist/grpc.js" \
		--js ./node_modules/grpc-web-client/**.js \
		--js ./node_modules/path/path.js \
		--js ./node_modules/util/util.js \
		--js ./node_modules/util/node_modules/inherits/inherits_browser.js \
		--js ./node_modules/util/support/isBuffer.js --help

clean:
	rm grpc-web.inc.js
