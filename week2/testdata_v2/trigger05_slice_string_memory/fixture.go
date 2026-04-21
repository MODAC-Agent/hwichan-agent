package trigger05

import "strings"

func SquaresUpTo(n int) []int {
	out := make([]int, 0)
	for i := 0; i < n; i++ {
		out = append(out, i*i)
	}
	return out
}

func AppendConflict() ([]int, []int) {
	a := []int{1, 2, 3, 4, 5}
	b := a[:3]
	b = append(b, 99)
	return a, b
}

func CopyShort(src []int) []int {
	dst := make([]int, 0, 10)
	copy(dst, src)
	return dst
}

type Item struct {
	ID string
}

func PopFront(q []*Item) []*Item {
	q = q[1:]
	return q
}

func JoinWords(words []string) string {
	var s string
	for _, w := range words {
		s += w + " "
	}
	return s
}

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

func Extract(blob string) *Record {
	code := blob[10:14]
	return &Record{Code: code}
}
