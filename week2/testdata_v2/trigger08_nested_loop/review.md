# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 1건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L8** `ProcessEvents` — [Rule 34] 루프 탈출 레이블
  - 근거: `switch` 내부의 `break`는 바깥쪽 `for` 루프가 아닌 `switch` 문만 탈출합니다.
  - 제안: 루프 전체를 탈출하려면 `Loop:` 레이블을 선언하고 `break Loop`를 사용해야 합니다.
