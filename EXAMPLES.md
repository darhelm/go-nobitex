# Nobitex SDK Examples

This document contains full practical usage examples for all available SDK methods.

---

# Authentication

```go
client, err := nobitex.NewClient(nobitex.ClientOptions{
    Username:    "user@example.com",
    Password:    "strong-password",
    OtpCode:     "123456",
    Remember:    "yes",
    UserAgent:   "MyBot/1.0",
    AutoRefresh: true,
})
```

# Market Information

## Get Nobitex Config

```go
cfg, err := client.GetNobitexConfig()
fmt.Println(cfg.Nobitex.ActiveCurrencies)
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
orderBook, err := client.GetOrderBook("BTCUSDT")
fmt.Println("Best Ask:", ob.Asks[0])
fmt.Println("Best Bid:", ob.Bids[0])
```

## Get Recent Trades

```go
recentTrades, err := client.GetRecentTrades("BTCUSDT")
fmt.Println(recent.Trades[0])
```

# Wallet Operations

## Get Wallets

```go
balances, err := client.GetWallets(types.GetWalletParams{
    Currencies: []string{"BTC", "USDT"},
})
fmt.Println(wallets.Wallets["BTC"].Balance)
```

# Trading

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
fmt.Println(order.Id)
```

## Cancel Order

```go
cancelOrder, err := client.CancelOrder(types.CancelOrderParams{
    Id: 555,
})
```

## Bulk Cancel Orders

```go
bulkCancel, err := client.CancelOrderBulk(types.CancelOrderBulkParams{
    Hours:       2,
    Execution:   "limit",
    TradeType:   types.TradeTypeSpot,
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
```

## Get Orders History

```go
orderHistory, err := client.GetOrdersHistory(types.GetOrdersListParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
fmt.Println(history.Orders)
```

## Get Open Orders

```go
openOrders, err := client.GetOpenOrders(types.GetOrdersListParams{})
fmt.Println(open.Orders)
```

## Get Order Status

```go
orderStatus, err := client.GetOrderStatus(types.GetOrderStatusParams{
    Id: 9999,
})
fmt.Println(st)
```

# User Trades

```go
trades, err := client.GetUserTrades(types.GetUserTradesParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
fmt.Println(trades.Trades)
```

# Error Handling

```go
_, err := client.GetNobitexConfig()
if err != nil {
    if apiErr, ok := err.(*nobitex.APIError); ok {
        fmt.Println("Status:", apiErr.Status)
        fmt.Println("Code:", apiErr.Code)
        fmt.Println("Message:", apiErr.Message)
        fmt.Println("Detail:", apiErr.Detail)
    }
}
```
