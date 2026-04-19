package trigger07

import "fmt"

type User struct {
	Name string
}

func RenameAll(users []User) {
	for _, u := range users {
		u.Name = "X"
		_ = u
	}
}

func SumBigArray() int {
	var arr [1 << 16]int
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}

func FirstKey(m map[string]int) string {
	for k := range m {
		return k
	}
	return ""
}

func FirstByte(s string) byte {
	if len(s) == 0 {
		return 0
	}
	return s[0]
}

func PrintWithIndex(s string) {
	for i, r := range s {
		fmt.Printf("char #%d: %c\n", i, r)
	}
}
