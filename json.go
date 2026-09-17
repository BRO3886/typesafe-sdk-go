package typesafe

import (
	"bytes"
	"encoding/json"
)

func encodeJSON(payload any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(payload); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

// decodeBody parses a body as JSON, falling back to its text: servers and
// proxies do not always set a content type, and error pages are often HTML.
func decodeBody(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err == nil {
		return decoded
	}
	return string(raw)
}
