package trigger07

import "fmt"

type User struct {
	Name string
}

// Rule 31 위반: range 값 u는 원소의 복제본. 필드 수정이 원본에 반영되지 않음.
// 권장: for i := range users { users[i].Name = ... }.
func RenameAll(users []User) {
	for _, u := range users {
		u.Name = "X"
		_ = u
	}
}

// Rule 32 위반: 배열(슬라이스 아님) range → 배열 전체 값 복제 발생.
// 권장: `for i := range &arr` 또는 슬라이스로 변환.
func SumBigArray() int {
	var arr [1 << 16]int
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}

// Rule 33 위반: map 순회 결과가 "특정 순서로 온다"고 가정한 첫 키 반환.
// Go map은 순회 순서를 보장하지 않음 — 호출마다 달라질 수 있음.
func FirstKey(m map[string]int) string {
	for k := range m {
		return k
	}
	return ""
}

// Rule 36 위반: 문자열을 바이트 단위로 인덱싱. 다국어 입력에서 룬을 깨뜨림.
func FirstByte(s string) byte {
	if len(s) == 0 {
		return 0
	}
	return s[0]
}

// Rule 37 위반: `for i, r := range s`의 i는 룬의 시작 바이트 인덱스인데
// "몇 번째 문자인지" 의미로 프린트해 혼동을 유발.
func PrintWithIndex(s string) {
	for i, r := range s {
		fmt.Printf("char #%d: %c\n", i, r)
	}
}
