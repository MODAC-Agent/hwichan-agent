# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 2건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L5** `FindByName` — [Rule 50] 센티널 에러 vs 에러 타입
  - 근거: 에러의 원인 데이터(`name`)가 동반된 에러인데 단순 문자열로 생성되어 호출 측이 에러의 원인 데이터를 구조적으로 파악할 수 없습니다.
  - 제안: 예상치 못한 정보가 필요하다면 구조체 타입 커스텀 에러(`type NotFoundError struct { Name string }`)로 정의하세요.
- **L8** `NotConnectedError` — [Rule 50] 센티널 에러 vs 에러 타입
  - 근거: 부가적인 데이터 없이 단순 "연결되지 않음"이라는 사실만 전달하므로 구조체보다는 센티널 에러 변수가 더 적합합니다.
  - 제안: `var ErrNotConnected = errors.New("not connected")` 형식의 센티널 에러 값으로 변경하세요.
