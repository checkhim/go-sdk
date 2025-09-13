// Package checkhim provides a Go SDK for phone number verification using the CheckHim API.
//
// This SDK allows you to verify phone numbers by checking if they are valid and active.
// It provides a simple interface to interact with the CheckHim verification service.
//
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
	DefaultBaseURL = "http://api.checkhim.tech"
	
	DefaultTimeout = 30 * time.Second
	
	APIVersion = "v1"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Config struct {
	BaseURL string
	
	Timeout time.Duration
	
	HTTPClient *http.Client
}

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

type VerifyRequest struct {
	Number string `json:"number"`
}

type internalVerifyRequest struct {
	Number string `json:"number"`
	
	Type string `json:"type"`
}

type VerifyResponse struct {
	Carrier string `json:"carrier"`
	
	Valid bool `json:"valid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	
	Code string `json:"code,omitempty"`
	
	Details map[string]interface{} `json:"details,omitempty"`
}

type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Details    map[string]interface{}
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("checkhim: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("checkhim: %s (status: %d)", e.Message, e.StatusCode)
}

func (c *Client) Verify(req VerifyRequest) (*VerifyResponse, error) {
	return c.VerifyWithContext(context.Background(), req)
}

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
