package trigger04

import "errors"

type User struct {
	Active bool
	Admin  bool
}

type Resource struct {
	Public bool
}

var ErrForbidden = errors.New("forbidden")

// Rule 2 위반: 5단 중첩. happy path가 가장 안쪽에 있음.
// 가드 절(early return)로 평탄화 필요.
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
