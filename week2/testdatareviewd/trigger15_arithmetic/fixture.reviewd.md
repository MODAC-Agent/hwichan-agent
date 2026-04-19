# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 3건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L7 `return int32(n)`
**[Rule 18] 정수 오버/언더플로**
- 근거: `strconv.Atoi`는 플랫폼 기본 `int` 크기(64비트 환경에서 최대 `math.MaxInt64`)를 반환한다. 이를 아무런 범위 검증 없이 `int32`로 좁히면, 입력값이 `math.MaxInt32`(2,147,483,647)를 초과하거나 `math.MinInt32`(-2,147,483,648) 미만일 때 조용히 wrap around된다. `raw`는 외부 입력 문자열이므로 실제 위험 경로다.
- 제안:
  ```go
  // before
  func ParseAge(raw string) int32 {
      n, _ := strconv.Atoi(raw)
      return int32(n)
  }

  // after
  func ParseAge(raw string) (int32, error) {
      n, err := strconv.Atoi(raw)
      if err != nil {
          return 0, err
      }
      if n < math.MinInt32 || n > math.MaxInt32 {
          return 0, fmt.Errorf("age out of int32 range: %d", n)
      }
      return int32(n), nil
  }
  ```

---

#### L11 `return price == expected`
**[Rule 19] 부동소수점 직접 비교 지양**
- 근거: `price`와 `expected` 모두 `float64`이며, `==`로 직접 비교하고 있다. IEEE 754 표현상 부동소수점 연산 경로에 따라 미세 오차가 발생할 수 있어 같은 논리적 값이라도 `false`를 반환할 수 있다. 특히 함수 이름이 `IsExpectedPrice`로 가격 비교에 사용되는 맥락이라면 금융 계산의 위험도가 더 높다.
- 제안:
  ```go
  // before
  func IsExpectedPrice(price, expected float64) bool {
      return price == expected
  }

  // after (엡실론 절대 오차)
  import "math"

  const eps = 1e-9

  func IsExpectedPrice(price, expected float64) bool {
      return math.Abs(price-expected) <= eps
  }

  // after (상대 오차 — 금액 크기가 다양할 때 권장)
  func IsExpectedPrice(price, expected float64) bool {
      if expected == 0 {
          return math.Abs(price) <= 1e-9
      }
      return math.Abs(price-expected)/math.Abs(expected) <= 1e-9
  }

  // after (금융 정밀 계산 — shopspring/decimal 사용)
  // import "github.com/shopspring/decimal"
  // func IsExpectedPrice(price, expected decimal.Decimal) bool {
  //     return price.Equal(expected)
  // }
  ```

---

#### L14–L19 `AccumulateSmall` 함수
**[Rule 20] 부동소수점 연산 순서**
- 근거: `bigVal`이라는 이름에서 알 수 있듯 첫 번째 인자는 큰 값이고, `tinies`는 작은 값들의 슬라이스다. 현재 구현은 `s = bigVal`로 시작해 작은 값들을 하나씩 더하고 있다. 큰 수와 작은 수를 직접 더할 때 자릿수 손실(catastrophic cancellation)이 발생할 수 있다. 즉, 작은 값들의 누적합이 `bigVal`에 비해 무시될 정도로 작으면 결과의 정밀도가 낮아진다.
- 제안:
  ```go
  // before
  func AccumulateSmall(bigVal float64, tinies []float64) float64 {
      s := bigVal
      for _, t := range tinies {
          s = s + t
      }
      return s
  }

  // after — 작은 값들을 먼저 합산한 뒤 bigVal에 더하기
  func AccumulateSmall(bigVal float64, tinies []float64) float64 {
      var sum float64
      for _, t := range tinies {
          sum += t
      }
      return bigVal + sum
  }

  // 더 높은 정밀도가 필요한 경우 — Kahan 보상 합산
  func AccumulateSmall(bigVal float64, tinies []float64) float64 {
      var sum, c float64
      for _, t := range tinies {
          y := t - c
          tmp := sum + y
          c = (tmp - sum) - y
          sum = tmp
      }
      return bigVal + sum
  }
  ```
