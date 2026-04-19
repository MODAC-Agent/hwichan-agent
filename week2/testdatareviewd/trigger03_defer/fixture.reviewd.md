# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 5건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L10 `defer fmt.Println("i was:", i)`
**[Rule 46] defer 인자 평가 시점**
- 근거: `defer`에 넘긴 인자 `i`는 **등록 시점**에 즉시 평가되어 값 `0`이 캡처된다. 이후 L11에서 `i = 42`로 바뀌어도 출력은 `i was: 0`이 된다. 함수 종료 시점의 값을 출력하려는 의도라면 버그다.
- 제안:
  ```go
  // before
  i := 0
  defer fmt.Println("i was:", i)
  i = 42

  // after — 클로저로 감싸 호출 시점에 i를 읽도록
  i := 0
  defer func() { fmt.Println("i was:", i) }()
  i = 42
  ```

#### L16–L23 `for _, p := range paths { ... defer f.Close() ... }`
**[Rule 35] 루프 내 defer 주의 (함수 추출)**
- 근거: `defer f.Close()`가 `for` 루프 안에 있어 `ProcessAll`이 끝날 때까지 모든 파일 핸들이 누적된다. `paths`가 커지면 fd 고갈로 이어진다. 각 반복이 끝나는 즉시 닫혀야 한다.
- 제안:
  ```go
  // before
  func ProcessAll(paths []string) error {
      for _, p := range paths {
          f, err := os.Open(p)
          if err != nil { return err }
          defer f.Close()
          _ = f
      }
      return nil
  }

  // after — 반복 로직을 함수로 추출해 defer 스코프를 한 반복으로 좁힌다
  func ProcessAll(paths []string) error {
      for _, p := range paths {
          if err := processOne(p); err != nil {
              return err
          }
      }
      return nil
  }

  func processOne(path string) (err error) {
      f, err := os.Open(path)
      if err != nil { return err }
      defer func() {
          if cerr := f.Close(); cerr != nil && err == nil {
              err = cerr
          }
      }()
      _ = f
      return nil
  }
  ```

#### L21 `defer f.Close()`
**[Rule 53] defer 에러 처리**
- 근거: `f.Close()`가 반환하는 에러가 완전히 버려진다. 읽기 파일이라도 명시적으로 무시 의도를 표현하지 않아 의도가 불명확하다. 쓰기 파일이라면 버퍼 플러시 실패를 놓칠 수 있어 특히 위험하다.
- 제안:
  ```go
  // before
  defer f.Close()

  // after (읽기 전용 — 명시적 무시)
  defer func() { _ = f.Close() }()

  // after (기명 반환값으로 에러 상위 전달 — Rule 35 개선안과 함께 적용)
  defer func() {
      if cerr := f.Close(); cerr != nil && err == nil {
          err = cerr
      }
  }()
  ```

#### L29 `panic("b is zero")`
**[Rule 47] recover / panic 사용 규칙 — panic을 일반 에러 흐름에 사용**
- 근거: `b == 0`은 호출자가 입력값으로 충분히 예측·방어할 수 있는 일반 에러 흐름이다. `panic`은 "복구 불가능한 프로그래머 실수/불변조건 위반"에만 써야 하며, 일반 입력 검증은 `error` 반환값으로 처리해야 한다.
- 제안:
  ```go
  // before
  func Divide(a, b int) int {
      if b == 0 {
          panic("b is zero")
      }
      defer handleRecover()
      return a / b
  }

  // after — 에러를 반환값으로 표현
  func Divide(a, b int) (int, error) {
      if b == 0 {
          return 0, fmt.Errorf("divide: b is zero")
      }
      return a / b, nil
  }
  ```

#### L31, L35–L38 `defer handleRecover()` / `func handleRecover() { recover() ... }`
**[Rule 47] recover / panic 사용 규칙 — defer 간접 호출 안에서 recover()**
- 근거: `recover()`는 `defer`로 등록된 함수 리터럴의 **바로 그 프레임**에서 호출돼야만 유효하다. `defer handleRecover()`는 `handleRecover`를 defer 직접 호출하지만, recover는 그 함수 내부(한 단계 아래)에서 불린다. 실제로는 동작할 수 있으나, 이 구조는 공식 안티패턴으로 간주된다. 래핑 레이어가 추가되면 `recover()`는 `nil`을 반환하여 panic을 잡지 못하고 프로세스가 종료된다.
- 제안:
  ```go
  // before
  defer handleRecover()

  func handleRecover() {
      if r := recover(); r != nil {
          fmt.Println("recovered:", r)
      }
  }

  // after — recover는 defer에 등록된 함수 리터럴 바로 안에서 호출
  defer func() {
      if r := recover(); r != nil {
          fmt.Println("recovered:", r)
      }
  }()
  ```
- 참고: Rule 47에 따라 `b == 0` panic 자체를 `error` 반환으로 대체하는 것이 우선이다. `recover`는 HTTP 핸들러 등 라이브러리 경계에서 예상치 못한 panic을 로깅·변환할 때 사용한다.

---
