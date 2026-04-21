# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 1건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L9 `cache = map[string]*big.Int{}` (+ L19-25 `Evict`)

**[Rule 29] 맵 메모리 단조 증가**

- 근거: `cache`는 패키지 레벨 전역 변수(장수명 맵)이며, 이름 자체에 `cache`가 포함되어 트리거 조건에 해당한다. `Evict` 함수는 `keys []string` 슬라이스를 받아 루프를 돌며 `delete(cache, k)`를 수행하는 **대량 삭제 패턴**이 명시적으로 존재한다. Go의 맵은 `delete`로 항목을 지워도 내부 해시 버킷(백킹 스토리지)은 축소되지 않으므로, 캐시가 한 번 커진 뒤 대량 삭제되어도 프로세스 메모리는 줄지 않고 단조 증가한다. **장수명 + 대량 삭제** 두 조건이 모두 충족되어 실질적 메모리 누수 위험이 있다. 값 타입이 `*big.Int` 포인터라는 점은 값 자체 메모리는 GC 가능하게 해 주지만, 맵 버킷의 키(`string`) 및 버킷 슬롯은 그대로 유지된다.
- 제안:
  ```go
  // before
  var (
      cache   = map[string]*big.Int{}
      cacheMu sync.Mutex
  )

  func Evict(keys []string) {
      cacheMu.Lock()
      defer cacheMu.Unlock()
      for _, k := range keys {
          delete(cache, k)
      }
  }

  // after — 옵션 A: 대량 삭제 시 살아있는 항목만 복사하여 맵 재생성
  func Evict(keys []string) {
      cacheMu.Lock()
      defer cacheMu.Unlock()

      toDelete := make(map[string]struct{}, len(keys))
      for _, k := range keys {
          toDelete[k] = struct{}{}
      }

      // 삭제 비율이 일정 임계치 이상이면 재생성으로 버킷까지 반납
      if len(toDelete) >= len(cache)/4 {
          fresh := make(map[string]*big.Int, len(cache)-len(toDelete))
          for k, v := range cache {
              if _, drop := toDelete[k]; drop {
                  continue
              }
              fresh[k] = v
          }
          cache = fresh
          return
      }
      for k := range toDelete {
          delete(cache, k)
      }
  }

  // after — 옵션 B: 용도가 LRU/TTL 캐시라면 전용 라이브러리로 교체하여 크기 상한 강제
  //   (hashicorp/golang-lru, ristretto 등)
  ```

  검토 권장: 캐시 최대 크기, 삭제 빈도/비율, 호출 패턴(배치 무효화인지 개별 삭제인지)을 실측한 뒤 옵션 A(주기적 재생성) 또는 옵션 B(상한 있는 캐시 자료구조) 중 선택. `sync.Mutex` + 단순 `map` 조합을 유지하더라도, 재생성 전략을 병행하지 않으면 장기 운용 시 RSS가 peak에 고정된다는 점을 문서화 또는 테스트로 남길 것.

---
