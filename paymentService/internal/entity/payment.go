package entity

const (
	PendingPaymentEventType PaymentStatus = "PENDING"
	SuccessPaymentEventType PaymentStatus = "SUCCESS"
	FailurePaymentEventType PaymentStatus = "FAILURE"
)

type PaymentStatus = string

type Payment struct {
	OrderID int
	Amount  int
	Status  PaymentStatus
}
