# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 (`fixture.go`)
- 새로 발견된 문제: 7건 (변경 라인 — 전체 🔴 처리)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L6 `out := make([]int, 0)`
**[Rule 21] len/cap 이해**
- 근거: `SquaresUpTo(n int)`는 최종 크기가 `n`으로 이미 알려져 있는데 용량 힌트 없이 `make([]int, 0)`으로 시작해 `n`번 `append`한다. 매 반복마다 재할당·복사가 발생할 수 있다.
- 제안:
  ```go
  // before
  out := make([]int, 0)
  for i := 0; i < n; i++ {
      out = append(out, i*i)
  }

  // after — 직접 인덱스 대입 (재할당 0회)
  out := make([]int, n)
  for i := 0; i < n; i++ {
      out[i] = i * i
  }

  // 또는 용량 힌트만 주기
  out := make([]int, 0, n)
  for i := 0; i < n; i++ {
      out = append(out, i*i)
  }
  ```

#### L13-18 `AppendConflict`
**[Rule 23] append 충돌 / 슬라이스 공유 백킹 배열**
- 근거: `b := a[:3]`는 `a`와 동일한 백킹 배열을 공유한다. `b`의 cap은 5이므로 `b = append(b, 99)`는 새 배열을 할당하지 않고 `a[3]`을 `99`로 덮어쓴다. 결과적으로 반환된 `a`는 `{1, 2, 3, 99, 5}`가 되어 호출자가 예상하지 못한 변형이 발생한다.
- 제안:
  ```go
  // before
  b := a[:3]
  b = append(b, 99)

  // after — 풀 슬라이스 표현식으로 cap 제한 (append 시 새 배열 할당 강제)
  b := a[:3:3]
  b = append(b, 99)

  // 또는 명시적 복제
  b := slices.Clone(a[:3])
  b = append(b, 99)
  ```

#### L20-24 `CopyShort`
**[Rule 24] copy 동작 인지**
- 근거: `dst := make([]int, 0, 10)`은 `len(dst)==0`이다. `copy(dst, src)`는 `min(len(dst), len(src)) == 0`개를 복사하므로 항상 빈 슬라이스를 반환한다. 함수 이름 및 의도와 실제 동작이 일치하지 않는다.
- 제안:
  ```go
  // before
  dst := make([]int, 0, 10)
  copy(dst, src)
  return dst

  // after — len을 src와 같게 설정
  dst := make([]int, len(src))
  copy(dst, src)
  return dst

  // 또는 append 활용
  return append([]int(nil), src...)
  ```

#### L30-33 `PopFront`
**[Rule 25] 포인터 원소 누수**
- 근거: `q []*Item` 슬라이스(선언: L30)를 `q = q[1:]`로 앞에서 축소하면서 잘려나간 `q[0]`을 nil 처리하지 않는다. 백킹 배열의 0번째 슬롯이 여전히 `*Item`을 참조하므로 해당 객체는 GC되지 않고 메모리에 남는다.
- 제안:
  ```go
  // before
  func PopFront(q []*Item) []*Item {
      q = q[1:]
      return q
  }

  // after — nil 처리 후 축소
  func PopFront(q []*Item) []*Item {
      if len(q) == 0 {
          return q
      }
      q[0] = nil // 백킹 배열의 포인터 제거 → GC 가능
      q = q[1:]
      return q
  }

  // 또는 slices.Delete 활용 (내부에서 zero 처리)
  // return slices.Delete(q, 0, 1)
  ```

#### L35-41 `JoinWords`
**[Rule 39] strings.Builder로 반복 결합**
- 근거: `for ... range words` 루프 안에서 `s += w + " "`로 문자열을 누적한다. Go 문자열은 불변이므로 매 반복마다 새 메모리 할당과 전체 복사가 발생해 N개 단어 기준 O(N²) 비용이 된다.
- 제안:
  ```go
  // before
  var s string
  for _, w := range words {
      s += w + " "
  }
  return s

  // after — strings.Builder 사용
  var b strings.Builder
  for _, w := range words {
      b.WriteString(w)
      b.WriteByte(' ')
  }
  return b.String()

  // 또는 구분자 결합이 목적이면 (끝 공백 동작 재검토 필요)
  // return strings.Join(words, " ")
  ```

#### L43-49 `ContainsToken`
**[Rule 40] bytes 패키지 활용**
- 근거: `data []byte`를 받아 `s := string(data)`(1회 할당)로 변환한 뒤 `strings.Contains`를 호출하고, 다시 `string([]byte(s))`(2회 할당)로 왕복 변환해 두 번째 `Contains`를 호출한다. 두 번의 불필요한 메모리 할당이 발생하며 두 번째 `Contains`는 동일 입력에 대한 중복 호출이다. `bytes` 패키지 함수를 사용하면 변환 자체가 불필요하다.
- 제안:
  ```go
  // before
  func ContainsToken(data []byte, tok string) bool {
      s := string(data)
      if !strings.Contains(s, tok) {
          return false
      }
      return strings.Contains(string([]byte(s)), tok)
  }

  // after — bytes.Contains로 변환 제거
  func ContainsToken(data []byte, tok string) bool {
      return bytes.Contains(data, []byte(tok))
  }
  // tok이 고정 상수라면 호출처에서 []byte(tok)를 미리 준비해 재사용
  ```

#### L55-58 `Extract`
**[Rule 41] 서브스트링 누수**
- 근거: `blob string`(선언: L55 파라미터)에서 `blob[10:14]`로 슬라이싱한 4바이트 부분 문자열을 `Record.Code` 필드에 대입해 함수 외부로 반환한다. Go 문자열 슬라이싱은 원본의 백킹 배열을 공유하므로 `Record`가 살아 있는 동안 `blob` 전체가 GC되지 않는다. `blob`이 큰 문자열일수록 누수 규모가 커진다.
- 제안:
  ```go
  // before
  func Extract(blob string) *Record {
      code := blob[10:14]
      return &Record{Code: code}
  }

  // after — strings.Clone으로 새 백킹 배열 확보
  func Extract(blob string) *Record {
      code := strings.Clone(blob[10:14])
      return &Record{Code: code}
  }
  ```
