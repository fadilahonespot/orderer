package controller

import "time"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type OrderRequest struct {
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []ItemData  `json:"items"`
}

type ItemData struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
