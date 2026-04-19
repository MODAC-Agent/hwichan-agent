# Trigger 03: defer (`defer ...`)

> 이 파일은 gosu-review SKILL의 트리거 03 레퍼런스다.
> 메인 SKILL.md의 워크플로우 3단계에서 이 트리거가 발동하면 읽는다.

## 발동 조건

- `defer` 키워드가 있는 모든 호출
- `recover()` 호출 (Rule 47)
- `panic(...)` 호출 (Rule 47)

## 검사 절차

1. `defer` 호출의 인자/리시버를 확인 — 등록 시점 이후에 변경되는 변수가 들어가면 Rule 46
2. `defer`가 **`for` 루프 안**에 있으면 Rule 35 (스택 쌓임)
3. `defer`에서 호출하는 함수가 에러를 반환하는데 무시하고 있으면 Rule 53
4. `recover()`가 `defer` 함수 밖에서 호출되면 Rule 47
5. `panic(...)`이 "프로그래머 실수" 수준이 아닌 일반 에러 흐름에 쓰이면 Rule 47

## 룰

### Rule 46 — defer 인자 평가 시점
- 핵심 질문: `defer`로 넘기는 인자/리시버가 호출 시점이 아닌 **defer 등록 시점**에 평가되는 것을 의식했는가?
- 문제 패턴:
  ```go
  i := 0
  defer fmt.Println(i)  // 0이 출력됨
  i = 42
  ```
- 등록 시점 즉시 평가됨. 나중 값을 캡처하려면:
  - 클로저로 감싸기: `defer func() { fmt.Println(i) }()`
  - 포인터 전달: `defer fmt.Println(&i)` (단, fmt가 포인터 그대로 출력하므로 위 클로저가 보통 정답)
- 메서드 리시버도 동일 — 값 리시버 메서드를 defer하면 등록 시점 값으로 호출됨

### Rule 35 — 루프 내 defer 주의 (함수 추출)
- 핵심 질문: `for` 루프 안에서 `defer`를 사용해 함수 종료까지 자원 해제가 미뤄지고 있지 않은가?
- 문제: `defer`는 **함수 종료 시점**에 실행 — 루프 반복마다 defer가 쌓임
  ```go
  for _, path := range paths {
      f, err := os.Open(path)
      if err != nil { return err }
      defer f.Close()   // 루프 끝날 때까지 안 닫힘 — 수천 개면 fd 고갈
      // ... f 사용
  }
  ```
- 수정: 반복 로직을 함수로 추출해 defer 스코프를 한 반복으로 좁힌다
  ```go
  for _, path := range paths {
      if err := processOne(path); err != nil { return err }
  }

  func processOne(path string) error {
      f, err := os.Open(path)
      if err != nil { return err }
      defer f.Close()
      // ...
      return nil
  }
  ```

### Rule 47 — recover / panic 사용 규칙
- 핵심 질문: `recover()`가 `defer` 함수 안에서 호출되는가? `panic`이 특수한 상황에만 쓰이는가?
- `recover()` 규칙:
  - `defer` 함수 **바로 안**에서만 유효 — 그 외 위치에서 호출하면 `nil` 반환, 효과 없음
  - 안티패턴: `defer`에서 다른 함수를 호출하고 그 함수 안에서 `recover()` → 일부 경우 제대로 동작 안 함. `defer func() { if r := recover(); r != nil { ... } }()` 형태로 직접 작성
- `panic` 규칙:
  - "정말 복구 불가능한 프로그래머 실수"에만 사용: 의존성 초기화 실패, 불변조건 위반
  - 일반 에러 흐름은 `error` 반환값으로 처리 — `panic`을 제어 흐름으로 쓰지 말 것
- 검사: 라이브러리 경계(특히 HTTP 핸들러)에서 `recover()` 없이 고루틴이 panic하면 프로세스 종료

### Rule 53 — defer 에러 처리
- 핵심 질문: `defer`로 호출한 함수가 에러를 반환하는데 무시하고 있지 않은가?
- 안티패턴: `defer f.Close()` — `Close`의 에러는 버려짐. 특히 **쓰기 파일**에서는 버퍼 플러시 실패를 놓칠 수 있음
- 옵션:
  - 클로저로 래핑해 로깅: `defer func() { if err := f.Close(); err != nil { log.Warn(...) } }()`
  - 기명 반환값으로 상위에 전달:
    ```go
    func process() (err error) {
        defer func() {
            if cerr := f.Close(); cerr != nil && err == nil {
                err = cerr
            }
        }()
        ...
    }
    ```
- 주의: 본 함수가 반환하는 `err`과 `defer`에서 발생한 에러가 **둘 다 있는 경우** — 둘을 병합한 커스텀 에러를 만들거나, 우선순위를 정해 하나만 반환하되 나머지는 로깅
- 읽기 전용 자원(읽기 파일 닫기 등)은 `_ = f.Close()`로 명시적 무시해도 무방