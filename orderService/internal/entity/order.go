package entity

const (
	NewOrderEventType    EventType = "NewOrder"
	FailedOrderEventType EventType = "FailedOrder"
)

type Order struct {
	ID    int   `json:"id"`
	Goods []int `json:"goods"`
}

type FailedOrderEvent struct {
	ID    int   `json:"order_id"`
	Goods []int `json:"goods"`
}

type OrderEvent struct {
	Info    EventInfo `json:"event_info"`
	Payload Order     `json:"payload"`
}
