package gw

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type JSONPB struct {
	runtime.JSONPb
	WrappCode bool
}

// Marshal marshals "v" into JSON
func (j *JSONPB) Marshal(v interface{}) ([]byte, error) {

	// fmt.Println("marshal xxxx,", reflect.TypeOf(v))

	if _, ok := v.(*status.Status); !ok {
		if j.WrappCode {

			v = struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Payload any    `json:"payload"`
			}{
				Code:    0,
				Message: "OK",
				Payload: v,
			}

		}
	}

	return j.JSONPb.Marshal(v)
}
