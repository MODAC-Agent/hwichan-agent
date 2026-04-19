package trigger10

import "reflect"

// Rule 30 위반: 기본 타입 슬라이스 비교에 reflect.DeepEqual.
// slices.Equal이 더 빠르고 의도가 명확함.
func SameIDs(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}
