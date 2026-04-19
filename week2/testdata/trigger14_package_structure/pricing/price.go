// Rule 12 대조 예시: 도메인 단위 이름.
// 위의 util/common과 달리 `pricing`은 책임이 명확해 파일 배치가 자명해진다.
package pricing

type Price struct {
	Cents int64
}

func Sum(a, b Price) Price {
	return Price{Cents: a.Cents + b.Cents}
}
