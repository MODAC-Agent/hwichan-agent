package trigger13

import "fmt"

func Describe[T fmt.Stringer](x T) string {
	return "value=" + x.String()
}
