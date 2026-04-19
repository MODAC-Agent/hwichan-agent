# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 (`fixture.go`)
- 새로 발견된 문제: 5건 (변경 라인 — 전체를 🔴로 간주)
- 기존 문제: 0건

---

## fixture.go

### 🔴 변경 라인의 문제

#### L10–13 `for _, u := range users`

**[Rule 31] range 값 원소는 복제본**
- 근거: `for _, u := range users`에서 `u`는 `User` 구조체의 복제본이다. `u.Name = "X"` 대입은 복제본의 필드만 수정하며, 원본 슬라이스 `users`의 원소에는 전혀 반영되지 않는다. 함수명 `RenameAll`이 "전체 이름 변경"을 의도하지만 실제로는 아무 효과가 없는 코드다.
- 제안:
  ```go
  // before
  func RenameAll(users []User) {
      for _, u := range users {
          u.Name = "X"
          _ = u
      }
  }

  // after (인덱스 접근 — 원본 직접 수정)
  func RenameAll(users []User) {
      for i := range users {
          users[i].Name = "X"
      }
  }

  // after (포인터 슬라이스로 시그니처 변경)
  func RenameAll(users []*User) {
      for _, u := range users {
          u.Name = "X"
      }
  }
  ```

---

#### L19–20 `for _, v := range arr`

**[Rule 32] range 표현식 단일 평가 — 큰 배열 값 복제**
- 근거: `arr`은 슬라이스가 아닌 `[1 << 16]int` 배열(65,536개 원소, 약 512KB)이다. `for _, v := range arr`는 루프 시작 시 배열 전체를 한 번 복제한다. 값을 읽기만 할 뿐 수정하지 않으므로 불필요한 대용량 복사가 발생한다.
- 제안:
  ```go
  // before
  var arr [1 << 16]int
  for _, v := range arr {
      s += v
  }

  // after (포인터로 range — 복제 없음)
  var arr [1 << 16]int
  for _, v := range &arr {
      s += v
  }

  // after (슬라이스로 변환)
  for _, v := range arr[:] {
      s += v
  }
  ```

---

#### L26–29 `for k := range m`

**[Rule 33] Map 순회 특성 — 순서 비결정적**
- 근거: `for k := range m { return k }`는 map에서 "첫 번째"로 나오는 키를 반환하지만, Go의 map 순회 순서는 실행마다 무작위로 달라진다. 함수명 `FirstKey`는 결정적인 "첫 번째 키"를 암시하므로, 호출자가 일관된 결과를 기대할 경우 오동작 버그로 이어진다.
- 제안:
  ```go
  // before
  func FirstKey(m map[string]int) string {
      for k := range m {
          return k
      }
      return ""
  }

  // after (의도가 "사전순 첫 키"라면 — 정렬 후 반환)
  func FirstKey(m map[string]int) string {
      keys := make([]string, 0, len(m))
      for k := range m {
          keys = append(keys, k)
      }
      sort.Strings(keys)
      if len(keys) == 0 {
          return ""
      }
      return keys[0]
  }

  // after (의도가 "임의의 키 하나"라면 — 이름과 주석 명확화)
  // AnyKey returns an arbitrary key from m, or "" if m is empty.
  func AnyKey(m map[string]int) string {
      for k := range m {
          return k
      }
      return ""
  }
  ```

---

#### L35 `return s[0]`

**[Rule 36] 룬과 바이트 — 바이트 인덱싱으로 다국어 첫 문자 깨짐 위험**
- 근거: `s[0]`은 문자열의 첫 번째 바이트를 반환한다. ASCII 문자열에서는 문제없으나, 한글·이모지 등 멀티바이트 UTF-8 문자를 포함한 문자열에서는 첫 룬의 일부 바이트만 반환하여 의미 없는 값이 된다. 함수명 `FirstByte`가 바이트 반환을 명시하지만, 호출자가 "첫 문자"로 오해할 위험이 있다.
- 제안:
  ```go
  // before
  func FirstByte(s string) byte {
      if len(s) == 0 {
          return 0
      }
      return s[0]
  }

  // after (ASCII 전용임을 문서화)
  // FirstByte returns the first byte of s. Only safe for ASCII strings.
  func FirstByte(s string) byte {
      if len(s) == 0 {
          return 0
      }
      return s[0]
  }

  // after (첫 룬 반환이 의도라면)
  import "unicode/utf8"

  func FirstRune(s string) rune {
      if len(s) == 0 {
          return 0
      }
      r, _ := utf8.DecodeRuneInString(s)
      return r
  }
  ```

---

#### L41 `fmt.Printf("char #%d: %c\n", i, r)`

**[Rule 37] 문자열 range — `i`는 룬 순번이 아닌 바이트 오프셋**
- 근거: `for i, r := range s`에서 `i`는 룬 `r`의 시작 바이트 오프셋이지, 0·1·2·3 형태의 룬 순번이 아니다. 다국어 문자열(예: 한글)에서는 `i`가 0, 3, 6, 9...처럼 건너뛴다. `"char #%d"`가 연속된 문자 번호처럼 보이므로, ASCII 문자열에서만 우연히 올바른 출력이 나오고 다국어 문자열에서는 오해를 유발한다.
- 제안:
  ```go
  // before
  func PrintWithIndex(s string) {
      for i, r := range s {
          fmt.Printf("char #%d: %c\n", i, r)
      }
  }

  // after (룬 순번이 필요한 경우 — 별도 카운터 사용)
  func PrintWithIndex(s string) {
      idx := 0
      for _, r := range s {
          fmt.Printf("char #%d: %c\n", idx, r)
          idx++
      }
  }

  // after (바이트 오프셋임을 명시)
  func PrintWithIndex(s string) {
      for byteOffset, r := range s {
          fmt.Printf("byte offset %d: %c\n", byteOffset, r)
      }
  }
  ```

---
