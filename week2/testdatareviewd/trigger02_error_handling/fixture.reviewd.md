# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 2건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L20-23 `LoadUser` — `if err != nil` 블록
**[Rule 51] 에러 이중 처리**
- 근거: `fetchUser`에서 받은 에러를 `log.Printf`로 로깅하면서 동시에 동일한 에러를 호출자에게 반환한다. 호출 스택 상위에서 다시 처리/로깅하면 같은 에러가 중복 기록된다. 처리 책임이 있는 최상위 계층에서만 로깅하거나, 래핑해서 반환 중 하나만 선택해야 한다.
- 제안:
  ```go
  // before
  func LoadUser(id string) (string, error) {
      name, err := fetchUser(id)
      if err != nil {
          log.Printf("failed to load user %s: %v", id, err)
          return "", err
      }
      return name, nil
  }
  // after — 래핑해서 반환만 수행, 로깅은 상위에서 일괄 처리
  func LoadUser(id string) (string, error) {
      name, err := fetchUser(id)
      if err != nil {
          return "", fmt.Errorf("LoadUser(%s): %w", id, err)
      }
      return name, nil
  }
  ```

#### L30 `LoadUserWrapped` — `fmt.Errorf("...: %v", err)`
**[Rule 48] 에러 래핑 전략 (`%w` vs `%v`)**
- 근거: 패키지 레벨 센티넬 `ErrUserNotFound`가 선언되어 있고 `fetchUser`가 이를 반환한다. 호출자가 `errors.Is(err, ErrUserNotFound)`로 분기할 수 있도록 경계 에러를 보존하는 것이 자연스러운데, 현재 `%v`를 사용해 원본 에러와의 결합을 끊어버린다. 센티넬 에러를 외부에 노출한 설계와 `%v`가 상충한다. 호출자가 타입 분기를 해야 하는 경계 에러라면 `%w`가 의도에 맞다.
- 제안:
  ```go
  // before
  return "", fmt.Errorf("LoadUserWrapped(%s): %v", id, err)
  // after — 센티넬을 보존해 호출자가 errors.Is로 분기 가능하게
  return "", fmt.Errorf("LoadUserWrapped(%s): %w", id, err)
  ```
- 참고: 반대로 `fetchUser` 내부 에러가 구현 디테일이라 외부 결합을 차단하고 싶은 의도라면 `%v` 유지가 옳다. 공개 API의 계약에 따라 결정할 것 — 검토 권장.
