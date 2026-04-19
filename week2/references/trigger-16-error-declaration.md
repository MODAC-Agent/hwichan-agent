# Trigger 16: 에러 타입/값 선언

> 이 파일은 gosu-review SKILL의 트리거 16 레퍼런스다.
> 메인 SKILL.md의 워크플로우 3단계에서 이 트리거가 발동하면 읽는다.

## 발동 조건

- **센티널 에러 선언**: `var Err... = errors.New(...)` / `var Err... = fmt.Errorf(...)` (전역/패키지 레벨)
- **커스텀 에러 타입**: `type ...Error struct { ... }` 또는 `Error() string` 메서드를 가진 타입

## 검사 절차

1. 에러가 **예상 가능한(expected) 상황**을 표현하는지, **예상 못한(unexpected) 구체 정보**를 담아야 하는지 판단
2. 각 경우에 맞는 표현 방식(값 vs 타입)이 쓰였는가?
3. 사용처에서 `errors.Is`(값) / `errors.As`(타입)으로 검사하고 있는가? (트리거 2의 Rule 49와 연계)

## 룰

### Rule 50 — 센티널 에러 vs 에러 타입
- 핵심 질문: 예상 가능한 에러는 **값**(센티널)으로, 예상 못한 / 구체 정보가 필요한 에러는 **타입**으로 표현하고 있는가?
- 원칙:
  - **값 (센티널 에러)** — 단순 "이 조건이 발생했다"는 사실만 전달하면 충분한 경우
    ```go
    var ErrNotFound = errors.New("not found")
    var ErrUnauthorized = errors.New("unauthorized")
    ```
    호출 측: `if errors.Is(err, ErrNotFound) { ... }`
  - **타입 (커스텀 에러)** — 에러에 동반 데이터가 있거나, 복구 로직이 데이터를 필요로 하는 경우
    ```go
    type ValidationError struct {
        Field  string
        Reason string
    }
    func (e *ValidationError) Error() string { ... }
    ```
    호출 측: `var ve *ValidationError; if errors.As(err, &ve) { log(ve.Field) }`
- 안티패턴:
  - 센티널 에러만으로 충분한데 구조체 타입을 만들어 복잡도 증가
  - 반대로 컨텍스트 데이터가 필요한데 `errors.New("... " + varName)` 식으로 문자열에 묻어버림 (호출 측이 파싱 불가)
- 네이밍 규칙:
  - 센티널: `ErrXxx` (변수) — `Err` 접두사
  - 타입: `XxxError` (타입) — `Error` 접미사
- 공개 범위: 외부에서 분기할 에러만 export, 내부 구현 디테일은 소문자 시작

## 출력 시 주의사항

이 트리거는 단일 파일만 보고는 판단하기 어려울 수 있다. 특히 "예상 가능한지" 여부는 **사용처 패턴**을 봐야 함. 리포트의 근거 란에 "이 에러가 호출 측에서 분기되는지 확인 필요" 같은 단서를 달아 거짓 양성을 줄일 것.

Rule 49(트리거 2)와 짝을 이룬다 — 에러를 **선언**하는 쪽(Rule 50)과 **소비**하는 쪽(Rule 49)은 함께 본다.
