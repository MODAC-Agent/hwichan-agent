package trigger05

import (
	"strings"
)

// Rule 21 위반: n이 예상 크기인데 make([]int, 0)로 시작 → 반복 재할당.
// 권장: make([]int, 0, n).
func SquaresUpTo(n int) []int {
	out := make([]int, 0)
	for i := 0; i < n; i++ {
		out = append(out, i*i)
	}
	return out
}

// Rule 23 위반: b는 a의 백킹 배열을 공유. b에 append하면 a[3]을 덮어씀.
func AppendConflict() ([]int, []int) {
	a := []int{1, 2, 3, 4, 5}
	b := a[:3]
	b = append(b, 99)
	return a, b
}

// Rule 24 위반: dst의 len이 0이라 copy가 0개만 복사.
// 권장: make([]int, len(src)) 또는 append(dst, src...).
func CopyShort(src []int) []int {
	dst := make([]int, 0, 10)
	copy(dst, src)
	return dst
}

type Item struct {
	ID string
}

// Rule 25 위반: []*Item 큐에서 앞을 제거하면서 잘려나간 원소를 nil 처리하지 않음.
// 백킹 배열이 여전히 포인터를 참조 → GC 회수 불가.
func PopFront(q []*Item) []*Item {
	q = q[1:]
	return q
}

// Rule 39 위반: for 루프 안에서 문자열 += 결합. O(N²) 비용.
// 권장: strings.Builder 또는 strings.Join.
func JoinWords(words []string) string {
	var s string
	for _, w := range words {
		s += w + " "
	}
	return s
}

// Rule 40 위반: []byte → string → []byte 왕복 변환. bytes 패키지로 대체 가능.
func ContainsToken(data []byte, tok string) bool {
	s := string(data)
	if !strings.Contains(s, tok) {
		return false
	}
	return strings.Contains(string([]byte(s)), tok)
}

type Record struct {
	Code string
}

// Rule 41 위반: blob(큰 문자열)의 서브슬라이스를 구조체 필드에 장기 보관.
// blob 전체의 백킹 배열이 GC되지 않음. strings.Clone 필요.
func Extract(blob string) *Record {
	code := blob[10:14]
	return &Record{Code: code}
}
