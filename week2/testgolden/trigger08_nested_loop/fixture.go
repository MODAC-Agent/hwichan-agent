package trigger08

type Event struct {
	Kind string
}

const stopKind = "stop"

// Rule 34 위반: for 안 switch의 break는 switch만 탈출. 의도는 for 탈출이라면 레이블 필요.
// 권장: Loop: for ... { switch ... { case stopKind: break Loop } }.
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
