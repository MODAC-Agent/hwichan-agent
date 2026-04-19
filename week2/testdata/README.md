# gosu-review 테스트 픽스처

`week2/` gosu-review 스킬의 동작 검증용 Go 픽스처 모음. 각 하위 디렉토리는 `week2/Skill.md`의 트리거 번호와 1:1 매핑된다.

## 사용법

1. 각 `fixture.go` 파일을 변경 파일로 지정해 `gosu-review` 스킬을 호출한다.
2. 스킬 리포트가 아래 "기대 발견 항목"을 모두 잡았는지 확인한다.
3. 누락은 recall 결함, 심지 않은 위반을 보고하면 precision 결함이다.

각 파일의 `// Rule N 위반:` 주석은 심은 위반의 힌트다. 스킬은 이 주석을 근거로 삼지 않고도 같은 룰을 검출해야 한다.

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
