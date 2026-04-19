package trigger12

type Notifier interface {
	Notify(msg string) error
}

type EmailNotifier struct {
	From string
}

func (e *EmailNotifier) Notify(msg string) error {
	_ = msg
	return nil
}
