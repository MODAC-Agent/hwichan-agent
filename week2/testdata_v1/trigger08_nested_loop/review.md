# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 1건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L12 `break`
**[Rule 34] 루프 탈출 레이블**
- 근거: `for` 루프 안의 `switch` 문에서 `break`를 호출하고 있습니다. 이는 `switch` 문만 탈출할 뿐 바깥쪽 `for` 루프는 계속 실행됩니다.
- 제안: 의도한 것이 `for` 루프 탈출이라면 레이블을 사용하세요.
  ```go
  // before
  		case stopKind:
  			break
  // after
  	Loop:
  		for _, e := range events {
  			switch e.Kind {
  			case stopKind:
  				break Loop
  ...
  ```
