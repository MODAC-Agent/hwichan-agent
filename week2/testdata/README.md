# gosu-review 테스트 픽스처

`week2/` gosu-review 스킬의 동작 검증용 Go 픽스처 모음. 각 하위 디렉토리는 `week2/Skill.md`의 트리거 번호와 1:1 매핑된다.

## 파일 구성

픽스처는 블라인드 입력과 정답지를 물리적으로 분리한다.

- `week2/testdata/triggerNN_*/fixture.go` — **답안 없음**. 스킬이 힌트 없이도 같은 룰을 검출하는지 검증하는 주 입력.
- `week2/testgolden/triggerNN_*/fixture.go` — **답안 포함**. `// Rule N 위반: ...` 주석으로 심어둔 위반을 표시. 리뷰어가 기대 항목을 눈으로 확인할 때 참고.

트리거 14는 디렉토리 구조 자체가 입력이므로 `testdata/trigger14_package_structure/{util,common,pricing}/`(블라인드)과 `testgolden/trigger14_package_structure/{util,common,pricing}/`(정답지) 두 갈래로 제공한다.

## 사용법

1. 먼저 `testdata/triggerNN_*/fixture.go`를 변경 파일로 지정해 `gosu-review` 스킬을 호출한다 (블라인드 테스트).
2. 스킬 리포트가 아래 "기대 발견 항목"을 모두 잡았는지 확인한다.
3. 누락은 recall 결함, 심지 않은 위반을 보고하면 precision 결함이다.
4. 필요 시 `testgolden/triggerNN_*/fixture.go`의 `// Rule N 위반:` 주석과 대조해 정답을 확인한다.

> `testdata/`는 Go 툴체인이 자동 제외하는 관례 디렉토리라 상위 빌드/테스트에 영향을 주지 않는다.

## 트리거별 기대 발견 항목

| 디렉토리 | 적용 룰 |
|---|---|
| trigger01_function_signature | R7, R8, R11, R42, R43, R44, R45 |
| trigger02_error_handling | R48, R51 |
| trigger03_defer | R35, R46, R47, R53 |
| trigger04_guard_nesting | R2 |
| trigger05_slice_string_memory | R21, R23, R24, R25, R39, R40, R41 |
| trigger06_nil_empty | R26, R27, R28 |
| trigger07_range_iteration | R31, R32, R33, R36, R37 |
| trigger08_nested_loop | R34 |
| trigger09_struct_definition | R4, R10 |
| trigger10_equality | R30 |
| trigger11_map_usage | R29 |
| trigger12_interface_declaration | R5, R6 |
| trigger13_generics | R9 |
| trigger14_package_structure | R12, R13 (디렉토리 전체 입력 필요) |
| trigger15_arithmetic | R18, R19, R20 |
| trigger16_error_declaration | R50 |


