# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 0건 (변경 라인)
- 기존 문제: 2건 (참고)

---

## fixture.go

### 🟡 기존 라인의 문제 (참고)

#### L19 `log.Printf("failed to load user %s: %v", id, err)`
**[Rule 51] 에러 이중 처리**
- 근거: 에러를 로깅하고 바로 반환하고 있습니다. 호출 스택 위에서 다시 처리할 경우 동일한 에러가 여러 번 로깅됩니다.
- 제안: 처리 책임이 있는 곳에서만 로깅하고, 그 외에는 래핑하여 반환하세요.
  ```go
  // before
  	if err != nil {
  		log.Printf("failed to load user %s: %v", id, err)
  		return "", err
  	}
  // after
  	if err != nil {
  		return "", fmt.Errorf("failed to load user %s: %w", id, err)
  	}
  ```

#### L28 `return "", fmt.Errorf("LoadUserWrapped(%s): %v", id, err)`
**[Rule 48] 에러 래핑 전략 (`%w` vs `%v`)**
- 근거: `ErrUserNotFound`와 같이 호출자가 확인할 가능성이 있는 에러를 `%v`로 포맷팅하여 원본 에러 체인이 유실됩니다. 호출자가 `errors.Is` 등을 사용할 수 있게 `%w`를 사용하는 것이 좋습니다.
- 제안:
  ```go
  // before
  		return "", fmt.Errorf("LoadUserWrapped(%s): %v", id, err)
  // after
  		return "", fmt.Errorf("LoadUserWrapped(%s): %w", id, err)
  ```
