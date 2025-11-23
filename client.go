package nobitex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	t "github.com/darhelm/go-nobitex/types"
	u "github.com/darhelm/go-nobitex/utils"
)

// Constants defining the API base URL and version.
const (
	// BaseUrl is the root URL for the Nobitex Market API.
	BaseUrl = "https://apiv2.nobitex.ir"
)

// ClientOptions represents the configuration options for creating a new API client.
// These options allow customization of the HTTP client, authentication tokens,
// API credentials, and automatic authentication/refresh behaviors.
type ClientOptions struct {
	// HttpClient is the custom HTTP client to be used for API requests.
	// If nil, the default HTTP client is used.
	HttpClient *http.Client

	// Timeout specifies the request timeout duration for the HTTP client.
	Timeout time.Duration

	// BaseUrl is the base URL of the API. Defaults to the constant BaseUrl
	// if not provided.
	BaseUrl string

	Username  string
	Password  string
	OtpSecret string
	OtpCode   string

	Remember  string
	UserAgent string

	// ApiKey is the token used for authenticated API requests.
	ApiKey string

	// AutoAuth enables automatic authentication if no valid tokens are provided.
	AutoAuth bool

	// AutoRefresh enables automatic refreshing of the access token when it expires.
	AutoRefresh bool
}

// Client represents the API client for interacting with the Nobitex Market API.
// It manages authentication, base URL, and API requests.
type Client struct {
	// HttpClient is the HTTP client used for API requests.
	// Defaults to the Go standard library's http.DefaultClient.
	HttpClient *http.Client

	// BaseUrl is the base URL of the API used by this client.
	// Defaults to the constant BaseUrl.
	BaseUrl string

	Username  string
	Password  string
	OtpSecret string
	OtpCode   string

	Remember  string
	UserAgent string

	AuthTime time.Time

	// ApiKey is the API key for authentication.
	ApiKey string

	// AutoAuth enables automatic authentication if no valid tokens are provided.
	AutoAuth bool

	// AutoRefresh enables automatic refreshing of the access token when it expires.
	AutoRefresh bool
}

