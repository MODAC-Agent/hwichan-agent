# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 7건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L6** `SquaresUpTo` — [Rule 21] len/cap 이해
  - 근거: 반복문을 돌며 `append`하지만 예상 길이를 아는데도 용량 힌트 없이 슬라이스를 생성했습니다.
  - 제안: `out := make([]int, 0, n)`으로 초기 용량을 지정하여 재할당을 방지합니다.
- **L16** `AppendConflict` — [Rule 23] append 충돌 / 슬라이스 축소 누수
  - 근거: 동일한 백킹 배열을 공유하는 원본 슬라이스의 축소판에 `append`를 수행하여 원본 데이터를 덮어씁니다.
  - 제안: `b := slices.Clone(a[:3])` 등으로 백킹 배열을 분리합니다.
- **L22** `CopyShort` — [Rule 24] copy 동작 인지
  - 근거: `dst`의 길이가 0이므로 `copy(dst, src)` 호출 시 아무 데이터도 복사되지 않습니다.
  - 제안: `dst := make([]int, len(src))`로 목적지 슬라이스의 길이를 지정합니다.
- **L30** `PopFront` — [Rule 25] 포인터 원소 누수
  - 근거: 포인터 슬라이스 `q`를 축소할 때 잘려나간 `q[0]`을 `nil`로 만들지 않아 메모리 누수가 발생합니다.
  - 제안: `q[0] = nil` 처리 후 슬라이스를 축소합니다.
- **L36** `JoinWords` — [Rule 39] strings.Builder로 반복 결합
  - 근거: `for` 루프 내부에서 `+` 연산자로 문자열을 반복 결합하여 매번 새 메모리 할당이 발생합니다.
  - 제안: `strings.Builder`를 사용하거나 `strings.Join`을 활용합니다.
- **L43** `ContainsToken` — [Rule 40] bytes 패키지 활용
  - 근거: `[]byte`를 불필요하게 `string`으로 변환하여 처리해 메모리 할당 비용을 발생시킵니다.
  - 제안: 문자열 변환 없이 `bytes.Contains(data, []byte(tok))`를 사용합니다.
- **L53** `Extract` — [Rule 41] 서브스트링 누수
  - 근거: 큰 문자열에서 작은 부분만 슬라이싱한 결과를 구조체 필드에 대입해 장기 보관하여 원본 배열 전체 누수를 일으킬 수 있습니다.
  - 제안: `strings.Clone(blob[10:14])`를 사용하여 메모리를 복제합니다.
