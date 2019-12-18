package interceptor

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientHeaderInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return errors.New("metadata.FromOutgoingContext failed!")
		}
		md["http-client"] = []string{"interceptor/0.1"}
		ctx = metadata.NewOutgoingContext(ctx, md)
		log.Infof("%+v", md)
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
