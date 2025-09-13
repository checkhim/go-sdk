package checkhim

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates client with default config", func(t *testing.T) {
		client := New("test-api-key")
		
		assert.Equal(t, "test-api-key", client.apiKey)
		assert.Equal(t, DefaultBaseURL, client.baseURL)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, DefaultTimeout, client.httpClient.Timeout)
	})
	
	t.Run("creates client with custom config", func(t *testing.T) {
		customTimeout := 60 * time.Second
		customBaseURL := "https://custom.api.com"
		customHTTPClient := &http.Client{Timeout: 120 * time.Second}
		
		client := New("test-api-key", Config{
			BaseURL:    customBaseURL,
			Timeout:    customTimeout,
			HTTPClient: customHTTPClient,
		})
		
		assert.Equal(t, "test-api-key", client.apiKey)
		assert.Equal(t, customBaseURL, client.baseURL)
		assert.Equal(t, customHTTPClient, client.httpClient)
	})
}

func TestClient_Verify(t *testing.T) {
	t.Run("successful verification", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method and path
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/verify", r.URL.Path)
			
			// Verify headers
			assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Contains(t, r.Header.Get("User-Agent"), "checkhim-go-sdk")
			
			// Verify request body - check internal structure
			var internalReq struct {
				Number string `json:"number"`
				Type   string `json:"type"`
			}
			err := json.NewDecoder(r.Body).Decode(&internalReq)
			require.NoError(t, err)
			assert.Equal(t, "+1234567890", internalReq.Number)
			assert.Equal(t, "frontend", internalReq.Type)
			
			// Send successful response
			response := VerifyResponse{
				Carrier: "UNITEL",
				Valid:   true,
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()
		
		client := New("test-api-key", Config{BaseURL: server.URL})
		
		result, err := client.Verify(VerifyRequest{Number: "+1234567890"})
		
		require.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Equal(t, "UNITEL", result.Carrier)
	})
	
	t.Run("invalid phone number", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := VerifyResponse{
				Carrier: "",
				Valid:   false,
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()
		
		client := New("test-api-key", Config{BaseURL: server.URL})
		
		result, err := client.Verify(VerifyRequest{Number: "+invalid"})
		
		require.NoError(t, err)
		assert.False(t, result.Valid)
		assert.Equal(t, "", result.Carrier)
	})
	
	t.Run("empty phone number", func(t *testing.T) {
		client := New("test-api-key")
		
		result, err := client.Verify(VerifyRequest{Number: ""})
		
		require.Error(t, err)
		assert.Nil(t, result)
		
		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, 400, apiErr.StatusCode)
		assert.Contains(t, apiErr.Message, "phone number is required")
		assert.Equal(t, "invalid_request", apiErr.Code)
	})
	
	t.Run("API error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			response := ErrorResponse{
				Error: "Invalid API key",
				Code:  "unauthorized",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()
		
		client := New("invalid-key", Config{BaseURL: server.URL})
		
		result, err := client.Verify(VerifyRequest{Number: "+1234567890"})
		
		require.Error(t, err)
		assert.Nil(t, result)
		
		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, 401, apiErr.StatusCode)
		assert.Equal(t, "Invalid API key", apiErr.Message)
		assert.Equal(t, "unauthorized", apiErr.Code)
	})
	
	t.Run("network error", func(t *testing.T) {
		client := New("test-api-key", Config{BaseURL: "http://invalid-url-that-does-not-exist.local"})
		
		result, err := client.Verify(VerifyRequest{Number: "+1234567890"})
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to execute request")
	})
	
	t.Run("malformed JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()
		
		client := New("test-api-key", Config{BaseURL: server.URL})
		
		result, err := client.Verify(VerifyRequest{Number: "+1234567890"})
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})
}

func TestClient_VerifyWithContext(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		
		client := New("test-api-key", Config{BaseURL: server.URL})
		
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		
		result, err := client.VerifyWithContext(ctx, VerifyRequest{Number: "+1234567890"})
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}

func TestAPIError_Error(t *testing.T) {
	t.Run("error with code", func(t *testing.T) {
		err := &APIError{
			StatusCode: 400,
			Message:    "Bad request",
			Code:       "bad_request",
		}
		
		expected := "checkhim: Bad request (code: bad_request, status: 400)"
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error without code", func(t *testing.T) {
		err := &APIError{
			StatusCode: 500,
			Message:    "Internal server error",
		}
		
		expected := "checkhim: Internal server error (status: 500)"
		assert.Equal(t, expected, err.Error())
	})
}

// Example tests that demonstrate usage
func ExampleNew() {
	client := New("your-api-key")
	
	result, err := client.Verify(VerifyRequest{
		Number: "+1234567890",
	})
	if err != nil {
		panic(err)
	}
	
	println(result.Valid, result.Carrier)
}

func ExampleClient_Verify() {
	client := New("your-api-key")
	
	result, err := client.Verify(VerifyRequest{
		Number: "+5511984339000",
	})
	if err != nil {
		panic(err)
	}
	
	println(result.Valid, result.Carrier)
}

func ExampleClient_VerifyWithContext() {
	client := New("your-api-key")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := client.VerifyWithContext(ctx, VerifyRequest{
		Number: "+1234567890",
	})
	if err != nil {
		panic(err)
	}
	
	println(result.Valid, result.Carrier)
}

// Benchmark tests
func BenchmarkClient_Verify(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := VerifyResponse{
			Carrier: "UNITEL",
			Valid:   true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := New("test-api-key", Config{BaseURL: server.URL})
	req := VerifyRequest{Number: "+1234567890"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Verify(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
