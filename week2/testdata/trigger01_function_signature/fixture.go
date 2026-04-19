package trigger01

import (
	"io"
	"os"
)

type Order struct {
	ID     string
	Amount int
}

type Shippable interface {
	Ship() error
}

type Parcel struct{ Tracking string }

func (p Parcel) Ship() error { return nil }

// Rule 7 위반: 구현체가 Parcel 하나뿐인데 인터페이스를 반환.
// 호출 측에서 구체 타입이 필요한 경우가 많음.
func NewParcel(tracking string) Shippable {
	return Parcel{Tracking: tracking}
}

// Rule 8 위반: 단순히 "여러 타입을 반환하고 싶어서" any 사용.
func Lookup(key string) any {
	if key == "" {
		return 0
	}
	return "value"
}

type MyError struct{ Msg string }

func (e *MyError) Error() string { return e.Msg }

// Rule 44 위반: nil *MyError를 error 인터페이스로 반환.
// 호출 측에서 err != nil이 참이 되는 함정.
func FindOrder(id string) error {
	var e *MyError
	if id == "" {
		e = &MyError{Msg: "empty id"}
	}
	return e
}

// Rule 45 위반: 파일 경로(string)를 받아 내부에서 os.Open. io.Reader로 받아야 테스트/재사용 용이.
func LoadConfig(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

type Client struct {
	Host    string
	Port    int
	User    string
	Timeout int
}

// Rule 11 위반: 파라미터 6개 + 옵션성 필드 다수. 함수형 옵션 패턴이 적합.
func NewClient(host string, port int, user, password string, timeoutSec int, retries int) *Client {
	_ = password
	_ = retries
	return &Client{Host: host, Port: port, User: user, Timeout: timeoutSec}
}

// Rule 42 위반: 같은 타입 Client에 값 리시버와 포인터 리시버가 섞임.
func (c Client) Address() string   { return c.Host }
func (c *Client) SetHost(h string) { c.Host = h }

// Rule 43 위반: 단순 함수에 기명 반환값 + naked return — 읽기 어려움.
func Split(s string) (first, rest string) {
	if len(s) == 0 {
		return
	}
	first, rest = s[:1], s[1:]
	return
}
