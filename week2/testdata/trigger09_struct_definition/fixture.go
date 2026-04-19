package trigger09

import "sync"

type Account struct {
	name string
}

func (a *Account) GetName() string  { return a.name }
func (a *Account) SetName(n string) { a.name = n }

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
