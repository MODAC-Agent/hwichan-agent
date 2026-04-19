# Trigger 08: 중첩 루프/switch/select

> 이 파일은 gosu-review SKILL의 트리거 08 레퍼런스다.
> 메인 SKILL.md의 워크플로우 3단계에서 이 트리거가 발동하면 읽는다.

## 발동 조건

`for` 안에 다른 `for`/`switch`/`select`가 있고 그 안에서 `break`/`continue`를 사용

## 검사 절차

1. 중첩된 제어 구조 내부의 `break`/`continue`를 찾는다
2. 각 키워드가 어느 문장을 탈출하려는 의도인지 판단
3. `switch`/`select` 내부의 `break`가 바깥 `for`를 탈출하려는 의도면 레이블 필요

## 룰

### Rule 34 — 루프 탈출 레이블
- 핵심 질문: 의도한 바깥 문장을 빠져나가기 위해 레이블을 사용했는가?
- 문제:
  ```go
  for _, x := range items {
      switch x.Kind {
      case Stop:
          break  // switch만 탈출, for는 계속
      }
  }
  ```
- 수정:
  ```go
  Loop:
  for _, x := range items {
      switch x.Kind {
      case Stop:
          break Loop
      }
  }
  ```
- `select` 안의 `break`도 동일 함정
