package gw

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
)

func (s *Server) WrapServeMux(grpcAddr string, fn func(*runtime.ServeMux, *grpc.ClientConn) error, opts ...OptionFunc) {

	option := &option{}
	for _, opt := range opts {
		opt(option)
	}

	conn, err := grpc.DialContext(
		s.ctx,
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	redirect := func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {

		headers := w.Header()
		if location, ok := headers["Grpc-Metadata-Location"]; ok {
			w.Header().Set("Location", location[0])
			w.WriteHeader(http.StatusFound)
		}

		return nil
	}

	muxOpts := []runtime.ServeMuxOption{}
	// muxOpts = append(muxOpts, runtime.WithIncomingHeaderMatcher(CustomMatcher))
	muxOpts = append(muxOpts, runtime.WithForwardResponseOption(redirect))
	muxOpts = append(muxOpts, runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(conn), "/health"))
	muxOpts = append(muxOpts, WithHealthEndpointAt(grpc_health_v1.NewHealthClient(conn), http.MethodHead, "/health"))
	muxOpts = append(muxOpts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &JSONBuiltin{WrappCode: option.WrapCode, CodeBase: option.CodeBase}))
	// muxOpts = append(muxOpts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
	// 	MarshalOptions: protojson.MarshalOptions{
	// 		UseProtoNames:  true,
	// 		UseEnumNumbers: true, // 枚举作为数字
	// 		// EmitUnpopulated: true, // 显示未填充的字段
	// 	},
	// }))
	// Marshaler: &JSONPb{
	// 	MarshalOptions: protojson.MarshalOptions{
	// 		EmitUnpopulated: true,
	// 	},
	// 	UnmarshalOptions: protojson.UnmarshalOptions{
	// 		DiscardUnknown: true,
	// 	},
	// },
	if option.metadataAnnotator != nil {

		muxOpts = append(muxOpts, runtime.WithMetadata(option.metadataAnnotator))
	}

	if option.matcher != nil {
		muxOpts = append(muxOpts, runtime.WithIncomingHeaderMatcher(option.matcher))
	}

	if option.errHandler != nil {
		muxOpts = append(muxOpts, runtime.WithErrorHandler(option.errHandler))
	}

	if option.routingErrorHandler != nil {

		muxOpts = append(muxOpts, runtime.WithRoutingErrorHandler(option.routingErrorHandler))
	}

	gwmux := runtime.NewServeMux(
		muxOpts...,
	)

	err = fn(gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	s.router.NoRoute(gin.WrapF(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		gwmux.ServeHTTP(writer, request)
	}))
	// s.router.Any("/v1/*path", gin.WrapH(gwmux))
}
