# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 1건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L5 `func Describe[T fmt.Stringer](x T) string`
**[Rule 9] 제네릭 남용**
- 근거: 제네릭 타입 `T`를 사용하여 `fmt.Stringer` 제약을 주었으나, 타입 파라미터가 반환 타입이나 다른 인자와의 의존성이 없습니다. 단순 인터페이스 파라미터로 충분히 표현 가능하므로 제네릭을 사용하는 것은 불필요한 복잡성입니다.
- 제안: 일반 인터페이스 파라미터로 변경하세요.
  ```go
  // before
  func Describe[T fmt.Stringer](x T) string {
  // after
  func Describe(x fmt.Stringer) string {
  ```
