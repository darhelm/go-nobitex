package types

import "time"

// Nobitex holds general market-wide configuration values,
// including supported currencies and precision settings.
type Nobitex struct {
	// AllCurrencies lists every currency supported by the platform.
	AllCurrencies []string `json:"allCurrencies"`

	// ActiveCurrencies lists currencies currently active for trading.
	ActiveCurrencies []string `json:"activeCurrencies"`

	// AmountPrecisions specifies fractional precision limits for amounts per currency.
	AmountPrecisions map[string]string `json:"amountPrecisions"`

	// PricePrecisions specifies fractional precision limits for prices per currency.
	PricePrecisions map[string]string `json:"pricePrecisions"`
}

// Ticker represents real-time market data for a trading pair,
// including pricing, volume, and daily statistics.
type Ticker struct {
	// IsClosed indicates whether the market is currently halted.
	IsClosed bool `json:"isClosed"`

	// BestSell is the lowest ask price available.
	BestSell string `json:"bestSell"`

	// BestBuy is the highest bid price available.
	BestBuy string `json:"bestBuy"`

	// VolumeSrc is the traded volume of the base currency.
	VolumeSrc string `json:"volumeSrc"`

	// VolumeDst is the traded volume of the quote currency.
	VolumeDst string `json:"volumeDst"`

	// Latest is the most recent trade price.
	Latest string `json:"latest"`

	// Mark is the platform-computed mark price for the market.
	Mark string `json:"mark"`

	// DayLow is the lowest price traded in the last 24 hours.
	DayLow string `json:"dayLow"`

	// DayHigh is the highest price traded in the last 24 hours.
	DayHigh string `json:"dayHigh"`

	// DayOpen is the opening price for the 24h period.
	DayOpen string `json:"dayOpen"`

	// DayClose is the closing price for the last 24h.
	DayClose string `json:"dayClose"`

	// DayChange is the net price change over the 24h period.
	DayChange string `json:"dayChange"`
}

// GetTickersParams specifies filters for retrieving ticker
// data for a specific base and quote currency pair.
type GetTickersParams struct {
	// SrcCurrency is the base currency being traded.
	SrcCurrency string `json:"srcCurrency"`

	// DstCurrency is the quote currency used for pricing.
	DstCurrency string `json:"dstCurrency"`
}

// OrderBook represents the current order book for a market,
// including aggregated lists of bids and asks.
type OrderBook struct {
	// Status indicates the response state for the order book.
	Status string `json:"status"`

	// LastUpdate is a Unix timestamp of the last book update.
	LastUpdate int64 `json:"lastUpdate"`

	// LastTradePrice is the price of the most recent trade.
	LastTradePrice string `json:"lastTradePrice"`

	// Asks lists the available sell orders at each price level.
	// Each entry is [price, quantity].
	Asks [][]string `json:"asks"`

	// Bids lists the available buy orders at each price level.
	// Each entry is [price, quantity].
	Bids [][]string `json:"bids"`
}

// Trade represents a single executed trade, including
// timestamp, price, size, and direction.
type Trade struct {
	// Time is the execution timestamp of the trade.
	Time time.Time `json:"time"`

	// Price is the executed trade price.
	Price string `json:"price"`

	// Volume is the traded amount of the base currency.
	Volume string `json:"volume"`

	// Type is the trade direction, such as "buy" or "sell".
	Type string `json:"type"`
}

// Config wraps high-level platform configuration under a Nobitex key.
type Config struct {
	Nobitex Nobitex `json:"nobitex"`
}

// Tickers represents multiple ticker entries,
// indexed by market symbol.
type Tickers struct {
	// Status indicates the response result state.
	Status string `json:"status"`

	// Stats contains per-market ticker data.
	Stats map[string]Ticker `json:"stats"`
}

// Trades represents a list of recent trades for a market.
type Trades struct {
	// Status indicates the response result state.
	Status string `json:"status"`

	// Trades is the list of individual trade records.
	Trades []Trade `json:"trades"`
}
