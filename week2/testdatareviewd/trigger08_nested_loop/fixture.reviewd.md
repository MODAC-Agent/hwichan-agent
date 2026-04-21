# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 1건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L14 `break`
**[Rule 34] 루프 탈출 레이블**
- 근거: `for _, e := range events` 안의 `switch e.Kind`에서 `case stopKind:` 분기가 `break`을 사용하지만, 이 `break`은 바깥 `for`가 아닌 `switch`만 탈출한다. 함수 의미상 `stop` 이벤트를 만나면 처리를 중단하려는 의도로 보이는데, 현재 코드에서는 `switch`만 빠져나가고 루프는 계속 돌아 다음 이벤트로 넘어간다. 결과적으로 `stop` 이후 이벤트도 계속 `count++`로 집계되어 의도와 불일치한다.
- 제안:
  ```go
  // before
  func ProcessEvents(events []Event) int {
      count := 0
      for _, e := range events {
          switch e.Kind {
          case stopKind:
              break
          default:
              count++
          }
      }
      return count
  }

  // after
  func ProcessEvents(events []Event) int {
      count := 0
  Loop:
      for _, e := range events {
          switch e.Kind {
          case stopKind:
              break Loop
          default:
              count++
          }
      }
      return count
  }
  ```
  - 만약 "stop 이벤트만 건너뛰고 루프는 계속"하는 의도였다면 `continue`로 바꾸는 편이 명확하다. 현재 `break`은 `switch`의 기본 동작과 동일해 사실상 no-op에 가까워 독자에게 혼동을 준다.
