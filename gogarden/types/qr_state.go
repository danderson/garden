package types

type QRState string

const (
	QRStateWant    QRState = "wanted"
	QRStateApplied         = "applied"
	QRStateIgnore          = "none"
)
