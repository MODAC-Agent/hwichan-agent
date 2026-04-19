# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 2건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L3 `type Notifier interface`
**[Rule 5] 인터페이스 선제 생성 금지**
- 근거: 이 파일에서 확인 가능한 구현체는 `*EmailNotifier` 단 1개이며, 테스트에서의 모킹이나 명시적 다형성 요구도 확인되지 않는다. "필요할지도 모른다"는 이유로 선제적으로 도입된 추상화일 가능성이 높다. (이 파일만으로는 프로젝트 전체 구현체 수를 확인할 수 없으므로 프로젝트 차원 검토 권장. 구현체가 다수이거나 외부 경계 — DB, HTTP 클라이언트 등 — 대체가 필요한 경우는 예외)
- 제안:
  ```go
  // before
  type Notifier interface {
      Notify(msg string) error
  }

  type EmailNotifier struct {
      From string
  }

  func (e *EmailNotifier) Notify(msg string) error {
      _ = msg
      return nil
  }

  // after — 구현체가 1개뿐이고 테스트 모킹 수요가 없다면 인터페이스 제거
  type EmailNotifier struct {
      From string
  }

  func (e *EmailNotifier) Notify(msg string) error {
      _ = msg
      return nil
  }
  // 이후 다형성이나 테스트 더블이 실제로 필요해지는 시점에
  // 클라이언트 패키지에서 필요한 메서드만 담은 인터페이스를 정의한다.
  ```

#### L3 `type Notifier interface`
**[Rule 6] 인터페이스 위치 (클라이언트 측 정의)**
- 근거: `Notifier` 인터페이스가 구현체 `EmailNotifier`와 동일한 패키지(`trigger12`)에 함께 선언되어 있다. Go 관례상 인터페이스는 이를 소비하는 클라이언트 패키지에 두어 필요한 메서드만 좁게 정의하는 편이 결합도를 낮추고 작은 인터페이스를 자연스럽게 유도한다. (이 파일만으로는 호출 패키지를 확인할 수 없으므로 프로젝트 차원 검토 권장. `io.Reader`처럼 매우 광범위하게 공유되는 표준 패턴이라면 예외)
- 제안:
  ```go
  // before — 구현체 패키지(trigger12)에서 인터페이스도 함께 export
  package trigger12

  type Notifier interface {
      Notify(msg string) error
  }

  type EmailNotifier struct{ From string }

  func (e *EmailNotifier) Notify(msg string) error { return nil }

  // after — 구현체 패키지는 구체 타입만 제공하고,
  //         인터페이스는 이를 사용하는 클라이언트 패키지에서 정의
  package trigger12

  type EmailNotifier struct{ From string }

  func (e *EmailNotifier) Notify(msg string) error { return nil }

  // 클라이언트 패키지 예시
  package alerting

  type notifier interface {
      Notify(msg string) error
  }

  func Send(n notifier, msg string) error {
      return n.Notify(msg)
  }
  ```

---
