# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 4건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L9 `defer fmt.Println("i was:", i)`
**[Rule 46] defer 인자 평가 시점**
- 근거: `defer`로 넘긴 인자 `i`는 `defer` 등록 시점에 평가되므로 0이 출력됩니다. 이후 변경된 값을 출력하려면 클로저를 사용해야 합니다.
- 제안:
  ```go
  // before
  	defer fmt.Println("i was:", i)
  // after
  	defer func() { fmt.Println("i was:", i) }()
  ```

#### L20 `defer f.Close()`
**[Rule 35] 루프 내 defer 주의 (함수 추출)**
- 근거: `for` 루프 안에서 `defer`를 호출하면 함수 종료 시점까지 자원 해제가 미뤄져 파일 디스크립터가 고갈될 수 있습니다.
- 제안: 루프 내부 로직을 별도 함수로 추출하여 `defer` 스코프를 좁히세요.

#### L28 `panic("b is zero")`
**[Rule 47] recover / panic 사용 규칙**
- 근거: 0으로 나누기와 같은 에러 상황은 일반적인 에러 흐름(`error` 반환)으로 처리해야 합니다. 제어 흐름에 `panic`을 사용하는 것은 권장되지 않습니다.
- 제안: `error` 반환값으로 변경하세요.

#### L30 `defer handleRecover()`
**[Rule 47] recover / panic 사용 규칙**
- 근거: `recover()`는 `defer`로 직접 호출된 함수 내부에서만 유효합니다. `handleRecover`처럼 한 단계 더 들어간 곳에서 `recover()`를 호출하면 `nil`을 반환하여 동작하지 않습니다.
- 제안:
  ```go
  // before
  	defer handleRecover()
  // after
  	defer func() {
  		if r := recover(); r != nil {
  			fmt.Println("recovered:", r)
  		}
  	}()
  ```
