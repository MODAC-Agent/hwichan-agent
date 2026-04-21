# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 2건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L6 `return errors.New("not found: " + name)`
**[Rule 50] 센티널 에러 vs 에러 타입**
- 근거: 에러 메시지에 컨텍스트 데이터(`name`)를 결합해 반환하고 있습니다. 호출 측에서 특정 이름에 대한 에러인지 프로그램적으로 판단(분기)하기 어렵게 만듭니다. (호출 측 확인 필요)
- 제안: 호출 측에서 에러 세부 정보가 필요하다면 커스텀 에러 구조체 타입을 정의하거나, 센티널 에러를 포맷팅 래핑(`%w`)하여 반환하세요.

#### L9 `type NotConnectedError struct{}`
**[Rule 50] 센티널 에러 vs 에러 타입**
- 근거: 추가적인 컨텍스트 데이터를 담지 않는 단순한 상태를 나타내는 에러임에도 불구하고 커스텀 구조체 타입으로 정의되었습니다. 단순 센티널 에러 값으로 충분합니다.
- 제안:
  ```go
  // before
  type NotConnectedError struct{}
  func (e *NotConnectedError) Error() string { return "not connected" }
  // after
  var ErrNotConnected = errors.New("not connected")
  ```
