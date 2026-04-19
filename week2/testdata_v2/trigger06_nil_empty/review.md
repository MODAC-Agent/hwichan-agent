# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 3건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L8** `DefaultPayload` — [Rule 26] nil vs empty 의미 차이
  - 근거: `Payload.Tags`가 초기화되지 않은 nil 슬라이스인 상태로 `json.Marshal`되어 `"null"`로 직렬화됩니다.
  - 제안: 빈 배열 `"[]"` 출력을 의도했다면 `var p = Payload{Tags: []string{}}` (empty 선언)으로 변경합니다.
- **L12** `IsEmpty` — [Rule 27] 컬렉션 비어있음 검사
  - 근거: 슬라이스가 비어있는지 확인할 때 `s == nil`을 사용하여 빈 슬라이스(`[]string{}`)를 처리하지 못합니다.
  - 제안: `len(s) == 0`으로 변경하여 nil과 empty 슬라이스를 모두 검사하도록 수정합니다.
- **L16** `MergeTags` — [Rule 28] API 모호성 제거
  - 근거: `primary == nil`과 `len(primary) == 0` 분기를 따로 두어 nil과 empty 슬라이스를 다르게 취급하고 있습니다.
  - 제안: 의도적인 구분이 아니라면 `len(primary) == 0` 하나로 통합하여 일관되게 처리합니다.
