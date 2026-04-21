package trigger08

type Event struct {
	Kind string
}

const stopKind = "stop"

func ProcessEvents(events []Event) int {
	count := 0
	for _, e := range events {
		switch e.Kind {
		case stopKind:
			break
		default:
			count++
		}
	}
	return count
}
