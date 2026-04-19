# Trigger 06: nil/empty 선언 및 비교

> 이 파일은 gosu-review SKILL의 트리거 06 레퍼런스다.
> 메인 SKILL.md의 워크플로우 3단계에서 이 트리거가 발동하면 읽는다.

## 발동 조건

- `var s []T`, `var m map[K]V` (nil 선언)
- `s := []T{}`, `m := map[K]V{}` (empty 선언)
- `len(x) == 0`, `x == nil` (슬라이스/맵 대상)

## 검사 절차

1. 해당 변수가 직렬화(json/proto), 리플렉션, 외부 API 경계를 타는지 확인
2. nil 검사 vs len 검사의 의도 일치성

## 룰

### Rule 26 — nil vs empty 의미 차이
- 핵심 질문: nil 슬라이스/맵과 빈 슬라이스/맵을 외부 시스템(json, reflect)이 다르게 본다는 점을 의식했는가?
- 예: `json.Marshal(nilSlice)` → `"null"`, `json.Marshal(emptySlice)` → `"[]"`
- 의도가 "값이 없음"이면 nil, "비어있는 컬렉션"이면 empty
- API 경계에서 한 가지로 통일하는 게 안전

### Rule 27 — 컬렉션 비어있음 검사
- 핵심 질문: 슬라이스/맵이 비었는지 볼 때 `s == nil`을 쓰는가?
- 문제: nil 슬라이스와 empty 슬라이스 모두 비어 있지만 `nil` 비교는 둘 중 하나만 잡음
- 권장: `len(s) == 0` (둘 다 잡음, nil 안전)

### Rule 28 — API 모호성 제거
- 핵심 질문: 함수가 nil 인자와 empty 인자를 다르게 취급하도록 설계되었는가?
- 가이드: 다르게 취급할 합리적 이유가 없다면 동일하게 다뤄야 함
- 의도적으로 구분한다면 문서화 필수
