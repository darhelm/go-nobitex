# Go Nobitex

[![Go Reference](https://pkg.go.dev/badge/github.com/darhelm/go-nobitex.svg)](https://pkg.go.dev/github.com/darhelm/go-nobitex)
[![Go Report Card](https://goreportcard.com/badge/github.com/darhelm/go-nobitex)](https://goreportcard.com/report/github.com/darhelm/go-nobitex)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/darhelm/go-nobitex)](https://golang.org/dl/)

A comprehensive, type-safe, and fully documented Go SDK for interacting with the **Nobitex** cryptocurrency exchange API.  
This SDK provides a clean and intuitive interface for authentication, wallet access, order execution, and real-time market data.

## Disclaimer

This SDK is **unofficial** and not officially affiliated with Nobitex.  
Use at your own risk — the author(s) assume no responsibility for any losses resulting from API usage or trading logic.

## Features

- Full implementation of public and private Nobitex endpoints
- Strongly typed request/response models
- Full OTP authentication support (TOTP, compatible with `github.com/pquerna/otp`)
- Auto-refresh of API keys (4 hours or ~30 days depending on `remember` mode)
- Real-time market data: tickers, order books, recent trades
- Wallet management
- Full order lifecycle: create, cancel, bulk-cancel, history, status
- Structured error handling (`APIError`, `RequestError`)
- Clean and maintainable Go codebase

## Installation

```bash
go get github.com/darhelm/go-nobitex
```

## Quick Start

```go
package main

import (
    "fmt"
    nobitex "github.com/darhelm/go-nobitex"
)

func main() {
    client, err := nobitex.NewClient(nobitex.ClientOptions{
        Username:    "user@example.com",
        Password:    "your-password",
        OtpCode:     "123456",
        Remember:    "yes",
        UserAgent:   "MyBot/1.0",
        AutoRefresh: true,
    })
    if err != nil {
        panic(err)
    }

    cfg, err := client.GetNobitexConfig()
    if err != nil {
        panic(err)
    }

    fmt.Println(cfg.Nobitex.ActiveCurrencies)
}
```

## Documentation

- Go SDK docs: https://pkg.go.dev/github.com/darhelm/go-nobitex
- Nobitex API: https://apidocs.nobitex.ir/
- Full examples: `EXAMPLES.md`

---

# Examples

## Authentication

```go
client, err := nobitex.NewClient(nobitex.ClientOptions{
    Username:  "user@example.com",
    Password:  "pass",
    OtpCode:   "123456",
    Remember:  "yes",
    UserAgent: "Bot/1.0",
})
```

## Get Nobitex Config

```go
cfg, err := client.GetNobitexConfig()
fmt.Println(cfg.Nobitex.AllCurrencies)
```

## Get Tickers

```go
tickers, err := client.GetTickers(types.GetTickersParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
for symbol, t := range tickers.Stats {
    fmt.Println(symbol, t.Latest)
}
```

## Get Order Book

```go
orderBook, err := client.GetOrderBook("btc-usdt")
fmt.Println(ob.Asks[0], ob.Bids[0])
```

## Get Recent Trades

```go
recentTrades, err := client.GetRecentTrades("btc-usdt")
```

## Get Wallets

```go
balances, err := client.GetWallets(types.GetWalletParams{
    Currencies: []string{"BTC", "USDT"},
})
```

## Create Order

```go
createOrder, err := client.CreateOrder(types.CreateOrderParams{
    Execution:   "limit",
    Type:        "buy",
    SrcCurrency: "btc",
    DstCurrency: "usdt",
    Amount:      "0.01",
    Price:       "2000000000",
})
```

## Cancel Order

```go
cancel, err = client.CancelOrder(types.CancelOrderParams{
    Id: 123,
})
```

## Bulk Cancel

```go
bulkCancel, err = client.CancelOrderBulk(types.CancelOrderBulkParams{
    Hours:       1,
    Execution:   "limit",
    TradeType:   types.TradeTypeSpot,
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
```

## Order History

```go
orderHistory, _ := client.GetOrdersHistory(types.GetOrdersListParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
```

## Open Orders

```go
openOrders, _ := client.GetOpenOrders(types.GetOrdersListParams{})
```

## Order Status

```go
orderStatus, _ := client.GetOrderStatus(types.GetOrderStatusParams{
    Id: 123456,
})
```

## User Trades

```go
userTrades, _ := client.GetUserTrades(types.GetUserTradesParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
```

## Error Handling

```go
if err != nil {
    if apiErr, ok := err.(*nobitex.APIError); ok {
        fmt.Println(apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Detail)
    }
}
```

## Contributing

1. Fork the repository
2. Create a branch in this example `feat/new-feature`
3. Commit changes
4. Open Pull Request

Before pushing:

```bash
go vet ./...
golangci-lint run
```

## License

MIT License.