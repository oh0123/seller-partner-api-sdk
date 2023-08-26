package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/logging"
	"github.com/google/uuid"
	order "github.com/oh0123/seller-partner-api-sdk/codegen/ordersv0"
	sign "github.com/oh0123/seller-partner-api-sdk/pkg/sign"
)

func main() {
	cfg := sign.SignConfig{
		AccessKeyID: "XXXX",
		SecretKey:   "XXXX",
		Region:      "XXXX",
		RoleArn:     "XXXX",
		MaxRetry:    1,
		LogMode:     []aws.ClientLogMode{},
		Logger:      logging.Nop{},
	}

	signer := sign.NewSigner(cfg)

	endpoint := "XXXX"

	reqEditor := order.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("X-Amzn-Requestid", uuid.New().String())
		if err := signer.Sign(req); err != nil {
			return fmt.Errorf("sign error: %s", err)
		}
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return fmt.Errorf("dump request error: %s", err)
		}
		log.Printf("DumpRequest = %s", dump)
		return nil
	})

	rsqEditor := order.WithResponseEditorFn(func(ctx context.Context, rsp *http.Response) error {
		dump, err := httputil.DumpResponse(rsp, true)
		if err != nil {
			return fmt.Errorf("dump response error: %s", err)
		}
		log.Printf("DumpResponse = %s", dump)
		return nil
	})

	client, err := order.NewClientWithResponses(endpoint, reqEditor, rsqEditor)
	if err != nil {
		log.Fatal(err)
	}

	p := order.GetOrdersParams{
		MarketplaceIds: []string{"XXXX"},
	}

	response, err := client.GetOrdersWithResponse(context.Background(), &p)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response.HTTPResponse.StatusCode)
}