// NewClient initializes a new Nobitex API client using the provided configuration
// options. It prepares the HTTP client, authentication settings, TOTP handling,
// base URL, and automatic re-authentication behavior.
//
// Parameters:
//   - opts: A ClientOptions struct containing configuration fields such as:
//   - HttpClient: optional custom HTTP client.
//   - Timeout: request timeout used when no custom client is provided.
//   - BaseUrl: optional override for the Nobitex base URL.
//   - Username / Password: credentials for API login.
//   - OtpSecret / OtpCode: TOTP configuration for X-TOTP header.
//   - ApiKey: an already-issued Nobitex API key (optional).
//   - Remember: long-lived ("yes") vs short-lived ("no") login sessions.
//   - AutoAuth: whether to auto-authenticate if ApiKey is missing.
//   - AutoRefresh: whether to re-authenticate when Remember expires.
//
// Returns:
//   - A pointer to an initialized Client.
//   - An error if automatic authentication fails or TOTP generation fails.
//
// Behavior:
//   - If opts.BaseUrl is provided, it overrides the default base URL
//     ("https://apiv2.nobitex.ir/").
//   - If opts.HttpClient is nil, a new http.Client with opts.Timeout is created.
//   - If opts.ApiKey is empty and username/password+TOTP are provided,
//     NewClient performs an immediate login by calling Authenticate().
//   - If opts.OtpSecret is provided, NewClient automatically generates a TOTP
//     one-time password using utils.GenerateOtpCode.
//   - Sets AuthTime to the time of successful authentication.
//   - If AutoRefresh is enabled, immediately checks whether the API key must be
//     refreshed based on Remember ("yes" = ~30 days, "no" = ~4 hours).
//
// Dependencies:
//   - utils.GenerateOtpCode for TOTP generation.
//   - Authenticate() for obtaining a Nobitex API key.
//   - handleAutoRefresh() for Remember-based token rotation.
//
// Errors:
//   - Returned directly from Authenticate() on login failure.
//   - Returned if TOTP generation fails.
//   - Returned from handleAutoRefresh() if refresh logic cannot proceed.
//
// Example:
//
//	opts := ClientOptions{
//	    Username:   "user@example.com",
//	    Password:   "password",
//	    OtpSecret:  "AuthSecret",
//	    Remember:   "yes",
//	    AutoAuth:   true,
//	    AutoRefresh: true,
//	}
//	client, err := NewClient(opts)
//	if err != nil {
//	    log.Fatalf("failed to initialize client: %v", err)
//	}
func NewClient(opts ClientOptions) (*Client, error) {
	client := &Client{
		AutoRefresh: opts.AutoRefresh,
		BaseUrl:     BaseUrl,
		Remember:    opts.Remember,
	}

	if opts.BaseUrl != "" {
		client.BaseUrl = opts.BaseUrl
	}

	if opts.UserAgent != "" {
		client.UserAgent = opts.UserAgent
	}

	if opts.HttpClient != nil {
		client.HttpClient = opts.HttpClient
	} else {
		client.HttpClient = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	client.Username = opts.Username
	client.Password = opts.Password

	if opts.ApiKey == "" && opts.Username != "" && opts.Password != "" && (opts.OtpCode != "" || opts.OtpSecret != "") {
		if opts.OtpSecret != "" {
			client.OtpSecret = opts.OtpSecret
			code, err := u.GenerateOtpCode(opts.OtpSecret)
			if err != nil {
				return nil, err
			}
			client.OtpCode = code
		}

		if opts.OtpCode != "" {
			client.OtpCode = opts.OtpCode
		}

		if _, err := client.Authenticate(opts.Username, opts.Password); err != nil {
			return nil, err
		}
	} else {
		client.ApiKey = opts.ApiKey
	}

	client.AuthTime = time.Now()

	if err := client.handleAutoRefresh(); err != nil {
		return nil, err
	}

	return client, nil
}

// assertAuth validates that the client is currently authenticated by checking
// whether an ApiKey is available.
//
// Parameters:
//   - client: A pointer to the Client instance.
//
// Returns:
//   - nil if the client contains a non-empty ApiKey.
//   - A *GoNobitexError if ApiKey is empty.
//
// Behavior:
//   - This function does not perform any network I/O.
//   - It is used internally before making authenticated requests.
//
// Errors:
//   - "access token is empty" when ApiKey is missing.
//
// Example:
//
//	if err := assertAuth(client); err != nil {
//	    return err
//	}
func assertAuth(client *Client) error {
	if client.ApiKey == "" {
		return &GoNobitexError{
			Message: "access token is empty",
			Err:     nil,
		}
	}
	return nil
}

// createApiURI constructs a fully qualified Nobitex API URL using the client's
// base URL, an optional version prefix, and the raw endpoint path.
//
// Parameters:
//   - endpoint: The API endpoint (MUST begin with a leading slash), for example:
//     "/market/orders/add"
//     "/options"
//     "/orderbook/BTCUSDT/"
//   - version: Optional API version string such as "v2" or "v3".
//     If empty, no version segment is inserted.
//
// Returns:
//   - A fully qualified URL as:
//   - Without version:  {BaseUrl}{endpoint}
//   - With version:     {BaseUrl}/{version}{endpoint}
//
// Behavior:
//   - BaseUrl must NOT have a trailing slash (e.g. "https://apiv2.nobitex.ir").
//   - Endpoint MUST begin with "/", and is appended as-is.
//   - Version MUST NOT begin with "/", the function prepends one automatically.
//
// Examples:
//
//	c.BaseUrl = "https://apiv2.nobitex.ir"
//
//	createApiURI("/market/stats", "")
//	→ "https://apiv2.nobitex.ir/market/stats"
//
//	createApiURI("/options", "v2")
//	→ "https://apiv2.nobitex.ir/v2/options"
//
//	createApiURI("/orderbook/BTCUSDT", "v3")
//	→ "https://apiv2.nobitex.ir/v3/orderbook/BTCUSDT"
func (c *Client) createApiURI(endpoint string, version string) string {
	if version == "" {
		return fmt.Sprintf("%s%s", c.BaseUrl, endpoint)
	}

	return fmt.Sprintf("%s/%s%s", c.BaseUrl, version, endpoint)
}

// handleAutoRefresh enforces Nobitex's Remember-based session lifetime rules.
// It automatically re-authenticates when the API key is too old.
//
// Returns:
//   - nil if the API key is still valid or was refreshed successfully.
//   - A *GoNobitexError if the TOTP secret is missing or refresh fails.
//
// Behavior:
//   - Nobitex API keys are not JWT tokens; their freshness must be managed
//     manually. The SDK uses the Remember field to infer expected validity:
//   - RememberNo  → expires in ~4 hours
//   - RememberYes → expires in ~30 days
//   - If the elapsed time since AuthTime exceeds the expected TTL,
//     handleAutoRefresh generates a new TOTP and re-executes Authenticate().
//   - Updates c.OtpCode automatically when using OtpSecret.
//
// Dependencies:
//   - utils.GenerateOtpCode
//   - Authenticate()
//
// Errors:
//   - "Otp secret is empty" when OtpSecret is required.
//   - Errors from TOTP generation or Authenticate().
//
// Example:
//
//	if c.AutoRefresh {
//	    if err := c.handleAutoRefresh(); err != nil {
//	        return err
//	    }
//	}
func (c *Client) handleAutoRefresh() error {
	if c.OtpSecret == "" {
		return &GoNobitexError{
			Message: "Otp secret is empty, can't refresh Api Key!",
		}
	}

	var ttl time.Duration

	switch c.Remember {
	case "no", "":
		ttl = 4 * time.Hour
	case "yes":
		ttl = 30 * 24 * time.Hour
	default:
		return &GoNobitexError{
			Message: fmt.Sprintf("unknown Remember value: %s", c.Remember),
		}
	}

	if time.Since(c.AuthTime) > ttl {
		code, err := u.GenerateOtpCode(c.OtpSecret)
		c.OtpCode = code
		if err != nil {
			return &GoNobitexError{
				Message: fmt.Sprintf("failed to generate Otp code: %v", err),
				Err:     err,
			}
		}

		if _, err := c.Authenticate(c.Username, c.Password); err != nil {
			return &GoNobitexError{
				Message: fmt.Sprintf("failed to refresh Api Key with Remember = '%s'", c.Remember),
				Err:     err,
			}
		}
	}

	return nil
}

// Request sends an HTTP request to a full Nobitex URL and handles:
//   - optional authentication
//   - automatic TOTP header placement
//   - request serialization (JSON or query params)
//   - response deserialization
//   - structured API error handling
//
// Parameters:
//   - method: "GET" or "POST".
//   - url: Fully constructed URL (already includes version when needed).
//   - auth: Whether the request requires an Authorization header.
//   - otpRequired: Whether X-TOTP should be sent for this specific request.
//   - body: For GET, serialized into URL params; for POST, JSON-encoded.
//   - result: Optional pointer to the output struct into which JSON is unmarshaled.
//
// Returns:
//   - nil on success.
//   - *RequestError on network/encoding issues.
//   - *APIError when Nobitex returns status != 2xx.
//
// Behavior:
//   - GET: struct → ?a=b&c=d via StructToURLParams.
//   - POST: struct → JSON in request body.
//   - If auth=true:
//   - handleAutoRefresh() is executed when AutoRefresh is enabled.
//   - assertAuth() ensures ApiKey is set.
//   - User-Agent and Authorization headers are required.
//   - X-TOTP header is added when otpRequired=true.
//   - On error HTTP status, parseErrorResponse() maps Nobitex JSON error objects
//     into APIError (fields: status, code, message, detail).
//
// Dependencies:
//   - StructToURLParams
//   - assertAuth()
//   - handleAutoRefresh()
//   - parseErrorResponse()
//
// Errors:
//   - "failed to marshal request body"
//   - "failed to convert struct to URL params"
//   - "failed to send request"
//   - "failed to unmarshal response"
//   - APIError (status, code, message, detail)
//
// Example:
//
//	var res t.OrderStatus
//	err := client.Request("POST", url, true, true, params, &res)
//	if err != nil {
//	    return err
//	}
func (c *Client) Request(method string, url string, auth bool, otpRequired bool, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if method == "GET" {
		if body != nil {
			urlParams, err := u.StructToURLParams(body)
			if err != nil {
				return &RequestError{
					GoNobitexError: GoNobitexError{
						Message: "failed to convert struct to URL params",
						Err:     err,
					},
					Operation: "preparing request parameters",
				}
			}
			url += "?" + urlParams
		}
	}

	if method == "POST" {
		if body != nil {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return &RequestError{
					GoNobitexError: GoNobitexError{
						Message: "failed to marshal request body",
						Err:     err,
					},
					Operation: "preparing request body",
				}
			}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &RequestError{
			GoNobitexError: GoNobitexError{
				Message: "failed to create request",
				Err:     err,
			},
			Operation: "creating request",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if auth {
		if c.AutoRefresh {
			if err := c.handleAutoRefresh(); err != nil {
				return &GoNobitexError{
					Message: "failed to refresh authentication",
					Err:     err,
				}
			}
		}

		if err := assertAuth(c); err != nil {
			return &GoNobitexError{
				Message: "authentication validation failed",
				Err:     err,
			}
		}

		if c.UserAgent == "" {
			return &GoNobitexError{
				Message: "UserAgent is empty! please set UserAgent",
				Err:     nil,
			}
		}

		req.Header.Set("User-Agent", "TraderBot/"+c.UserAgent)
		req.Header.Set("Authorization", "Token "+c.ApiKey)
	}

	if otpRequired {
		if c.UserAgent == "" {
			return &GoNobitexError{
				Message: "UserAgent is empty! please set UserAgent",
				Err:     nil,
			}
		}

		req.Header.Set("User-Agent", "TraderBot/"+c.UserAgent)
		req.Header.Set("X-TOTP", c.OtpCode)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &RequestError{
			GoNobitexError: GoNobitexError{
				Message: "failed to send request",
				Err:     err,
			},
			Operation: "sending request",
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestError{
			GoNobitexError: GoNobitexError{
				Message: "failed to read response body",
				Err:     err,
			},
			Operation: "reading response",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	if result != nil {
		if err = json.Unmarshal(respBody, result); err != nil {
			return &RequestError{
				GoNobitexError: GoNobitexError{
					Message: "failed to unmarshal response",
					Err:     err,
				},
				Operation: "parsing response",
			}
		}
	}

	return nil
}

// ApiRequest is a convenience wrapper that builds a Nobitex API URL using
// createApiURI() and delegates the actual HTTP call to Request().
//
// Parameters:
//   - method: HTTP method ("GET", "POST").
//   - endpoint: The endpoint path, such as "/market/orders/add".
//   - version: Nobitex version string ("v2", "v3"). May be empty.
//   - auth: Whether this call requires Authorization: Token <key>.
//   - otpRequired: Whether this endpoint requires X-TOTP.
//   - body: Struct for GET params or POST JSON body.
//   - result: Destination struct for response JSON.
//
// Returns:
//   - nil on success.
//   - See Request() for structured errors.
//
// Behavior:
//   - Constructs URL = BaseUrl + "/api/{version}/{endpoint}".
//   - Passes all fields to Request() unchanged.
//
// Example:
//
//	var stats t.Tickers
//	err := client.ApiRequest("GET", "/market/stats", "", false, false, params, &stats)
func (c *Client) ApiRequest(method, endpoint string, version string, auth bool, otpRequired bool, body interface{}, result interface{}) error {
	url := c.createApiURI(endpoint, version)
	return c.Request(method, url, auth, otpRequired, body, result)
}

// Authenticate logs in to Nobitex using username, password, captcha="api",
// and an X-TOTP code. It retrieves a session API key for authenticated endpoints.
//
// Endpoint:
//
//	POST /auth/login/
//
// Parameters:
//   - Username: account email.
//   - Password: account password.
//
// The function automatically sends:
//   - captcha="api"
//   - remember=<client.Remember>
//   - X-TOTP header (generated or provided)
//
// Returns:
//   - *AuthenticationResponse containing:
//   - Status ("ok")
//   - Key    (API key used in Authorization: Token <key>)
//   - Device (session device ID)
//   - An error if login fails.
//
// Behavior:
//   - Constructs a JSON object with username, password, captcha, remember.
//   - Sends X-TOTP via Request().
//   - On success, updates c.ApiKey with response.Key.
//   - Does NOT manage refresh tokens (Nobitex does not use them).
//
// Errors:
//   - "Username and/or password are empty"
//   - APIError with status/code/message/detail
//
// Example:
//
//	resp, err := client.Authenticate("me@example.com", "secret")
//	if err != nil {
//	    return err
//	}
//	fmt.Println("API key:", resp.Key)
func (c *Client) Authenticate(Username, Password string) (*t.AuthenticationResponse, error) {
	if Username == "" || Password == "" {
		return nil, &GoNobitexError{
			Message: "Username and/or password are empty",
			Err:     nil,
		}
	}

	reqBody := map[string]string{
		"username": Username,
		"password": Password,
		"captcha":  "api",
		"remember": string(c.Remember),
	}

	var authResponse t.AuthenticationResponse
	err := c.ApiRequest("POST", "/auth/login/", "", false, true, reqBody, &authResponse)

	if err != nil {
		// Check for specific API errors here
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 401:
				return nil, err
			case 429:
				return nil, err
			default:
				return nil, err
			}
		}
		return nil, err
	}

	c.AuthTime = time.Now()
	// Update the client's tokens with the newly received ones
	c.ApiKey = authResponse.Key

	return &authResponse, nil
}

// GetNobitexConfig fetches global Nobitex configuration including:
//   - allCurrencies
//   - activeCurrencies
//   - amountPrecisions[asset]
//   - pricePrecisions[asset]
//
// Endpoint:
//
//	GET /api/v2/options
//
// Returns:
//   - *t.Config containing t.Nobitex metadata.
//
// Behavior:
//   - No authentication required.
//   - Calls ApiRequest("GET", "options", "v2").
//
// Example:
//
//	cfg, err := client.GetNobitexConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(cfg.Nobitex.ActiveCurrencies)
func (c *Client) GetNobitexConfig() (*t.Config, error) {
	var config *t.Config
	err := c.ApiRequest("GET", "/options", "v2", false, false, nil, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetTickers retrieves market statistics for one or more trading pairs.
// It delegates to the Nobitex market stats endpoint.
//
// Endpoint:
//
//	GET /api/market/stats
//
// Parameters:
//   - params: t.GetTickersParams
//   - SrcCurrency
//   - DstCurrency
//
// Returns:
//   - *t.Tickers containing:
//     Status
//     Stats[market] → t.Ticker
//
// Behavior:
//   - No authentication required.
//   - Serializes params as query string.
//
// Example:
//
//	stats, err := client.GetTickers(t.GetTickersParams{SrcCurrency:"btc",DstCurrency:"usdt"})
//	if err != nil { ... }
//	fmt.Println(stats.Stats["BTCUSDT"].Latest)
func (c *Client) GetTickers(params t.GetTickersParams) (*t.Tickers, error) {
	var tickers *t.Tickers
	err := c.ApiRequest("GET", "/market/stats", "", false, false, params, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetOrderBook retrieves aggregated asks and bids for a specific market symbol.
//
// Endpoint:
//
//	GET /api/v3/orderbook/{symbol}/
//
// Parameters:
//   - symbol: market pair in Nobitex format (e.g. "BTCUSDT" or "BTCIRT").
//
// Returns:
//   - *t.OrderBook containing:
//     Status
//     LastUpdate
//     LastTradePrice
//     Asks [][]string
//     Bids [][]string
//
// Behavior:
//   - No authentication required.
//   - Uses version "v3" (Nobitex newest orderbook specification).
//
// Example:
//
//	ob, _ := client.GetOrderBook("BTCUSDT")
//	fmt.Println(ob.Asks[0], ob.Bids[0])
func (c *Client) GetOrderBook(symbol string) (*t.OrderBook, error) {
	var orderBook *t.OrderBook
	err := c.ApiRequest("GET", fmt.Sprintf("/orderbook/%s", symbol), "v3", false, false, nil, &orderBook)
	if err != nil {
		return nil, err
	}
	return orderBook, nil
}

// GetRecentTrades retrieves recent trade executions for a market.
//
// Endpoint:
//
//	GET /api/v2/trades/{symbol}/
//
// Parameters:
//   - symbol: market pair ("BTCUSDT", "ETHIRT", etc.)
//
// Returns:
//   - *t.Trades where each Trade has:
//     Time
//     Price
//     Volume
//     Type ("buy"/"sell")
//
// Behavior:
//   - No authentication required.
//   - Returns trades sorted newest-first.
//
// Example:
//
//	trades, _ := client.GetRecentTrades("BTCUSDT")
//	fmt.Println(trades[0].Price, trades[0].Type)
func (c *Client) GetRecentTrades(symbol string) (*t.Trades, error) {
	var trades *t.Trades
	err := c.ApiRequest("GET", fmt.Sprintf("/trades/%s", symbol), "v2", false, false, nil, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetWallets retrieves a list of wallets for the authenticated user from the API.
// It sends a GET request to the `/wallets` endpoint and returns wallet information
// based on the provided parameters.
//
// Parameters:
//   - params: A `GetWalletParams` struct containing optional filters for querying
//     specific wallets, such as by currency list or trade type ("spot"/"margin").
//
// Returns:
//   - A pointer to a `Wallets` struct containing the list of wallets with fields:
//   - id
//   - balance
//   - blocked (frozen)
//   - An error if the request fails, authentication is missing, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to `/api/v2/wallets`
//   - Requires authentication (`auth=true`)
//   - Response JSON is unmarshalled into `t.Wallets`.
//
// Example:
//
//	params := t.GetWalletParams{
//	    Currencies: []string{"BTC","USDT"},
//	    TradeType:  "spot",
//	}
//	wallets, err := client.GetWallets(params)
//	if err != nil { ... }
//	fmt.Println(wallets.Wallets["btc"].Balance)
//
// Dependencies:
//   - ApiRequest
//
// Errors:
//   - APIError (status, code, message, detail)
//   - RequestError (network/JSON issues)
func (c *Client) GetWallets(params t.GetWalletParams) (*t.Wallets, error) {
	var wallets *t.Wallets
	err := c.ApiRequest("GET", "/wallets", "v2", true, false, params, &wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateOrder submits a new trading order (spot or margin) to Nobitex.
//
// Endpoint:
//
//	POST /market/orders/add
//
// Parameters:
//   - params: t.CreateOrderParams containing:
//     Execution ("limit","market","stop_limit","stop_market")
//     SrcCurrency
//     DstCurrency
//     Type ("buy"/"sell")
//     Amount (optional for some execution types)
//     Price  (required for limit orders)
//     StopPrice / StopLimitPrice (for stop orders)
//     ClientOrderId (optional)
//
// Returns:
//
//   - *t.OrderStatus containing the created order state:
//
//     Status
//     Order → OrderStatusResponse
//
// Behavior:
//   - Requires authentication.
//   - No TOTP required for order placement (otpRequired=false).
//   - Returns server-evaluated matched/unmatched amounts, fees, timestamps.
//
// Example:
//
//	order, err := client.CreateOrder(t.CreateOrderParams{
//	    Execution:   "limit",
//	    SrcCurrency: "btc",
//	    DstCurrency: "usdt",
//	    Type:        "buy",
//	    Amount:      "0.01",
//	    Price:       "1500000000",
//	})
func (c *Client) CreateOrder(params t.CreateOrderParams) (*t.OrderStatus, error) {
	var orderStatus *t.OrderStatus
	err := c.ApiRequest("POST", "/market/orders/add", "", true, false, params, &orderStatus)
	if err != nil {
		return nil, err
	}
	return orderStatus, nil
}

// CancelOrder cancels a single existing order.
//
// Endpoint:
//
//	POST /market/orders/update-status
//
// Parameters:
//   - params: t.CancelOrderParams
//     Id or ClientOrderId
//     Status="canceled"
//
// Returns:
//   - *t.CancelOrderResponse on success, always contains "status": "ok" on status code 200
//   - APIError on failure
//
// Behavior:
//   - Requires authentication.
//   - The SDK enforces params.Status="canceled" automatically.
//
// Example:
//
//	err := client.CancelOrder(t.CancelOrderParams{Id:12345})
func (c *Client) CancelOrder(params t.CancelOrderParams) (*t.CancelOrderResponse, error) {
	params.Status = "canceled"

	var cancelOrderStatus *t.CancelOrderResponse
	err := c.ApiRequest("POST", "/market/orders/update-status", "", true, false, params, &cancelOrderStatus)
	if err != nil {
		return nil, err
	}
	return cancelOrderStatus, nil
}

// CancelOrderBulk cancels multiple orders based on filter criteria.
//
// Endpoint:
//
//	POST /market/orders/cancel-old
//
// Parameters:
//   - params: t.CancelOrderBulkParams
//     Hours: cancel orders older than X hours
//     Execution: limit/market/...
//     TradeType: spot/margin
//     SrcCurrency / DstCurrency
//
// Returns:
//   - *t.CancelOrderResponse on success, always contains "status": "ok" on status code 200
//   - APIError on failure
//
// Behavior:
//   - Requires authentication.
//   - Cancels all matching orders without returning per-order info.
//
// Example:
//
//	err := client.CancelOrderBulk(t.CancelOrderBulkParams{Hours: 6})
func (c *Client) CancelOrderBulk(params t.CancelOrderBulkParams) (*t.CancelOrderResponse, error) {
	var cancelOrderBulkStatus *t.CancelOrderResponse
	err := c.ApiRequest("POST", "/market/orders/cancel-old", "", true, false, params, &cancelOrderBulkStatus)
	if err != nil {
		return nil, err
	}
	return cancelOrderBulkStatus, nil
}

// GetOrdersHistory retrieves historical user orders using filter criteria.
//
// Endpoint:
//
//	GET /market/orders/list
//
// Parameters:
//   - params: t.GetOrdersListParams
//     Status ("open","closed","canceled", etc.)
//     Type ("buy"/"sell")
//     Execution ("limit","market")
//     TradeType ("spot","margin")
//     SrcCurrency / DstCurrency
//     Details, FromId, Order (sorting/id filters)
//
// Returns:
//
//   - *t.OrdersListResponse containing:
//
//     Status
//     Orders []OrdersListResponse
//
// Behavior:
//   - Requires authentication.
//   - Filters determine scope and ordering.
//
// Example:
//
//	resp, _ := client.GetOrdersHistory(t.GetOrdersListParams{
//	    Status:"closed",
//	    SrcCurrency:"btc",
//	    DstCurrency:"usdt",
//	})
func (c *Client) GetOrdersHistory(params t.GetOrdersListParams) (*t.OrdersListResponse, error) {
	var orders *t.OrdersListResponse
	err := c.ApiRequest("GET", "/market/orders/list", "", true, false, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOpenOrders retrieves only active/open orders.
//
// Endpoint:
//
//	GET /market/orders/list
//
// Parameters:
//   - params: t.GetOrdersListParams
//     Status is overridden to "open" internally.
//
// Returns:
//   - *t.OrdersListResponse
//
// Behavior:
//   - Requires authentication.
//   - Same filtering as GetOrdersHistory except Status="open".
//
// Example:
//
//	openOrders, _ := client.GetOpenOrders(t.GetOrdersListParams{})
func (c *Client) GetOpenOrders(params t.GetOrdersListParams) (*t.OrdersListResponse, error) {
	var orders *t.OrdersListResponse
	params.Status = "open" // Automatically filter for active (open) orders
	err := c.ApiRequest("GET", "/market/orders/list", "", true, false, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderStatus retrieves detailed status for a specific order
// using either Id or ClientOrderId.
//
// Endpoint:
//
//	POST /market/orders/status
//
// Parameters:
//   - params: t.GetOrderStatusParams
//
// Returns:
//
//   - *t.OrderStatus containing:
//
//     Status
//     Order → OrderStatusResponse
//
// Behavior:
//   - Requires authentication.
//   - POST is used because Nobitex expects request body.
//
// Example:
//
//	st, _ := client.GetOrderStatus(t.GetOrderStatusParams{Id: 12345})
func (c *Client) GetOrderStatus(params t.GetOrderStatusParams) (*t.OrderStatus, error) {
	var orders *t.OrderStatus
	err := c.ApiRequest("POST", "/market/orders/status", "", true, false, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserTrades retrieves the authenticated user's trade history.
//
// Endpoint:
//
//	GET /market/trades/list
//
// Parameters:
//   - params: t.GetUserTradesParams
//     SrcCurrency
//     DstCurrency
//     FromId (for pagination)
//
// Returns:
//   - *t.UserTrades containing:
//     Status
//     Trades []UserTradeResponse
//     HasNext bool
//
// Behavior:
//   - Requires authentication.
//   - Supports pagination via FromId.
//
// Example:
//
//	trades, _ := client.GetUserTrades(t.GetUserTradesParams{
//	    SrcCurrency:"btc",
//	    DstCurrency:"usdt",
//	})
func (c *Client) GetUserTrades(params t.GetUserTradesParams) (*t.UserTrades, error) {
	var trades *t.UserTrades
	err := c.ApiRequest("GET", "/market/trades/list", "", true, false, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
