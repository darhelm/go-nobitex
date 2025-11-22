package types

// Wallet represents a user’s wallet entry for a specific currency,
// including available and blocked balances.
type Wallet struct {
	// Id is the unique identifier assigned to this wallet entry.
	Id int `json:"id"`

	// Balance is the available amount of the asset that can be used
	// for trading, withdrawal, or transfer.
	Balance string `json:"balance"`

	// Blocked is the portion of the balance that is locked due to
	// open orders, withdrawals, or system holds.
	Blocked string `json:"blocked"`
}

// GetWalletParams defines optional filters for retrieving wallet data,
// such as narrowing results by currency or trade type.
type GetWalletParams struct {
	// Currencies specifies a list of asset symbols to filter on.
	// If empty, all wallets are returned.
	Currencies []string `json:"assets,omitempty"`

	// TradeType restricts results to wallets associated with a specific
	// market type, such as 'spot' or 'margin'.
	TradeType string `json:"type,omitempty"`
}

// Wallets represents a collection of wallet entries,
// keyed by currency symbol and grouped under a status field.
type Wallets struct {
	// Status indicates the result of the wallet retrieval operation.
	Status string `json:"status"`

	// Wallets is a map of currency symbols to their corresponding wallet data.
	Wallets map[string]Wallet `json:"wallets"`
}
