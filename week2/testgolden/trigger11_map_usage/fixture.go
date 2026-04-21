package trigger11

import (
	"math/big"
	"sync"
)

// Rule 29 위반: 전역 장수명 캐시 + 잦은 대량 delete.
// delete로 항목을 지워도 백킹 버킷은 유지 → 메모리 단조 증가.
// 권장: 주기적 재생성(새 맵에 복사 후 교체) 또는 LRU 자료구조 도입.
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
