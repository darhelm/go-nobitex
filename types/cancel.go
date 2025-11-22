package types

// CancelOrderParams defines the fields used to cancel a single order.
// An order may be referenced either by its numeric Id or by a ClientOrderId.
// The Status field must be set to "canceled".
type CancelOrderParams struct {
	// Id is the numeric identifier of the order being canceled.
	// Optional if ClientOrderId is provided.
	Id int `json:"id,omitempty"`

	// ClientOrderId is an optional client-assigned identifier used
	// to reference and cancel a specific order.
	ClientOrderId string `json:"clientOrderId,omitempty"`

	// Status indicates the desired final state of the order
	// and must be "canceled" for cancellation operations.
	Status string `json:"status"`
}

// CancelOrderBulkParams defines filter criteria for canceling multiple
// orders at once. Any of the fields may be provided to limit which
// orders are eligible for bulk cancellation.
type CancelOrderBulkParams struct {
	// Hours specifies the minimum age of orders to cancel.
	// Orders older than this value (in hours) are included.
	Hours float64 `json:"hours,omitempty"`

	// Execution filters orders by execution type, e.g. "limit" or "market".
	Execution string `json:"execution,omitempty"`

	// TradeType filters orders by market type, such as 'spot' or 'margin'.
	TradeType string `json:"tradeType,omitempty"`

	// SrcCurrency filters orders by base currency of the trading pair.
	SrcCurrency string `json:"srcCurrency,omitempty"`

	// DstCurrency filters orders by quote currency of the trading pair.
	DstCurrency string `json:"dstCurrency,omitempty"`
}

// CancelOrderResponse is always {"status": "ok"} if nothing other than status code 200 is returned
type CancelOrderResponse struct {
	Status string `json:"status"`
}
