package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bakins/protoc-gen-grpc-transcode/internal/example"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	transcoder := pb.NewGreeterTranscode(conn)

	mux := http.NewServeMux()

	// Methods returns a list of methods in the form
	// /<package>.<service>/<method> that is suitable
	// for using in a mux
	for _, method := range transcoder.Methods() {
		mux.Handle(method, transcoder)
	}

	http.ListenAndServe(":8080", mux)
}
