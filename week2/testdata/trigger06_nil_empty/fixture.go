package trigger06

import "encoding/json"

type Payload struct {
	Tags []string `json:"tags"`
}

func DefaultPayload() ([]byte, error) {
	var p Payload
	return json.Marshal(p)
}

func IsEmpty(s []string) bool {
	return s == nil
}

func MergeTags(primary, fallback []string) []string {
	if primary == nil {
		return fallback
	}
	if len(primary) == 0 {
		return primary
	}
	return append(primary, fallback...)
}
