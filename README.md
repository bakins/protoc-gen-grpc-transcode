# protoc-gen-grpc-transcode

`protoc-gen-grpc-transcode` is a protoc plugin for generating a simple HTTP POST+json
transcoder for grpc in Go.

The generated code does not depend on any packages in this repo.

It is meant to be used when along side of or instead of [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) when you need simple json to grpc transcoding.

## Usage

Install the plugin:

```shell
go install github.com/bakins/protoc-gen-grpc-transcode@latest
```

Generate the server code (using example in [internal/example](./internal/example/)):

```shell
protoc -I . \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-transcode_out=. \
    --grpc-transcode_opt=paths=source_relative \
    *.proto
```

Now use in your program

```go
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
```

The test from another terminal:

```shell
$ curl -H "Content-Type: application/json" --data-binary '{"name": "world" }' http://127.0.0.1:8080/helloworld.Greeter/SayHello

{"message": "hello world"}
```

## LICENSE

See [LICENSE](./LICENSE)

## TODO
- more test cases
