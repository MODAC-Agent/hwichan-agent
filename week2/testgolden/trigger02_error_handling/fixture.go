package trigger02

import (
	"errors"
	"fmt"
	"log"
)

var ErrUserNotFound = errors.New("user not found")

func fetchUser(id string) (string, error) {
	if id == "" {
		return "", ErrUserNotFound
	}
	return "alice", nil
}

// Rule 51 위반: 같은 에러를 로깅도 하고 반환도 함.
// 상위 레이어에서 또 로깅하면 중복 로그가 쌓임.
func LoadUser(id string) (string, error) {
	name, err := fetchUser(id)
	if err != nil {
		log.Printf("failed to load user %s: %v", id, err)
		return "", err
	}
	return name, nil
}

// Rule 48 위반: 경계 에러(ErrUserNotFound)인데 %v로 래핑해 errors.Is 분기 차단.
// %w로 감싸야 호출자가 errors.Is(err, ErrUserNotFound) 판정 가능.
func LoadUserWrapped(id string) (string, error) {
	name, err := fetchUser(id)
	if err != nil {
		return "", fmt.Errorf("LoadUserWrapped(%s): %v", id, err)
	}
	return name, nil
}
