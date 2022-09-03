package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	_, present := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if !present {
		log.Fatal("AWS_WEB_IDENTITY_TOKEN_FILE not set")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	crd, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("{ \"Version\": 1, "+
		"\"AccessKeyId\": \"%s\", "+
		"\"SecretAccessKey\": \"%s\", "+
		"\"SessionToken\": \"%s\", "+
		"\"Expiration\": \"%s\" }",
		crd.AccessKeyID,
		crd.SecretAccessKey,
		crd.SessionToken,
		crd.Expires.Format(time.RFC3339))
}