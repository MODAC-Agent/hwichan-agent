# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 2건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L6 `errors.New("not found: " + name)`
**[Rule 50] 센티널 에러 vs 에러 타입**
- 근거: `FindByName`은 실패 시 호출자에게 `name` 값이 필요한 상황임에도, 동적 데이터를 `errors.New("not found: " + name)` 식으로 문자열에 묻어버리고 있다. 호출 측에서 `name`을 꺼내 분기 처리하거나 로깅할 수 없으며, 매 호출마다 새로운 에러 인스턴스를 생성하므로 `errors.Is` 비교도 불가능하다. 컨텍스트 데이터가 필요한 에러는 커스텀 타입으로 표현해야 한다.
- 제안:
  ```go
  // before
  func FindByName(name string) error {
      return errors.New("not found: " + name)
  }

  // after — 컨텍스트 데이터가 필요하다면 커스텀 타입 사용
  type NotFoundError struct {
      Name string
  }

  func (e *NotFoundError) Error() string {
      return "not found: " + e.Name
  }

  func FindByName(name string) error {
      return &NotFoundError{Name: name}
  }
  // 호출 측: var nfe *NotFoundError; if errors.As(err, &nfe) { log(nfe.Name) }

  // 또는, 호출 측에서 name 값이 필요 없고 "못 찾았다"는 사실만 분기하면 되는 경우
  var ErrNotFound = errors.New("not found")

  func FindByName(name string) error {
      return ErrNotFound
  }
  // 호출 측: if errors.Is(err, ErrNotFound) { ... }
  ```

#### L9-11 `type NotConnectedError struct{}`
**[Rule 50] 센티널 에러 vs 에러 타입**
- 근거: `NotConnectedError`는 필드가 전혀 없는 빈 구조체(`struct{}`)이고, `Error()` 메서드도 정적 문자열만 반환한다. 에러에 동반 데이터가 없고 복구 로직이 데이터를 필요로 하지 않으므로, "연결되지 않았다"는 사실 자체만 전달하는 센티널 에러로 충분하다. 불필요한 구조체 타입 정의는 복잡도를 높이고, `errors.As` 대신 `errors.Is`를 쓸 수 있는 상황에서 호출 측 코드를 번거롭게 만든다.
- 제안:
  ```go
  // before
  type NotConnectedError struct{}

  func (e *NotConnectedError) Error() string { return "not connected" }

  func Connect() error {
      return &NotConnectedError{}
  }

  // after — 구체 데이터가 없으므로 센티널 에러로 단순화
  var ErrNotConnected = errors.New("not connected")

  func Connect() error {
      return ErrNotConnected
  }
  // 호출 측: if errors.Is(err, ErrNotConnected) { ... }
  ```

---
