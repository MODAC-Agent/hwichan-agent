package trigger16

import "errors"

func FindByName(name string) error {
	return errors.New("not found: " + name)
}

type NotConnectedError struct{}

func (e *NotConnectedError) Error() string { return "not connected" }

func Connect() error {
	return &NotConnectedError{}
}
