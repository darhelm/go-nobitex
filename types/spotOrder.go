package types

import "time"

// CreateOrderResponse represents the full details of an order after it is
// submitted, including pricing, amounts, status, and lifecycle information.
type CreateOrderResponse struct {
	// Type is the order side, such as "buy" or "sell".
	Type string `json:"type"`

	// SrcCurrency is the base asset of the trading pair.
	SrcCurrency string `json:"srcCurrency"`

	// DstCurrency is the quote asset of the trading pair.
	DstCurrency string `json:"dstCurrency"`

	// Price is the order price, represented as a string for precision.
	Price string `json:"price"`

	// Amount is the quantity of the base asset submitted with the order.
	Amount string `json:"amount"`

	// TotalPrice is the quote-value equivalent of the submitted order.
	TotalPrice string `json:"totalPrice"`

	// MatchedAmount is the portion of the order that has already been filled.
	MatchedAmount string `json:"matchedAmount"`

	// UnmatchedAmount is the unfilled portion of the order still remaining.
	UnmatchedAmount string `json:"unmatchedAmount"`

	// Id is the unique identifier assigned to the order.
	Id int `json:"id"`

	// ClientOrderId is an optional client-defined identifier for the order.
	ClientOrderId string `json:"clientOrderId,omitempty"`

	// Status indicates the current lifecycle state of the order,
	// such as "open", "filled", "partial", or "canceled".
	Status string `json:"status"`

	// CreatedAt is the timestamp when the order was created.
	CreatedAt time.Time `json:"created_at"`

	// Partial indicates partial execution details as returned by the API.
	Partial string `json:"partial"`

	// Fee is the total fee charged for executions on this order.
	Fee string `json:"fee"`
}

// CreateOrderParams defines the parameters used to submit a new order,
// including execution type, pricing, and asset information.
type CreateOrderParams struct {
	// Execution specifies how the order should execute, such as "limit" or "market".
	Execution string `json:"execution"`

	// StopPrice is the trigger price for stop or stop-market orders.
	StopPrice string `json:"stopPrice,omitempty"`

	// StopLimitPrice is the trigger price for stop-limit orders.
	StopLimitPrice string `json:"stopLimitPrice,omitempty"`

	// SrcCurrency is the base asset of the order.
	SrcCurrency string `json:"srcCurrency"`

	// DstCurrency is the quote asset of the order.
	DstCurrency string `json:"dstCurrency"`

	// Type specifies the order side, such as "buy" or "sell".
	Type string `json:"type"`

	// Amount specifies the quantity of the base asset to buy or sell.
	Amount string `json:"amount,omitempty"`

	// Price is the limit price for limit-type orders.
	Price string `json:"price,omitempty"`

	// ClientOrderId is an optional client-defined identifier for tracking the order.
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// GetOrderStatusParams specifies filters for retrieving the status of
// a specific order by its id or client-defined identifier.
type GetOrderStatusParams struct {
	// Id is the order identifier assigned by the system.
	Id int `json:"id,omitempty"`

	// ClientOrderId is the client-defined order identifier.
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// OrderStatusResponse contains detailed information about a single order,
// including execution amounts, price, and lifecycle state.
type OrderStatusResponse struct {
	UnmatchedAmount string    `json:"unmatchedAmount"`
	Fee             string    `json:"fee"`
	Partial         bool      `json:"partial"`
	Price           string    `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
	Id              int       `json:"id"`
	SrcCurrency     string    `json:"srcCurrency"`
	DstCurrency     string    `json:"dstCurrency"`
	TotalPrice      string    `json:"totalPrice"`
	Type            string    `json:"type"`
	IsMyOrder       bool      `json:"isMyOrder"`
	Status          string    `json:"status"`
	Amount          string    `json:"amount"`
	ClientOrderId   string    `json:"clientOrderId,omitempty"`
}

// GetOrdersListParams defines the filters used to retrieve a list of orders,
// allowing selection by state, currency, execution type, or id ranges.
type GetOrdersListParams struct {
	Status      string `json:"status"`
	Type        string `json:"type"`
	Execution   string `json:"execution"`
	TradeType   string `json:"tradeType"`
	SrcCurrency string `json:"srcCurrency"`
	DstCurrency string `json:"dstCurrency"`
	Details     int64  `json:"details"`
	FromId      int64  `json:"fromId"`
	Order       string `json:"order"`
}

// OrdersListResponse represents a single entry in an order list,
// providing summary-level data about the order.
type OrdersListResponse struct {
	Id            int       `json:"id,omitempty"`
	Type          string    `json:"type"`
	Execution     string    `json:"execution"`
	Status        string    `json:"status,omitempty"`
	SrcCurrency   string    `json:"srcCurrency"`
	DstCurrency   string    `json:"dstCurrency"`
	Price         string    `json:"price"`
	Amount        string    `json:"amount"`
	MatchedAmount string    `json:"matchedAmount"`
	AveragePrice  string    `json:"averagePrice,omitempty"`
	Fee           string    `json:"fee,omitempty"`
	ClientOrderId string    `json:"clientOrderId"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

// UserTradeResponse represents an executed trade associated with an order,
// including trade amounts, execution price, and fee data.
type UserTradeResponse struct {
	Id          int       `json:"id"`
	OrderId     string    `json:"orderId"`
	SrcCurrency string    `json:"srcCurrency"`
	DstCurrency string    `json:"dstCurrency"`
	Market      string    `json:"market"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	Price       string    `json:"price"`
	Amount      string    `json:"amount"`
	Total       int       `json:"total"`
	Fee         string    `json:"fee"`
}

// GetUserTradesParams defines filters used for retrieving a user’s trade history,
// such as currency filters or pagination identifiers.
type GetUserTradesParams struct {
	SrcCurrency string `json:"srcCurrency,omitempty"`
	DstCurrency string `json:"dstCurrency,omitempty"`
	FromId      string `json:"fromId,omitempty"`
}

// CreateOrderStatus wraps a single order response together with a status field.
type CreateOrderStatus struct {
	Status      string              `json:"status"`
	OrderStatus CreateOrderResponse `json:"order"`
}

// UserTrades represents a collection of user trades, along with pagination info.
type UserTrades struct {
	Status  string              `json:"status"`
	Trades  []UserTradeResponse `json:"trades"`
	HasNext bool                `json:"hasNext"`
}

// OrderStatusList represents a list of orders together with a status indicator.
type OrderStatusList struct {
	Status string               `json:"status"`
	Orders []OrdersListResponse `json:"orders"`
}

// OrderStatus wraps a single order status entry.
type OrderStatus struct {
	Status string              `json:"status"`
	Order  OrderStatusResponse `json:"order"`
}
