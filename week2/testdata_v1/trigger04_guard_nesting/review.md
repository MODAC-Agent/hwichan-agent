# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 1건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L15 `if user != nil {` (이하 중첩)
**[Rule 2] 중첩 최소화 (happy path 좌측 정렬)**
- 근거: `if`문이 3단 이상 깊게 중첩되어 있어 가독성이 떨어집니다. 가드 절(Guard Clause)을 사용하여 예외 상황을 일찍 반환하고 정상 흐름을 밖으로 빼는 것이 좋습니다.
- 제안:
  ```go
  // before
  	if user != nil {
  		if user.Active { ... }
  // after
  	if user == nil || !user.Active {
  		return ErrForbidden
  	}
  	if resource == nil {
  		return ErrForbidden
  	}
  	if !resource.Public && !user.Admin {
  		return ErrForbidden
  	}
  	if action != "read" && !user.Admin {
  		return ErrForbidden
  	}
  	return nil
  ```
