package trigger06

import "encoding/json"

type Payload struct {
	Tags []string `json:"tags"`
}

// Rule 26 위반: 영 값으로 시작한 nil 슬라이스가 JSON에 "null"로 직렬화됨.
// API 경계에서 "빈 배열"이 기대라면 []string{}로 초기화해야 함.
func DefaultPayload() ([]byte, error) {
	var p Payload
	return json.Marshal(p)
}

// Rule 27 위반: 빈 체크를 == nil로 수행. []string{} 빈 슬라이스는 잡지 못함.
// 권장: len(s) == 0.
func IsEmpty(s []string) bool {
	return s == nil
}

// Rule 28 위반: nil과 empty를 서로 다른 동작으로 취급하지만 문서화/시그니처에 드러나지 않음.
func MergeTags(primary, fallback []string) []string {
	if primary == nil {
		return fallback
	}
	if len(primary) == 0 {
		return primary
	}
	return append(primary, fallback...)
}
