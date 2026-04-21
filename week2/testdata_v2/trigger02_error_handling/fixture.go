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

func LoadUser(id string) (string, error) {
	name, err := fetchUser(id)
	if err != nil {
		log.Printf("failed to load user %s: %v", id, err)
		return "", err
	}
	return name, nil
}

func LoadUserWrapped(id string) (string, error) {
	name, err := fetchUser(id)
	if err != nil {
		return "", fmt.Errorf("LoadUserWrapped(%s): %v", id, err)
	}
	return name, nil
}
