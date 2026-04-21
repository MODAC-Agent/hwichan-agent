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

func NewParcel(tracking string) Shippable {
	return Parcel{Tracking: tracking}
}

func Lookup(key string) any {
	if key == "" {
		return 0
	}
	return "value"
}

type MyError struct{ Msg string }

func (e *MyError) Error() string { return e.Msg }

func FindOrder(id string) error {
	var e *MyError
	if id == "" {
		e = &MyError{Msg: "empty id"}
	}
	return e
}

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

func NewClient(host string, port int, user, password string, timeoutSec int, retries int) *Client {
	_ = password
	_ = retries
	return &Client{Host: host, Port: port, User: user, Timeout: timeoutSec}
}

func (c Client) Address() string   { return c.Host }
func (c *Client) SetHost(h string) { c.Host = h }

func Split(s string) (first, rest string) {
	if len(s) == 0 {
		return
	}
	first, rest = s[:1], s[1:]
	return
}
