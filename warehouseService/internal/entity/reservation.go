package entity

const (
	CompletedReservationEventType reservationEventType = "reservated"
	FailedReservationEventType    reservationEventType = "failed"
)

type reservationEventType = string

type Reservation struct {
	OrderID int   `json:"order_id"`
	Goods   []int `json:"goods"`
}
