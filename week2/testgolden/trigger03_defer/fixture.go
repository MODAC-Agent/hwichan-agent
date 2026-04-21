package trigger03

import (
	"fmt"
	"os"
)

// Rule 46 위반: defer 인자 i는 등록 시점(0)에 평가됨.
// 이후 i = 42 대입은 반영 안 됨.
func DeferArgEval() {
	i := 0
	defer fmt.Println("i was:", i)
	i = 42
	_ = i
}

// Rule 35 위반: for 루프 안 defer. 함수 종료까지 fd가 누적돼 파일 디스크립터 고갈 위험.
// Rule 53 위반: defer f.Close()의 에러를 무시. 쓰기 파일이면 버퍼 플러시 실패를 놓침.
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

// Rule 47 위반 (1): 일반 에러 흐름에 panic을 제어 흐름으로 사용.
// Rule 47 위반 (2): recover()가 defer 함수 직접 안이 아닌 다른 함수 호출 경유 → 케이스에 따라 동작 안 함.
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
