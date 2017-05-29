checkout:
	git submodule update --init protobuf
	npm install google-closure-compiler google-closure-library

build:
	java -jar node_modules/google-closure-compiler/compiler.jar \
		--dependency_mode STRICT \
		--use_types_for_optimization \
		--entry_point="goog:jspb.Message" \
		--entry_point="goog:jspb.BinaryReader" \
		--entry_point="goog:jspb.BinaryWriter" \
		--entry_point="goog:jspb.Map" \
		--js ./protobuf/js/**.js \
		--js ./protobuf/js/binary/**.js \
		--js ./node_modules/google-closure-library/closure/**.js \
		--js_output_file protobuf.inc.js
	echo '$$global.jspb = jspb;' >> protobuf.inc.js
