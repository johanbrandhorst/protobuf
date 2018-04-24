# GopherJS Bindings for ProtobufJS and gRPC-Web
[![Circle CI](https://circleci.com/gh/johanbrandhorst/protobuf.svg?style=shield)](https://circleci.com/gh/johanbrandhorst/protobuf)
[![Go Report Card](https://goreportcard.com/badge/github.com/johanbrandhorst/protobuf)](https://goreportcard.com/report/github.com/johanbrandhorst/protobuf)
[![GoDoc](https://godoc.org/github.com/johanbrandhorst/protobuf?status.svg)](https://godoc.org/github.com/johanbrandhorst/protobuf)
[![Gitter chat](https://badges.gitter.im/johanbrandhorst/protobuf.png)](https://gitter.im/gopherjs-protobuf)

![gRPC-Web radio operator Gopher by Egon Elbre (@egonelbre)](./logo.svg)
_gRPC-Web radio operator Gopher by Egon Elbre (@egonelbre)_

## Users
A list of some of the users of the library. Send
me a message on @jbrandhorst on Gophers Slack if you wish
to be added to this list

* https://github.com/anxiousmodernman/co-chair
* https://github.com/google/shenzhen-go

## Getting started
The easiest way to get started with gRPC-Web for Go is to clone
[the boilerplate repo](https://github.com/johanbrandhorst/grpcweb-boilerplate)
and start playing around with it.

## Components
### [GopherJS Protobuf Generator](./protoc-gen-gopherjs/README.md)
This is a GopherJS client code generator for the Google Protobuf format.
It generates code for interfacing with any gRPC services exposing a
gRPC-Web spec compatible interface. It uses `jspb` and `grpcweb`.
It is the main entrypoint for using the protobuf/gRPC GopherJS bindings.

### [GopherJS ProtobufJS Bindings](./jspb/README.md)
This is a simple GopherJS binding around the npm `google-protobuf` package.
Importing it into any GopherJS source allows usage of ProtobufJS functionality.

### [GopherJS gRPC-Web Client Bindings](./grpcweb/README.md)
This is a GopherJS binding around the Improbable gRPC-Web client.
It is not intended for public use.

## Contributions
Contributions are very welcome, please submit issues or PRs for review.

## Demo
See [the example repo](https://github.com/johanbrandhorst/grpcweb-example)
and [the demo website](https://grpcweb.jbrandhorst.com)
for an example use of the Protobuf and gRPC-Web bindings.
