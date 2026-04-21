# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 3건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L9 `var p Payload`
**[Rule 26] nil vs empty 의미 차이**
- 근거: `Payload` 내부의 `Tags` 슬라이스가 `nil`로 초기화됩니다. 이를 JSON 직렬화하면 `"tags": null`이 되어 클라이언트가 빈 배열 `"tags": []`를 예상할 경우 문제가 될 수 있습니다.
- 제안: 명시적으로 빈 슬라이스로 초기화하거나 JSON 태그에 `omitempty`를 고려하세요.

#### L14 `return s == nil`
**[Rule 27] 컬렉션 비어있음 검사**
- 근거: 슬라이스가 비었는지 확인할 때 `nil` 검사만 수행하면 할당되었지만 비어 있는(`[]string{}`) 슬라이스는 놓치게 됩니다.
- 제안:
  ```go
  // before
  	return s == nil
  // after
  	return len(s) == 0
  ```

#### L18 `if primary == nil {`
**[Rule 28] API 모호성 제거**
- 근거: `primary`가 `nil`일 때와 단순히 비어있을 때(`len(primary) == 0`)의 동작이 다릅니다. 이 차이가 의도된 것인지 검토가 필요합니다.
- 제안: 합리적인 이유가 없다면 동일하게 처리되도록 통합하세요.
