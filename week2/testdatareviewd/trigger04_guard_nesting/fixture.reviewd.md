# gosu-review 리뷰 리포트

## 요약
- 검토 파일: 1개
- 새로 발견된 문제: 1건 (변경 라인)
- 기존 문제: 0건 (참고)

---

## fixture.go

### 🔴 변경 라인의 문제

#### L16-L29 `func Authorize(user *User, resource *Resource, action string) error`
**[Rule 2] 중첩 최소화 (happy path 좌측 정렬)**
- 근거: `if` 중첩 깊이가 5단까지 들어가며(17→18→19→20→21), 정상 경로(`return nil`)가 함수의 가장 안쪽에 숨어 있다. 모든 바깥 분기는 실패 시 `ErrForbidden`으로 떨어지는 예외 처리이므로 가드 절로 뒤집어 happy path를 좌측 정렬해야 한다. 또한 `resource.Public || user.Admin`과 `action == "read" || user.Admin`처럼 조건이 중첩 분기와 OR 조합으로 섞여 있어 인가 로직의 의도가 읽히지 않는다.
- 제안:
  ```go
  // before
  func Authorize(user *User, resource *Resource, action string) error {
      if user != nil {
          if user.Active {
              if resource != nil {
                  if resource.Public || user.Admin {
                      if action == "read" || user.Admin {
                          return nil
                      }
                  }
              }
          }
      }
      return ErrForbidden
  }

  // after
  func Authorize(user *User, resource *Resource, action string) error {
      if user == nil || !user.Active {
          return ErrForbidden
      }
      if resource == nil {
          return ErrForbidden
      }

      // Admin은 모든 동작 허용
      if user.Admin {
          return nil
      }

      // 일반 사용자는 public 리소스의 read만 허용
      if resource.Public && action == "read" {
          return nil
      }

      return ErrForbidden
  }
  ```
  - 적용 우선순위대로 에러 분기를 먼저 조기 반환으로 분리하고, OR로 섞인 관문은 역할별(Admin 우회 / 일반 허용 규칙)로 쪼개 의도를 드러낸다. 필요 시 `action`의 유효값 검증 또는 상수화는 별도 검토 권장.
