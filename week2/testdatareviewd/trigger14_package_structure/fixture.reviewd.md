# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 3개 (`common/shared.go`, `pricing/price.go`, `util/misc.go`)
- 새로 발견된 문제: 3건 (변경 라인 — 전체를 🔴로 간주)
- 기존 문제: 0건

## 프로젝트 구조

디렉토리 구조 개요:
```
trigger14_package_structure/
├── common/
│   └── shared.go       (package common)
├── pricing/
│   └── price.go        (package pricing)
└── util/
    └── misc.go         (package util)
```

### 🔴 [Rule 13] 패키지 네이밍 — `common`, `util` 은 책임이 드러나지 않는 이름

**대상**: `common/shared.go` (package common), `util/misc.go` (package util)

- **근거**:
  - `common` 과 `util` 은 "어떤 책임이든 들어갈 수 있는 폴더" 이름의 전형적인 예다. Rule 13은 `common`, `util`, `shared`, `helper`, `misc` 같은 이름을 명시적으로 문제로 지적한다.
  - `common/shared.go` 는 `Pair` 타입과 `MakePair` 헬퍼를 제공하는데, 실제 사용 도메인이 무엇인지 패키지 이름만으로는 알 수 없다.
  - `util/misc.go` 는 `Zero`, `CoalesceString` 두 함수를 담고 있으며, 서로 관련 없는 함수가 "잡동사니" 패키지에 함께 묶인 형태다.

- **제안**:
  - `common` → 실제 책임에 맞는 이름으로 변경. `Pair` 타입이 설정/메타데이터용이라면 `meta`, KV 쌍을 다루는 도메인이 특정된다면 해당 도메인 이름 사용.
  - `util` → 기능별로 분리. 예를 들어 `CoalesceString` 은 `stringx` 또는 `strutil` 로, `Zero` 가 특정 도메인 기본값이라면 해당 도메인 패키지로 이동. 책임이 다른 함수를 한 패키지에 묶지 않도록 한다.

  ```go
  // before
  package common   // common/shared.go
  package util     // util/misc.go

  // after (예시)
  package meta     // meta/pair.go  — KV 쌍 관련 타입
  package stringx  // stringx/coalesce.go  — 문자열 유틸
  // Zero() 는 소속 도메인 패키지로 이동
  ```

### 🔴 [Rule 12] 프로젝트 레이아웃 — 표준 레이아웃 미적용

**대상**: 프로젝트 루트 (`trigger14_package_structure/`)

- **근거**:
  - 최상위에 `cmd/`, `internal/`, `pkg/` 같은 golang-standards/project-layout 관례가 없고, 도메인 패키지(`pricing`)와 유틸/공통 패키지(`common`, `util`)가 동일 레벨에 혼재한다.
  - `pricing` 은 명확한 도메인 책임을 가지므로 `internal/pricing` 또는 `pkg/pricing` 으로 배치하는 것이 표준 레이아웃에 부합한다.
  - 새 프로젝트 또는 초기 구조 설계 단계로 보이므로, 표준 레이아웃 도입을 권장한다.

- **제안**:
  ```
  // before
  trigger14_package_structure/
  ├── common/
  ├── pricing/
  └── util/

  // after (권장 예시)
  trigger14_package_structure/
  ├── cmd/
  │   └── main.go
  ├── internal/
  │   └── pricing/
  │       └── price.go
  └── pkg/
      └── meta/        (구 common → 책임 명칭으로 변경)
          └── pair.go
  ```
  - `internal/` 안에 두면 외부 모듈에서 직접 임포트 불가 — 라이브러리가 아닌 애플리케이션 내부 코드에 적합.
  - `pkg/` 는 외부 공개 목적의 재사용 가능한 패키지에 사용.

---

## common/shared.go

### 🔴 변경 라인의 문제

#### L1 `package common`
**[Rule 13] 패키지 네이밍**
- 근거: `common` 은 책임이 드러나지 않는 이름이다. `Pair`/`MakePair` 가 어느 도메인의 KV 쌍인지 알 수 없어 응집도가 낮아진다.
- 제안:
  ```go
  // before
  package common

  // after
  package meta  // 또는 실제 사용 도메인에 맞는 이름
  ```

---

## util/misc.go

### 🔴 변경 라인의 문제

#### L1 `package util`
**[Rule 13] 패키지 네이밍**
- 근거: `util` 은 Rule 13이 명시적으로 지적하는 모호한 이름이다. 더불어 파일명 `misc.go` 역시 "잡동사니" 를 암시한다. `Zero()` 와 `CoalesceString()` 은 서로 다른 책임이므로 하나의 패키지에 묶는 것은 응집도 저하를 초래한다.
- 제안:
  ```go
  // before
  package util  // util/misc.go

  // after
  package stringx  // stringx/coalesce.go — 문자열 관련
  // Zero() 는 해당 도메인 패키지로 이동하거나 제거 검토
  ```

---

## pricing/price.go

_트리거 14 범위 내에서 `package pricing` 은 명확한 도메인 책임(가격 계산)을 표현하고 있어 Rule 13 위반 없음._
