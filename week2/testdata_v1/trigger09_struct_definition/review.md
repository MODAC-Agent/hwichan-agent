# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 2건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L9 `func (a *Account) GetName() string`, L10 `func (a *Account) SetName(n string)`
**[Rule 4] 게터/세터 기계적 적용**
- 근거: 단순 접근을 위한 게터와 세터를 기계적으로 만들었습니다. 검증 등의 로직이 없다면 필드를 직접 노출(`Name string`)하는 것이 Go 관례에 맞습니다.
- 제안: 필드를 외부에 노출(`Name`)하거나, 메서드 이름을 `Name()`, `SetName()`으로 변경하세요.

#### L13 `sync.Mutex` (임베딩)
**[Rule 10] 임베딩 가시성**
- 근거: `Cache` 구조체에 `sync.Mutex`를 임베딩하여 `Lock()`과 `Unlock()` 메서드가 외부에 노출됩니다. 이는 캡슐화를 위반하고 외부에서 잘못 사용할 위험이 있습니다.
- 제안: 명명된 내부 필드로 변경하세요.
  ```go
  // before
  type Cache struct {
  	sync.Mutex
  // after
  type Cache struct {
  	mu sync.Mutex
  ```
