# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 7건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L6 `out := make([]int, 0)`
**[Rule 21] len/cap 이해**
- 근거: 루프 반복 횟수 `n`을 알 수 있으나 슬라이스 생성 시 용량(capacity)을 지정하지 않아 불필요한 재할당이 발생할 수 있습니다.
- 제안:
  ```go
  // before
  	out := make([]int, 0)
  // after
  	out := make([]int, 0, n)
  ```

#### L15 `b = append(b, 99)`
**[Rule 23] append 충돌 / 슬라이스 축소 누수**
- 근거: `b`는 `a`의 백킹 배열을 공유합니다. `append`를 수행하면 원본 배열 `a`의 데이터(`a[3]`)를 의도치 않게 덮어쓰게 됩니다.
- 제안: `slices.Clone`을 사용하거나 풀 슬라이스 표현식(`a[:3:3]`)을 고려하세요.

#### L21 `copy(dst, src)`
**[Rule 24] copy 동작 인지**
- 근거: `dst`의 `len`이 0이므로 `copy` 함수는 아무것도 복사하지 않습니다.
- 제안:
  ```go
  // before
  	dst := make([]int, 0, 10)
  // after
  	dst := make([]int, len(src))
  ```

#### L30 `q = q[1:]`
**[Rule 25] 포인터 원소 누수**
- 근거: 포인터 슬라이스를 축소할 때 잘려나가는 원소를 `nil`로 처리하지 않으면 GC가 해당 객체를 회수하지 못해 메모리 누수가 발생합니다.
- 제안:
  ```go
  // before
  	q = q[1:]
  // after
  	q[0] = nil
  	q = q[1:]
  ```

#### L36 `s += w + " "`
**[Rule 39] strings.Builder로 반복 결합**
- 근거: `for` 루프 안에서 문자열을 `+=`로 결합하면 매번 메모리를 재할당하게 되어 비용이 매우 큽니다.
- 제안: `strings.Builder`를 사용하거나 `strings.Join`을 사용하세요.

#### L42 `s := string(data)`
**[Rule 40] bytes 패키지 활용**
- 근거: `[]byte`를 `string`으로 변환하여 `strings.Contains`를 사용하고 있습니다. 불필요한 메모리 할당이 발생합니다.
- 제안: `bytes.Contains(data, []byte(tok))`를 사용하세요.

#### L53 `code := blob[10:14]`
**[Rule 41] 서브스트링 누수**
- 근거: 원본 문자열의 일부를 슬라이싱하여 장기 보관(구조체 필드)하면 원본 전체 문자열의 백킹 배열이 GC되지 않습니다.
- 제안: `strings.Clone(blob[10:14])`를 사용하여 복제본을 저장하세요.
