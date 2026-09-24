package gw

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func WithHealthEndpointAt(healthCheckClient grpc_health_v1.HealthClient, method string, endpointPath string) runtime.ServeMuxOption {
	return func(s *runtime.ServeMux) {
		// error can be ignored since pattern is definitely valid
		_ = s.HandlePath(
			method, endpointPath, func(w http.ResponseWriter, r *http.Request, _ map[string]string,
			) {
				_, outboundMarshaler := runtime.MarshalerForRequest(s, r)
				annotatedContext, err := runtime.AnnotateContext(r.Context(), s, r, grpc_health_v1.Health_Check_FullMethodName, runtime.WithHTTPPathPattern(endpointPath))
				if err != nil {
					runtime.HTTPError(r.Context(), s, outboundMarshaler, w, r, err)
					// s.errorHandler(r.Context(), s, outboundMarshaler, w, r, err)
					return
				}

				var md runtime.ServerMetadata
				resp, err := healthCheckClient.Check(annotatedContext, &grpc_health_v1.HealthCheckRequest{
					Service: r.URL.Query().Get("service"),
				}, grpc.Header(&md.HeaderMD), grpc.Trailer(&md.TrailerMD))
				annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
				if err != nil {

					runtime.HTTPError(annotatedContext, s, outboundMarshaler, w, r, err)
					// s.errorHandler(annotatedContext, s, outboundMarshaler, w, r, err)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
					switch resp.GetStatus() {
					case grpc_health_v1.HealthCheckResponse_NOT_SERVING, grpc_health_v1.HealthCheckResponse_UNKNOWN:
						err = status.Error(codes.Unavailable, resp.String())
					case grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN:
						err = status.Error(codes.NotFound, resp.String())
					}

					runtime.HTTPError(annotatedContext, s, outboundMarshaler, w, r, err)
					// s.errorHandler(annotatedContext, s, outboundMarshaler, w, r, err)
					return
				}

				_ = outboundMarshaler.NewEncoder(w).Encode(resp)
			})
	}
}
