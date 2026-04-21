# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개 / 새 위반: 3건 / 기존 위반: 0건

---

## fixture.go

### 🔴 변경 라인
- **L8** `DeferArgEval` — [Rule 46] defer 인자 평가 시점
  - 근거: `defer fmt.Println(..., i)`는 호출 시점의 `i` 값(0)을 평가하여 나중 값이 출력되지 않습니다.
  - 제안: 클로저 `defer func() { fmt.Println("i was:", i) }()`를 사용해 지연 평가되게 합니다.
- **L19** `ProcessAll` — [Rule 35] 루프 내 defer 주의 (함수 추출)
  - 근거: `for` 루프 안에서 `defer f.Close()`를 호출해 루프가 끝날 때까지 파일 디스크립터가 해제되지 않습니다.
  - 제안: 반복 몸체를 별도의 함수로 추출하여 반복마다 `defer`가 실행되게 합니다.
- **L27** `Divide` — [Rule 47] recover / panic 사용 규칙
  - 근거: `defer handleRecover()` 안의 다른 함수에서 `recover()`를 호출하여 복구가 제대로 동작하지 않습니다.
  - 제안: `defer func() { if r := recover(); r != nil { ... } }()` 형태로 직접 작성합니다.
