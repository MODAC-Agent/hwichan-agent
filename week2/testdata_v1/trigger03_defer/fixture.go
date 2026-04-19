package trigger03

import (
	"fmt"
	"os"
)

func DeferArgEval() {
	i := 0
	defer fmt.Println("i was:", i)
	i = 42
	_ = i
}

func ProcessAll(paths []string) error {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		_ = f
	}
	return nil
}

func Divide(a, b int) int {
	if b == 0 {
		panic("b is zero")
	}
	defer handleRecover()
	return a / b
}

func handleRecover() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r)
	}
}
