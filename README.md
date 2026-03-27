# Amazon Seller-Partner API (SP-API) Go SDK


A modernized, ultra-high-performance, and highly decoupled Go SDK for the [Amazon Seller-Partner API (SP-API)](https://developer-docs.amazon.com/sp-api).

This completely modernized and refactored generator template is tailor-made for high-concurrency cloud environments. We have strictly implemented three gold standards: **High Performance, High Decoupling, and High Extensibility**.

## 🚀 Features & Architecture

### ⚡️ 1. High Performance
* **Streaming API Parsing Engine**: We have completely eliminated the `io.ReadAll` approach from the original implementation, which previously caused a 2x memory overhead for massive JSON payloads (such as parsing huge order books). Replaced with `json.NewDecoder`, this reduces Garbage Collection (GC) pressure and minimizes memory fragmentation.
* **Zero-Alloc Logging**: The built-in log interceptor `middleware.Log` utilizes `io.TeeReader` black magic to implement "side-channel packet capture". As the JSON parser systematically consumes the underlying network stream, the data is concurrently streaming to the logging mechanism. It logs while it parses, effectively eradicating Out-of-Memory (OOM) errors caused by massive request/response payloads. Built-in structured Debug/Info/Error levels.

### 🧩 2. High Decoupling
* **Global JSON Engine Hot-Replacement**: An incredible boon for massive concurrent loads! We explicitly expose top-level variables (`JSONMarshal`, `JSONUnmarshal`, `JSONNewDecoder`, `JSONNewEncoder`) in the generated SDK. By adding a single line in your `main` function like `ordersV0.JSONMarshal = jsoniter.Marshal`, you can seamlessly swap out the standard `encoding/json` with ultra-fast libraries like `sonic`, `jsoniter`, or `fastjson` with absolutely zero code intrusion, boosting serialization speed by over 200%.
* **Onion Model Decoupling**: We ditched the rigid, hardcoded `RequestEditors` slice. All underlying interceptors are natively delegated to Go's `http.RoundTripper`. This ensures a clean separation between the HTTP protocol layer and your business logic. 

### 🔌 3. High Extensibility / Pluggability
* **Lego-like Middleware Chains**: We provide `middleware.Chain` for plug-and-play assembly. Need Authentication (`Auth`), Performance Monitoring (`APM`), or Logging (`Log`)? Snap them together like Lego bricks. Add what you need, drop what you don't.
* **Enterprise-Grade Smart Retries**: The generated SDK is equipped with a `RetryTransport` by default, expanding anti-jitter and rate-limit mitigation capabilities. You can fully customize your retry strategies based on precise `x-amzn-RateLimit-Limit` and `Retry-After` headers, or apply an exponential backoff to effortlessly handle server throttling and prevent cascading failures.

## 📦 Requirements

- **Go version**: `>= 1.26`
- Active Amazon SP-API Developer Credentials

## ⚙️ Installation

```bash
go get github.com/oh0123/seller-partner-api-sdk@latest
```

*(Note: Ensure you run `go mod tidy` after fetching the library to synchronize dependencies.)*

## 📖 How It Works

This SDK utilizes an **Onion Model (Russian Doll)** for its middleware execution. This mechanism ensures that wrappers around the original request execute in exactly the order they point outwards.

```text
 Client Request 
      │
      ▼
┌───────────────────────────────────────┐
│ Middleware A (e.g. Record Start Time) │
│   │                                   │
│   ▼                                   │
│  ┌─────────────────────────────────┐  │
│  │ Middleware B (e.g. Add Auth)    │  │
│  │   │                             │  │
│  │   ▼                             │  │
│  │  ┌───────────────────────────┐  │  │
│  │  │ Middleware C (e.g. Retry) │  │  │
│  │  │   │                       │  │  │
│  │  │   ▼                       │  │  │
│  │  │ [ Actual Network Request ]│  │  │
│  │  │   ▲                       │  │  │
│  │  │   │                       │  │  │
│  │  │ Process response after C  │  │  │
│  │  └───────────────────────────┘  │  │
│  │   ▲                             │  │
│  │ Process response after B         │  │
│  └─────────────────────────────────┘  │
│   ▲                                   │
│ Process response after A              │
└───────────────────────────────────────┘
      │
      ▼
 Return to Client
```

## 💻 Quick Start Example

Below is a complete example of how to configure the client, swap out the JSON engine, set up smart retries, chain middlewares, and invoke the `GetOrders` API.

```go
package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/oh0123/seller-partner-api-sdk/codegen/ordersV0"
	"github.com/oh0123/seller-partner-api-sdk/pkg/middleware"
	"github.com/oh0123/seller-partner-api-sdk/pkg/sign"
)

func main() {
	// ⚡️ GLOBAL JSON ENGINE HOT-REPLACEMENT
	// Replace the standard json lib with a high-performance engine (e.g., sonic)
	ordersV0.JSONMarshal = sonic.Marshal
	ordersV0.JSONUnmarshal = sonic.Unmarshal

	endpoint := "https://sellingpartnerapi-na.amazon.com"

	// 1. Initialize Retry Options
	retryOptions := ordersV0.RetryOptions{
		MaxRetries: 3,
		RetryOn: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode == 429 || resp.StatusCode >= 500
		},
		Backoff: func(attempt int, resp *http.Response) time.Duration {
			if resp != nil {
				// Smart backoff using Amazon's `Retry-After` header if rate-limited
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if seconds, err := strconv.Atoi(retryAfter); err == nil {
						log.Printf("[Retry] Server requested backoff for %d seconds\n", seconds)
						return time.Duration(seconds) * time.Second
					}
				}
			}
			// Fallback to exponential backoff
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			log.Printf("[Retry] Exponential backoff for %v\n", backoff)
			return backoff
		},
	}

	// 2. Wrap Interceptors into Middleware Functions
	apmMW := middleware.APM()

	retryMW := func(next http.RoundTripper) http.RoundTripper {
		return &ordersV0.RetryTransport{Options: retryOptions, Next: next}
	}

	headerMW := middleware.HeaderModifier(map[string]string{
		"X-Amz-Requestid":                    uuid.New().String(),
		sign.SIGNED_ACCESS_TOKEN_HEADER_NAME: "Atza|IwEB....", // Provide a valid LWA Token
	})

	middleware.MaxLogBodyLength = 4096
	logMW := middleware.Log(middleware.DefaultLog) // High-performance zero-alloc logging

	// 3. Chain the middlewares together
	// Execution order is inward: APM -> Retry -> Header -> Log -> DefaultTransport(Network)
	chainedTransport := middleware.Chain(http.DefaultTransport,
		apmMW,
		retryMW,
		headerMW,
		logMW,
	)

	// 4. Inject into the standard HTTP Client
	httpClient := &http.Client{
		Transport: chainedTransport,
	}

	// 5. Initialize the generated Client
	client, err := ordersV0.NewClientWithResponses(endpoint, ordersV0.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("Successfully created SP-API OrdersV0 client.")

	// 6. Execute an API Call!
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	marketplaceIds := []string{"ATVPDKIKX0DER"}
	createdAfter := "2023-01-01T00:00:00Z"
	params := &ordersV0.GetOrdersParams{
		MarketplaceIds: marketplaceIds,
		CreatedAfter:   &createdAfter,
	}

	log.Println("\n--- Invoking GetOrdersWithResponse ---")
	resp, err := client.GetOrdersWithResponse(ctx, params)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	// 7. Typed Error Handling Layout
	if resp.JSON200 != nil {
		log.Printf("\nSUCCESS! Retrieved orders.\n")
	} else if resp.JSON429 != nil {
		log.Println("\nRATE LIMITED! Exhausted max retries.")
	} else if resp.JSON403 != nil {
		log.Println("\nFORBIDDEN! Check your SP-API credentials.")
	} else {
		log.Printf("\nUNHANDLED RESPONSE: Status %d\n", resp.StatusCode())
	}
}
```

## 💡 Tips

* **Global JSON Engine Replacement:** Take maximum advantage of speed! Don't forget to map `JSONNewDecoder` and `JSONNewEncoder` if you are using stream parsing endpoints.
* **Authentication Layer:** Ensure you use the SDK's signing utilities or your own custom integration within the `HeaderModifier` slice to consistently inject proper STS AWS Signature V4 headers and LWA (Login with Amazon) Access Tokens into all your API calls.
* **Middleware Ordering Strategy:** Bear in mind that HTTP Roundtrippers run like an onion. Place your `APM` and `Log` middlewares appropriately. Retries should usually be *wrapped by* APM so that you can correctly measure the total time elapsed (including backoff waiting intervals) across all request lifetimes.

## 🔗 Official Links

* [Amazon Seller-Partner API Documentation](https://developer-docs.amazon.com/sp-api)

## 💬 Community & Contributing

Found a bug or interested in expanding features? Feel free to open an Issue or a Pull Request! Let's build the most robust and high-performance go-to SDK for the Amazon Seller Partner API ecosystem together.
