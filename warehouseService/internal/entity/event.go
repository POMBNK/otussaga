package entity

import "time"

type EventInfo struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderEvent struct {
	Info    EventInfo `json:"event_info"`
	Payload Order     `json:"payload"`
}

type Order struct {
	ID    int   `json:"id"`
	Goods []int `json:"goods"`
}
