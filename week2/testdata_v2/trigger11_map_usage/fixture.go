package trigger11

import (
	"math/big"
	"sync"
)

var (
	cache   = map[string]*big.Int{}
	cacheMu sync.Mutex
)

func Put(k string, v *big.Int) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache[k] = v
}

func Evict(keys []string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	for _, k := range keys {
		delete(cache, k)
	}
}
