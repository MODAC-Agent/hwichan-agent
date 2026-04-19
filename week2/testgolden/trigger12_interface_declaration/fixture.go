package trigger12

// Rule 5 위반: 구현체가 EmailNotifier 하나뿐이고 테스트 모킹도 없는 선제적 인터페이스.
// 현재 실제로 교체 가능성/다형성 요구가 없음 → 구체 타입을 직접 쓰면 충분.
//
// Rule 6 위반: 인터페이스가 구현체 패키지에 함께 선언됨.
// 권장: 사용자(클라이언트) 패키지에 작은 인터페이스를 두고 이 패키지는 구현체만 노출.
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
