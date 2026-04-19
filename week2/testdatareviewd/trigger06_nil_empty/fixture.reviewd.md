# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 3건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L9-12 `DefaultPayload()`
**[Rule 26] nil vs empty 의미 차이**
- 근거: `var p Payload`로 선언하면 `Payload.Tags` 필드가 nil 슬라이스가 된다. 이 상태로 `json.Marshal(p)`를 호출하면 `{"tags":null}`이 직렬화 결과에 포함된다. API 소비자가 `"tags":[]`(빈 배열)를 기대한다면 의도와 다른 응답이 되며, API 경계에서 nil/empty 중 하나로 통일하는 것이 안전하다.
- 제안:
  ```go
  // before
  func DefaultPayload() ([]byte, error) {
      var p Payload
      return json.Marshal(p)
  }

  // after — "비어있는 컬렉션"이 의도라면 empty 슬라이스로 초기화
  func DefaultPayload() ([]byte, error) {
      p := Payload{Tags: []string{}}
      return json.Marshal(p)
  }
  ```

---

#### L14-16 `IsEmpty()`
**[Rule 27] 컬렉션 비어있음 검사**
- 근거: `return s == nil`은 nil 슬라이스만 true로 판단한다. `s := []string{}`처럼 명시적으로 초기화된 empty 슬라이스는 `s == nil`이 false가 되므로, IsEmpty가 false를 반환한다. 함수 이름 `IsEmpty`("비어있는가")의 의미와 구현이 불일치한다.
- 제안:
  ```go
  // before
  func IsEmpty(s []string) bool {
      return s == nil
  }

  // after — nil과 empty 슬라이스 모두 "비어있음"으로 처리
  func IsEmpty(s []string) bool {
      return len(s) == 0
  }
  ```

---

#### L18-26 `MergeTags()`
**[Rule 28] API 모호성 제거**
- 근거: `primary == nil`이면 fallback을 반환하고(L19-21), `len(primary) == 0`이면 primary(empty 슬라이스)를 그대로 반환한다(L22-24). nil과 empty를 의도적으로 다르게 취급하고 있지만, 그 근거와 의미가 주석·문서로 명시되어 있지 않다. 호출자는 nil과 `[]string{}`을 구분해서 넘겨야 하는지 알 수 없어 API 계약이 모호하다.
- 제안:
  ```go
  // before
  func MergeTags(primary, fallback []string) []string {
      if primary == nil {
          return fallback
      }
      if len(primary) == 0 {
          return primary
      }
      return append(primary, fallback...)
  }

  // after (안 1) — nil/empty를 동일하게 다루는 경우
  func MergeTags(primary, fallback []string) []string {
      if len(primary) == 0 {
          return fallback
      }
      return append(primary, fallback...)
  }

  // after (안 2) — nil/empty를 의도적으로 구분해야 하는 경우, 문서 명시
  // MergeTags merges primary and fallback tag slices.
  // nil primary: "unset"으로 간주하여 fallback을 반환.
  // empty (non-nil) primary: "명시적으로 비어있음"으로 간주하여 병합하지 않고 반환.
  func MergeTags(primary, fallback []string) []string {
      if primary == nil {
          return fallback
      }
      if len(primary) == 0 {
          return primary
      }
      return append(primary, fallback...)
  }
  ```

---
