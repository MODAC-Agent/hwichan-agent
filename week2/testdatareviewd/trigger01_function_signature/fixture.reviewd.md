# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 7건 (변경 라인 — 파일 전체를 🔴 변경분으로 간주)
- 기존 문제: 0건

---

## fixture.go

### 🔴 변경 라인의 문제

#### L21 `func NewParcel(tracking string) Shippable`
**[Rule 7] 인터페이스 반환 지양**
- 근거: 생성자 함수가 구체 타입 `Parcel`을 생성하면서 반환 타입을 인터페이스 `Shippable`로 선언하고 있다. "입력은 인터페이스, 출력은 구체 타입" 원칙에 어긋나며, `Shippable` 구현체가 현재 `Parcel` 하나뿐인 상황에서 인터페이스 반환은 불필요한 추상화다. 호출처가 `Parcel.Tracking` 필드 등에 접근하려면 타입 단언이 강제된다.
- 제안:
  ```go
  // before
  func NewParcel(tracking string) Shippable {
      return Parcel{Tracking: tracking}
  }

  // after
  func NewParcel(tracking string) Parcel {
      return Parcel{Tracking: tracking}
  }
  // 호출 측에서 Shippable 인터페이스가 필요하면 호출 측의 변수 타입으로 받으면 됨
  ```

---

#### L25 `func Lookup(key string) any`
**[Rule 8] any 남용**
- 근거: 반환 타입이 `any`이나 실제 반환 값은 `int(0)`과 `string("value")`로 한정된다. 임의 타입 수용이 본질인 함수가 아님에도 `any`를 사용해 컴파일 타임 타입 안정성을 상실했다. 호출 측은 반드시 타입 단언을 해야 하며, 잘못된 타입 단언은 런타임 패닉으로 이어진다.
- 제안:
  ```go
  // before
  func Lookup(key string) any {
      if key == "" {
          return 0
      }
      return "value"
  }

  // after (옵션 1): 반환값을 단일 타입으로 통일하고 "없음"을 ok 패턴으로 표현
  func Lookup(key string) (string, bool) {
      if key == "" {
          return "", false
      }
      return "value", true
  }

  // after (옵션 2): 도메인에 따라 제네릭 적용 검토
  ```

---

#### L36 `func FindOrder(id string) error`
**[Rule 44] nil 인터페이스 반환 함정**
- 근거: `var e *MyError`로 nil 구체 포인터를 선언한 뒤, `id`가 비어 있지 않으면 `e`를 그대로 `return e`로 반환한다. `*MyError`는 `error` 인터페이스를 구현하므로, nil 포인터임에도 반환된 `error` 값은 "타입 정보는 있고 값만 nil인" 비-nil 인터페이스가 된다. 호출 측의 `if err != nil` 검사가 `id != ""` 경우에도 **항상 참**이 되어 정상 경로가 오류로 처리된다.
- 제안:
  ```go
  // before
  func FindOrder(id string) error {
      var e *MyError
      if id == "" {
          e = &MyError{Msg: "empty id"}
      }
      return e  // id != "" 일 때 nil 포인터가 비-nil error 인터페이스로 반환됨
  }

  // after
  func FindOrder(id string) error {
      if id == "" {
          return &MyError{Msg: "empty id"}
      }
      return nil  // 명시적 nil 반환으로 인터페이스 함정 회피
  }
  ```

---

#### L44 `func LoadConfig(path string) ([]byte, error)`
**[Rule 45] io.Reader 활용**
- 근거: 파일 경로 문자열(`path string`)을 받아 내부에서 `os.Open`을 직접 호출한다. 파일 시스템에 강하게 결합되어 단위 테스트 시 반드시 실제 파일이 필요하고, 메모리 버퍼(`strings.NewReader`, `bytes.NewReader`)나 네트워크 소스 등 다른 입력 소스로 교체가 불가능하다. 파일 열기와 읽기 책임이 한 함수에 혼재된다.
- 제안:
  ```go
  // before
  func LoadConfig(path string) ([]byte, error) {
      f, err := os.Open(path)
      if err != nil {
          return nil, err
      }
      defer f.Close()
      return io.ReadAll(f)
  }

  // after: io.Reader를 받아 입력 소스에 무관하게 동작
  func LoadConfig(r io.Reader) ([]byte, error) {
      return io.ReadAll(r)
  }

  // 호출 측에서 파일 열기 담당 (자원 관리 책임 분리)
  // f, err := os.Open(path)
  // if err != nil { ... }
  // defer f.Close()
  // cfg, err := LoadConfig(f)

  // 테스트 시:
  // cfg, err := LoadConfig(strings.NewReader(`{"key":"val"}`))
  ```

