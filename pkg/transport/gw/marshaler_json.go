package gw

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/status"
)

// JSONBuiltin is a Marshaler which marshals/unmarshals into/from JSON
// with the standard "encoding/json" package of Golang.
// Although it is generally faster for simple proto messages than JSONPb,
// it does not support advanced features of protobuf, e.g. map, oneof, ....
//
// The NewEncoder and NewDecoder types return *json.Encoder and
// *json.Decoder respectively.
type JSONBuiltin struct {
	runtime.JSONBuiltin
	WrappCode bool
	CodeBase  int32
}

// Marshal marshals "v" into JSON
func (j *JSONBuiltin) Marshal(v interface{}) ([]byte, error) {

	// fmt.Println("marshal xxxx,", reflect.TypeOf(v))

	if s, ok := v.(*status.Status); !ok {
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
	} else {

		code := s.Code
		if code <= 16 {
			code = j.CodeBase + s.Code
		}

		v = struct {
			Code    int32  `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: s.Message,
		}

	}

	return j.JSONBuiltin.Marshal(v)
}
