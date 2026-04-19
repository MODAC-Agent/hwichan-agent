package trigger13

import "fmt"

// Rule 9 위반: 타입 파라미터 T가 본문에서 fmt.Stringer의 메서드만 호출.
// 제네릭이 필요한 이득(타입별 보일러플레이트 제거, 타입 안전성)이 없음 → 인터페이스 파라미터로 충분.
//
//	권장: func Describe(x fmt.Stringer) string.
func Describe[T fmt.Stringer](x T) string {
	return "value=" + x.String()
}
