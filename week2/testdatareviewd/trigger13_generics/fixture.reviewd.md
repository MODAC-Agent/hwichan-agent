# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 1건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L5 `func Describe[T fmt.Stringer](x T) string`
**[Rule 9] 제네릭 남용**
- 근거: 타입 파라미터가 `T`이고 제약이 `fmt.Stringer` 인터페이스 하나뿐이다. 함수 본문에서는 `x.String()` 메서드 호출만 수행하므로, `fmt.Stringer`를 직접 매개변수 타입으로 사용하는 것과 동일한 의미다. 레퍼런스 위반 패턴("타입 파라미터가 1개뿐이고 본문에서 그 타입의 메서드만 호출 → 인터페이스가 더 단순")에 정확히 해당하며, 보일러플레이트 제거나 타입 안전성 이득이 없다.
- 제안:
  ```go
  // before
  func Describe[T fmt.Stringer](x T) string {
      return "value=" + x.String()
  }

  // after
  func Describe(x fmt.Stringer) string {
      return "value=" + x.String()
  }
  ```
