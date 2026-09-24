package gw

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	*http.Server
	ctx     context.Context
	address string
	router  *gin.Engine
}

// type ServeMux struct {
// 	*mux.Router
// 	*runtime.ServeMux
// }

// func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	mux.ServeMux.ServeHTTP(w, r)
// }

type OptionFunc func(*option)
type option struct {
	WrapCode            bool
	CodeBase            int32
	errHandler          runtime.ErrorHandlerFunc
	routingErrorHandler runtime.RoutingErrorHandlerFunc
	metadataAnnotator   func(context.Context, *http.Request) metadata.MD
	matcher             runtime.HeaderMatcherFunc
}

func WithWrapCode() OptionFunc {
	return func(o *option) {
		o.WrapCode = true
	}
}

func WithCodeBase(base int32) OptionFunc {
	return func(o *option) {
		if base <= 0 {
			return
		}
		o.CodeBase = base
	}
}

func WithMetadata(annotator func(context.Context, *http.Request) metadata.MD) OptionFunc {
	return func(o *option) {
		if annotator == nil {
			return
		}
		o.metadataAnnotator = annotator
	}
}

func WithHeaderMatcher(matcher runtime.HeaderMatcherFunc) OptionFunc {
	return func(o *option) {
		if matcher == nil {
			return
		}
		o.matcher = matcher
	}
}

func WithErrorHandler(errhandler runtime.ErrorHandlerFunc) OptionFunc {
	return func(o *option) {
		if errhandler == nil {
			return
		}
		o.errHandler = errhandler
	}
}

func WithRoutingErrorHandler(routingErrorHandler runtime.RoutingErrorHandlerFunc) OptionFunc {
	return func(o *option) {

		if routingErrorHandler == nil {
			return
		}
		o.routingErrorHandler = routingErrorHandler
	}
}

func NewServer(ctx context.Context, address string) *Server {

	// smux := &ServeMux{}
	// r := mux.NewRouter()
	// smux.Router = r
	// smux.ServeMux = gwmux

	// r.PathPrefix("/").Handler(gwmux)
	// mux := http.NewServeMux()
	// mux.Handle("/", gwmux)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("*")
	router.Use(cors.New(config))

	srv := Server{
		ctx:     ctx,
		address: address,
		router:  router,
		// Server: &http.Server{
		// 	Addr: address,
		// 	Handler: allowCORS(router),
		// },
	}

	srv.Server = &http.Server{
		Addr:    srv.address,
		Handler: router,
	}
	return &srv
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) Start(ctx context.Context) error {

	s.ctx = ctx
	slog.Info("[GW] server listening on: ", "address", s.address)
	// return s.Run(s.address)
	return s.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("[GW] server stopping")
	return s.Shutdown(ctx)
}
