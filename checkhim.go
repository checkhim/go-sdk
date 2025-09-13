// Package checkhim provides a Go SDK for phone number verification using the CheckHim API.
//
// This SDK allows you to verify phone numbers by checking if they are valid and active.
// It provides a simple interface to interact with the CheckHim verification service.
//
// Basic usage:
//
//	client := checkhim.New("your-api-key")
//	result, err := client.Verify(checkhim.VerifyRequest{
//		Number: "+1234567890",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Valid: %v, Carrier: %s\n", result.Valid, result.Carrier)
package checkhim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the CheckHim API
	DefaultBaseURL = "http://api.checkhim.tech"

	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second

	// APIVersion is the current API version
	APIVersion = "v1"
)

// Client represents a CheckHim API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Config holds configuration options for the Client
type Config struct {
	// BaseURL is the base URL for the CheckHim API (optional)
	BaseURL string

	// Timeout is the timeout for HTTP requests (optional)
	Timeout time.Duration

	// HTTPClient is a custom HTTP client (optional)
	HTTPClient *http.Client
}

// New creates a new CheckHim client with the provided API key
func New(apiKey string, configs ...Config) *Client {
	config := Config{
		BaseURL: DefaultBaseURL,
		Timeout: DefaultTimeout,
	}

	if len(configs) > 0 {
		if configs[0].BaseURL != "" {
			config.BaseURL = configs[0].BaseURL
		}
		if configs[0].Timeout > 0 {
			config.Timeout = configs[0].Timeout
		}
		if configs[0].HTTPClient != nil {
			config.HTTPClient = configs[0].HTTPClient
		}
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    config.BaseURL,
		httpClient: httpClient,
	}
}

// VerifyRequest represents a phone number verification request
type VerifyRequest struct {
	// Number is the phone number to verify (required)
	// Should include country code (e.g., "+1234567890")
	Number string `json:"number"`
}

// internalVerifyRequest is the internal request structure sent to the API
type internalVerifyRequest struct {
	// Number is the phone number to verify
	Number string `json:"number"`

	// Type is always "frontend" for this SDK
	Type string `json:"type"`
}

// VerifyResponse represents the response from a phone number verification
type VerifyResponse struct {
	// Carrier is the name of the mobile carrier
	Carrier string `json:"carrier"`

	// Valid indicates whether the phone number is valid and active
	Valid bool `json:"valid"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	// Error is the error message
	Error string `json:"error"`

	// Code is the error code
	Code string `json:"code,omitempty"`

	// Details provides additional error details
	Details map[string]interface{} `json:"details,omitempty"`
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Details    map[string]interface{}
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("checkhim: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("checkhim: %s (status: %d)", e.Message, e.StatusCode)
}

// Verify verifies a phone number using the CheckHim API
func (c *Client) Verify(req VerifyRequest) (*VerifyResponse, error) {
	return c.VerifyWithContext(context.Background(), req)
}

// VerifyWithContext verifies a phone number with a custom context
func (c *Client) VerifyWithContext(ctx context.Context, req VerifyRequest) (*VerifyResponse, error) {
	if req.Number == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "phone number is required",
			Code:       "invalid_request",
		}
	}

	internalReq := internalVerifyRequest{
		Number: req.Number,
		Type:   "frontend",
	}

	reqBody, err := json.Marshal(internalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/verify", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "checkhim-go-sdk/1.0")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    errorResp.Error,
				Code:       errorResp.Code,
				Details:    errorResp.Details,
			}
		}

		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	var verifyResp VerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &verifyResp, nil
}
