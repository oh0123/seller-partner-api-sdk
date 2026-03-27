package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	ordersV0 "github.com/oh0123/seller-partner-api-sdk/v2/codegen/ordersv0"
	"github.com/oh0123/seller-partner-api-sdk/v2/pkg/middleware"
	"github.com/oh0123/seller-partner-api-sdk/v2/pkg/sign"
)

func main() {
	// You can globally replace the JSON engine with a high-performance one here, e.g., sonic, jsoniter!
	// back to std json lib if you don't need it
	ordersV0.JSONMarshal = sonic.Marshal
	ordersV0.JSONUnmarshal = sonic.Unmarshal

	endpoint := "https://sellingpartnerapi-na.amazon.com"

	// 1. Initialize retry options
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
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if seconds, err := strconv.Atoi(retryAfter); err == nil {
						log.Printf("[Retry] Server requested backoff for %d seconds\n", seconds)
						return time.Duration(seconds) * time.Second
					}
				}
			}
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			log.Printf("[Retry] Exponential backoff for %v\n", backoff)
			return backoff
		},
	}

	// 2. Wrap each interceptor into a Middleware function
	apmMW := middleware.APM()

	// Wrap the internally generated retry struct with a unified middleware signature
	retryMW := func(next http.RoundTripper) http.RoundTripper {
		return &ordersV0.RetryTransport{Options: retryOptions, Next: next}
	}

	headerMW := middleware.HeaderModifier(map[string]string{
		"X-Amz-Requestid":                    uuid.New().String(),
		sign.SIGNED_ACCESS_TOKEN_HEADER_NAME: "Atza|IwEB....",
	})

	middleware.MaxLogBodyLength = 4096

	logMW := middleware.Log(log.Default()) // Default high-performance standard library logging

	// 3. Chain the middlewares together in a single line:
	// Execution order is inward: APM -> Retry -> headerMW -> Log -> DefaultTransport(Network)
	chainedTransport := middleware.Chain(http.DefaultTransport,
		apmMW,
		retryMW,
		headerMW,
		logMW,
	)

	// 4. Inject the chained middlewares into the HTTP Client
	httpClient := &http.Client{
		Transport: chainedTransport,
	}

	// 5. Initialize the generated Client
	// Passing "" utilizes the DefaultServer generated from the OpenAPI `servers` JSON block
	client, err := ordersV0.NewClientWithResponses(endpoint, ordersV0.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("Successfully created SP-API OrdersV0 client utilizing DefaultServer.")

	// 8. Example API Call
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create request parameters
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

	// The generated SDK provides strong types for JSON parsing (JSON200, JSON400, etc.)
	if resp.JSON200 != nil {
		log.Printf("\nSUCCESS! Retrieved orders.\n")
	} else if resp.JSON429 != nil {
		log.Println("\nRATE LIMITED! Exhausted max retries.")
	} else if resp.JSON403 != nil {
		log.Println("\nFORBIDDEN! Check your SP-API credentials.")
	} else {
		// Non-200, unmodeled, or raw bytes fallback
		log.Printf("\nUNHANDLED RESPONSE: Status %d\n", resp.StatusCode())
	}
}
