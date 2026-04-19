# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 7건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L20 `func NewParcel(tracking string) Shippable`
**[Rule 7] 인터페이스 반환 지양**
- 근거: `Shippable` 인터페이스를 반환하고 있습니다. 입력은 인터페이스, 출력은 구체 타입을 반환하는 것이 좋습니다.
- 제안:
  ```go
  // before
  func NewParcel(tracking string) Shippable {
  // after
  func NewParcel(tracking string) Parcel {
  ```

#### L24 `func Lookup(key string) any`
**[Rule 8] any 남용**
- 근거: 반환 타입으로 `any`를 사용하고 있습니다. 컴파일 타임의 타입 안정성을 잃을 수 있습니다.
- 제안: 제네릭을 사용하거나 다중 시그니처를 검토하세요.

#### L35 `func FindOrder(id string) error`
**[Rule 44] nil 인터페이스 반환 함정**
- 근거: `error` 인터페이스를 반환할 때 구체 타입의 포인터 `*MyError`를 그대로 반환합니다. `id`가 빈 문자열이 아닐 경우 `e`는 `nil`이 되지만 반환된 `error` 인터페이스는 `nil`이 아니게 되어 호출처에서 에러로 취급됩니다.
- 제안:
  ```go
  // before
  	return e
  // after
  	if e == nil {
  		return nil
  	}
  	return e
  ```

#### L43 `func LoadConfig(path string) ([]byte, error)`
**[Rule 45] io.Reader 활용**
- 근거: 파일 경로(`string`)를 파라미터로 받아 함수 내부에서 `os.Open`을 호출합니다. 테스트와 재사용성을 위해 `io.Reader`를 받는 것이 좋습니다.
- 제안:
  ```go
  // before
  func LoadConfig(path string) ([]byte, error) {
  // after
  func LoadConfig(r io.Reader) ([]byte, error) {
  ```

#### L59 `func NewClient(host string, port int, user, password string, timeoutSec int, retries int) *Client`
**[Rule 11] 함수형 옵션 패턴**
- 근거: 파라미터가 6개로 너무 많습니다. 선택적 파라미터가 있다면 함수형 옵션 패턴(Functional Options) 사용을 고려하세요.
- 제안: 필수 파라미터만 받고 나머지는 `opts ...Option` 형태로 변경합니다.

#### L65 `func (c Client) Address() string`, L66 `func (c *Client) SetHost(h string)`
**[Rule 42] 리시버 타입 결정 (값 vs 포인터)**
- 근거: 같은 타입 `Client`에 대해 값 리시버(`Address`)와 포인터 리시버(`SetHost`)를 혼용하고 있습니다. 일관성을 유지해야 합니다.
- 제안: 모두 포인터 리시버로 통일하세요.
  ```go
  // before
  func (c Client) Address() string
  // after
  func (c *Client) Address() string
  ```

#### L68 `func Split(s string) (first, rest string)`
**[Rule 43] 기명 반환값**
- 근거: 함수 구조가 단순함에도 불필요하게 기명 반환값을 사용하고 있습니다.
- 제안:
  ```go
  // before
  func Split(s string) (first, rest string) {
  // after
  func Split(s string) (string, string) {
  ```
