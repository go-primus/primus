package gw

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func CustomMatcher(key string) (string, bool) {

	// slog.Warn("custom key:", key)
	switch strings.ToLower(key) {
	case "token":
		return "token", true
	case "user-id", "userid", "user_id":
		return "user_id", true
	case "devicemodel", "device_model", "device-model":
		return "device_model", true
	case "devicename", "device_name", "device-name":
		return "device_name", true
	case "devicetype", "device_type", "device-type":
		return "device_type", true
	case "deviceid", "device_id", "device-id":
		return "device_id", true
	case "authori-zation":
		return "Authorization", true
	default:
		return runtime.DefaultHeaderMatcher(key) // Grpc-Metadata-UseNamedPath  --> UseNamedPath
	}
}
