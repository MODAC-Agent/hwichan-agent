package trigger15

import "strconv"

// Rule 18 위반: 외부 입력을 검증 없이 int32로 좁힘.
// 32비트 범위를 벗어나면 조용히 wrap around.
func ParseAge(raw string) int32 {
	n, _ := strconv.Atoi(raw)
	return int32(n)
}

// Rule 19 위반: float64 금액 비교에 ==.
// IEEE 754 오차로 연산 경로에 따라 결과가 달라짐. 엡실론 범위 비교 필요.
func IsExpectedPrice(price, expected float64) bool {
	return price == expected
}

// Rule 20 위반: 크기 격차가 큰 항(bigVal vs tinies)을 그대로 누적.
// 자릿수 손실이 누적됨 — tinies를 먼저 합산한 뒤 bigVal과 더하는 편이 안정적.
func AccumulateSmall(bigVal float64, tinies []float64) float64 {
	s := bigVal
	for _, t := range tinies {
		s = s + t
	}
	return s
}
