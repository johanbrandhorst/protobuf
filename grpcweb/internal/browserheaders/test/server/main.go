package main

import (
	"log"
	"net/http"

	"github.com/johanbrandhorst/protobuf/grpcweb/internal/browserheaders/test/shared"
)

func main() {
	log.Fatal(http.ListenAndServe(shared.ServerAddr, http.FileServer(http.Dir("./client/html"))))
}
