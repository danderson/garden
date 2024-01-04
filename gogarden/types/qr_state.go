package types

import "fmt"

type QRState int

const (
	QRStateWant QRState = iota
	QRStateApplied
	QRStateIgnore
	QRStateUnknown
)

//go:generate stringer -type=QRState -trimprefix=QRState

func (s *QRState) MarshalText() ([]byte, error) {
	switch *s {
	case QRStateWant:
		return []byte("want"), nil
	case QRStateApplied:
		return []byte("applied"), nil
	case QRStateIgnore:
		return []byte("ignore"), nil
	default:
		return []byte(""), nil
	}
}

func (s *QRState) UnmarshalText(text []byte) error {
	switch string(text) {
	case "want":
		*s = QRStateWant
	case "applied":
		*s = QRStateApplied
	case "ignore":
		*s = QRStateIgnore
	default:
		return fmt.Errorf("unknown QR state %q", text)
	}
	return nil
}
