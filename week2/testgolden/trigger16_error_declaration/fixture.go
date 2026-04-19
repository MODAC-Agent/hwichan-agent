package trigger16

import "errors"

// Rule 50 위반 (1): 호출 측이 분기해야 할 "찾지 못함" 신호를 문자열에 name을 매립해 반환.
// 호출 측은 메시지를 파싱해야만 분기 가능 → 센티널 에러로 값화해야 함.
//
//	권장: var ErrNotFound = errors.New("not found") 를 선언하고 errors.Is로 분기.
func FindByName(name string) error {
	return errors.New("not found: " + name)
}

// Rule 50 위반 (2): "연결 안 됨" 같은 단순 사실 신호라면 센티널이 적합.
// 동반 데이터/복구 컨텍스트가 없는데 구조체 타입으로 만들어 과설계.
//
//	권장: var ErrNotConnected = errors.New("not connected").
type NotConnectedError struct{}

func (e *NotConnectedError) Error() string { return "not connected" }

func Connect() error {
	return &NotConnectedError{}
}
