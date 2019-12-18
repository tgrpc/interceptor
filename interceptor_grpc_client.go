package interceptor

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientHeaderInterceptor(header http.Header) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := HeaderToMD(header)
		md["http-client"] = []string{"v0.1"}
		log.Infof("%+v", md)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func HeaderToMD(header http.Header) metadata.MD {
	kv := make(map[string]string, len(header))
	for key, vals := range header {
		kv[key] = strings.Join(vals, ";")
	}
	return metadata.New(kv)
}
