# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 1건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L6 `return reflect.DeepEqual(a, b)`
**[Rule 30] 동등성 비교 선택**
- 근거: `SameIDs`의 인자 타입은 `[]int`로, 슬라이스 원소가 비교 가능 타입(int)만 포함한다. `reflect.DeepEqual`은 리플렉션 기반으로 동작해 불필요한 성능 오버헤드가 발생한다. 또한 `reflect.DeepEqual`은 nil 슬라이스와 빈 슬라이스를 다르게 취급하는 미묘한 차이가 있어, 호출부의 의도에 따라 다른 결과를 낼 수 있다. 요소가 단순 comparable 타입인 `[]int` 비교에는 `slices.Equal`이 더 빠르고 의도가 명확하다.
- 제안:
  ```go
  // before
  import "reflect"

  func SameIDs(a, b []int) bool {
      return reflect.DeepEqual(a, b)
  }

  // after
  import "slices"

  func SameIDs(a, b []int) bool {
      return slices.Equal(a, b)
  }
  ```

### 🟡 기존 라인의 문제 (참고)

없음.

---
