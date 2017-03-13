// Copyright 2017 Johan Brandhorst. All Rights Reserved.
// See LICENSE for licensing terms.

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/johanbrandhorst/protoc-gen-gopherjs/filegenerator"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln("Could not read request from STDIN: ", err)
	}

	req := &plugin.CodeGeneratorRequest{}

	err = proto.Unmarshal(data, req)
	if err != nil {
		log.Fatalln("Could not unmarshal request: ", err)
	}

	resp := &plugin.CodeGeneratorResponse{}

	for _, inFile := range req.GetProtoFile() {
		for _, reqFile := range req.GetFileToGenerate() {
			if inFile.GetName() == reqFile {
				outFile, err := processFile(inFile)
				if err != nil {
					log.Fatalln("Could not process file: ", err)
				}
				resp.File = append(resp.File, outFile)
			}
		}
	}

	data, err = proto.Marshal(resp)
	if err != nil {
		log.Fatalf("Could not marshal response: %v [%v]\n", err, resp)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatalln("Could not write response to STDOUT: ", err)
	}
}

func processFile(inFile *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	outFile := &plugin.CodeGeneratorResponse_File{}
	outFile.Name = proto.String(strings.TrimSuffix(inFile.GetName(), ".proto") + ".gopherjs.pb.go")

	b := &bytes.Buffer{}
	fg := filegenerator.New(b)

	fg.Generate(inFile)

	outFile.Content = proto.String(b.String())

	return outFile, nil
}
