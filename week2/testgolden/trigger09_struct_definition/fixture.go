package trigger09

import "sync"

// Rule 4 위반: 검증/계산 없이 단순 필드 노출에 Get/Set 메서드를 기계적으로 붙임.
// Go 관례: 필드를 직접 노출하거나 Getter는 `Name()` (Get 접두사 없이).
type Account struct {
	name string
}

func (a *Account) GetName() string  { return a.name }
func (a *Account) SetName(n string) { a.name = n }

// Rule 10 위반: sync.Mutex 임베딩으로 Lock/Unlock이 Cache의 외부 API로 promotion.
// 내부 구현 디테일은 명명된 필드로 숨겨야 함: `mu sync.Mutex`.
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
