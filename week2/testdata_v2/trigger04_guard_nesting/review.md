# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 1건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L14** `Authorize` — [Rule 2] 중첩 최소화 (happy path 좌측 정렬)
  - 근거: 조건문이 4중으로 중첩되어 가독성이 떨어지며 정상 흐름이 숨겨져 있습니다.
  - 제안: `if user == nil { return ErrForbidden }` 등 가드 절을 사용해 조기 반환(early return) 형태로 리팩토링합니다.
