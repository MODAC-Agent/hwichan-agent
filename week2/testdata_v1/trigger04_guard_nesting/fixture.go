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
