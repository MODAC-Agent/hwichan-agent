package trigger10

import "reflect"

func SameIDs(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}