---

#### L60 `func NewClient(host string, port int, user, password string, timeoutSec int, retries int) *Client`
**[Rule 11] 함수형 옵션 패턴**
- 근거: 파라미터가 6개이며, `password`와 `retries`는 함수 본문에서 `_ = password`, `_ = retries`로 즉시 무시된다. 이는 해당 파라미터들이 실질적으로 미완성 옵션임을 나타낸다. 위치 인자 6개는 호출 시 순서 오류가 컴파일 타임에 잡히지 않으며, 향후 파라미터 추가 시 모든 호출부의 파괴적 변경이 불가피하다.
- 제안:
  ```go
  // before
  func NewClient(host string, port int, user, password string, timeoutSec int, retries int) *Client {
      _ = password
      _ = retries
      return &Client{Host: host, Port: port, User: user, Timeout: timeoutSec}
  }

  // after: 함수형 옵션 패턴
  type ClientOption func(*Client)

  func WithPort(p int) ClientOption      { return func(c *Client) { c.Port = p } }
  func WithUser(u string) ClientOption   { return func(c *Client) { c.User = u } }
  func WithTimeout(s int) ClientOption   { return func(c *Client) { c.Timeout = s } }

  func NewClient(host string, opts ...ClientOption) *Client {
      c := &Client{Host: host}
      for _, o := range opts {
          o(c)
      }
      return c
  }

  // 호출 측:
  // NewClient("localhost", WithPort(8080), WithUser("alice"), WithTimeout(30))
  ```

---

#### L66-L67 `func (c Client) Address()` / `func (c *Client) SetHost(h string)`
**[Rule 42] 리시버 타입 결정 (값 vs 포인터)**
- 근거: 같은 타입 `Client`의 메서드들이 값 리시버(`c Client`)와 포인터 리시버(`c *Client`)를 혼용하고 있다. `SetHost`는 상태 변경을 위해 포인터 리시버가 필수이므로, 일관성 원칙에 따라 `Address`도 포인터 리시버로 통일해야 한다. 혼용하면 인터페이스 구현 시 메서드 집합 불일치로 예상치 못한 컴파일 오류가 발생할 수 있다.
- 제안:
  ```go
  // before
  func (c Client) Address() string   { return c.Host }
  func (c *Client) SetHost(h string) { c.Host = h }

  // after: 포인터 리시버로 통일
  func (c *Client) Address() string  { return c.Host }
  func (c *Client) SetHost(h string) { c.Host = h }
  ```

---

#### L69 `func Split(s string) (first, rest string)`
**[Rule 43] 기명 반환값**
- 근거: 함수 본문이 단순하고 반환값의 의미가 직관적이나 기명 반환값 + bare `return`을 사용하고 있다. `defer`에서 반환값을 수정하거나 같은 타입의 반환값이 3개 이상인 복잡한 경우가 아니다. `len(s) == 0`일 때 bare `return`으로 zero-value(`"", ""`)를 반환하는 의도가 명시적이지 않아 독자가 zero-value 반환임을 추론해야 한다.
- 제안:
  ```go
  // before
  func Split(s string) (first, rest string) {
      if len(s) == 0 {
          return
      }
      first, rest = s[:1], s[1:]
      return
  }

  // after: 명시적 반환으로 의도 명확화
  func Split(s string) (string, string) {
      if len(s) == 0 {
          return "", ""
      }
      return s[:1], s[1:]
  }
  ```
