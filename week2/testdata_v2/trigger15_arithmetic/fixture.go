package trigger15

import "strconv"

func ParseAge(raw string) int32 {
	n, _ := strconv.Atoi(raw)
	return int32(n)
}

func IsExpectedPrice(price, expected float64) bool {
	return price == expected
}

func AccumulateSmall(bigVal float64, tinies []float64) float64 {
	s := bigVal
	for _, t := range tinies {
		s = s + t
	}
	return s
}
