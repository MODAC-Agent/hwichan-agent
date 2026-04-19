# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 2건 (변경 라인 — 전체 🔴 간주)
- 기존 문제: 0건

---

## fixture.go

### 🔴 변경 라인의 문제

#### L5-L10 `type Account struct { name string }` + `GetName`/`SetName`
**[Rule 4] 게터/세터 기계적 적용**
- 근거: `Account.name` 필드는 단순 값 보관용이며, 검증/계산/동시성 가드/인터페이스 충족 같은 정당한 사유가 없다. 이런 상태에서 `GetName()`/`SetName()`을 일률적으로 만드는 것은 Go 관례에 어긋난다. 특히 `GetX` 접두사는 비관용적이며, 접근자가 꼭 필요해도 `Name()` 형태여야 한다.
- 제안:
  ```go
  // before
  type Account struct {
      name string
  }
  func (a *Account) GetName() string  { return a.name }
  func (a *Account) SetName(n string) { a.name = n }

  // after (옵션 1: 단순 노출이 목적이라면 필드를 직접 export)
  type Account struct {
      Name string
  }

  // after (옵션 2: 접근자가 꼭 필요한 경우 — Get 접두사 제거)
  type Account struct {
      name string
  }
  func (a *Account) Name() string     { return a.name }
  func (a *Account) SetName(n string) { a.name = n }
  ```

#### L12-L24 `type Cache struct { sync.Mutex; ... }`
**[Rule 10] 임베딩 가시성**
- 근거: `sync.Mutex`를 익명 임베딩하면 `Lock()`/`Unlock()` 메서드가 `Cache`로 promotion되어 외부 패키지에서 `cache.Lock()`을 직접 호출할 수 있다. 뮤텍스는 `Cache`의 내부 동기화 수단일 뿐 외부 API가 아니므로, 이 임베딩은 캡슐화를 깨고 오용(외부 코드가 락을 잡은 채 다른 메서드를 호출하는 등)을 유발할 위험이 있다. 명명된 필드로 바꿔 내부 구현 세부로 숨겨야 한다.
- 제안:
  ```go
  // before
  type Cache struct {
      sync.Mutex
      data map[string]string
  }
  func (c *Cache) Put(k, v string) {
      c.Lock()
      defer c.Unlock()
      if c.data == nil {
          c.data = map[string]string{}
      }
      c.data[k] = v
  }

  // after
  type Cache struct {
      mu   sync.Mutex
      data map[string]string
  }
  func (c *Cache) Put(k, v string) {
      c.mu.Lock()
      defer c.mu.Unlock()
      if c.data == nil {
          c.data = map[string]string{}
      }
      c.data[k] = v
  }
  ```

---
