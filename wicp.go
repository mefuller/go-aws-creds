package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
}

type AssumeRoleWithWebIdentityResult struct {
	Creds Credentials `json:"Credentials"`
}

type AssumeRoleWithWebIdentityResponse struct {
	Result AssumeRoleWithWebIdentityResult `json:"AssumeRoleWithWebIdentityResult"`
}

type Response struct {
	Response AssumeRoleWithWebIdentityResponse `json:"AssumeRoleWithWebIdentityResponse"`
}

func main() {
	if _, ok := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE"); !ok {
		log.Fatal("AWS_WEB_IDENTITY_TOKEN_FILE not set")
	}

	defaultHttpClient := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		log.Fatal("AWS_REGION not set")
	}

	baseURL := fmt.Sprintf("https://sts.%s.amazonaws.com", awsRegion)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL,
		http.NoBody,
	)
	if err != nil {
		log.Fatal(err)
	}

	awsRoleArn, ok := os.LookupEnv("AWS_ROLE_ARN")
	if !ok {
		log.Fatal("AWS_ROLE_ARN not set")
	}

	awsWebIdentityToken, ok := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if !ok {
		log.Fatal("AWS_WEB_IDENTITY_TOKEN_FILE not set")
	}

	q := url.Values{}
	q.Add("api_key", "key_from_environment_or_flag")
	q.Add("RoleArn", awsRoleArn)
	q.Add("WebIdentityToken", awsWebIdentityToken)
	q.Add("RoleSessionName", "app1")
	q.Add("Version", "2011-06-15")

	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	resp, err := defaultHttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var awsResp Response
	if err := json.NewDecoder(resp.Body).Decode(&awsResp); err != nil {
		log.Printf("Could not decode JSON Response: %s", err)
	}

	crd := awsResp.Response.Result.Creds

	fmt.Printf("{ \"Version\": 1, "+
		"\"AccessKeyId\": \"%s\", "+
		"\"SecretAccessKey\": \"%s\", "+
		"\"SessionToken\": \"%s\", ",
		crd.AccessKeyId,
		crd.SecretAccessKey,
		crd.SessionToken)
}
