package helloworld_test

import (
	"bytes"
	"context"
	io "io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/bakins/protoc-gen-grpc-transcode/internal/example"
)

func TestHelloWorld(t *testing.T) {
	tests := map[string]struct {
		grpcError       error
		expectedMessage string
		expectedCode    codes.Code
		expectedStatus  int
	}{
		"ok": {
			expectedStatus:  http.StatusOK,
			expectedMessage: "hello world",
		},
		"unauthenticated": {
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "test error message",
			expectedCode:    codes.Unauthenticated,
			grpcError:       status.Error(codes.Unauthenticated, "test error message"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := server{
				err: test.grpcError,
			}

			g := grpc.NewServer()
			pb.RegisterGreeterServer(g, &s)
			grpcServer := httptest.NewServer(h2c.NewHandler(g, &http2.Server{}))
			defer grpcServer.Close()

			u, err := url.Parse(grpcServer.URL)
			require.NoError(t, err)

			conn, err := grpc.Dial(u.Host, grpc.WithInsecure())
			require.NoError(t, err)

			transcoder := pb.NewGreeterTranscode(conn)

			svr := httptest.NewServer(transcoder)
			defer svr.Close()

			req, err := http.NewRequest(
				http.MethodPost,
				svr.URL+"/helloworld.Greeter/SayHello",
				bytes.NewBufferString(`{"name": "world"}`),
			)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, test.expectedStatus, resp.StatusCode)
			require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			data, err := io.ReadAll(resp.Body)

			require.NoError(t, err)

			if test.expectedStatus == http.StatusOK {
				var response pb.HelloReply
				err := protojson.Unmarshal(data, &response)
				require.NoError(t, err)
				require.Equal(t, test.expectedMessage, response.Message)
			}

			var st spb.Status
			err = protojson.Unmarshal(data, &st)
			require.NoError(t, err)

			require.Equal(t, test.expectedCode, codes.Code(st.Code))
			require.Equal(t, test.expectedMessage, st.Message)
		})
	}
}

type server struct {
	pb.UnimplementedGreeterServer
	err error
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &pb.HelloReply{Message: "hello " + in.GetName()}, nil
}
