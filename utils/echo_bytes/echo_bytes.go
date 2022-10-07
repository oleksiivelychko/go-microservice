package echo_bytes

import (
	"bytes"
	"encoding/json"
)

func EchoBytes(b []byte, indent string) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", indent)

	return out.Bytes()
}
