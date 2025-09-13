//go:build integration
// +build integration

package checkhim

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_Verify tests the SDK against the real CheckHim API
// Run with: go test -tags=integration -v
// Requires CHECKHIM_API_KEY environment variable to be set
func TestIntegration_Verify(t *testing.T) {
	apiKey := os.Getenv("CHECKHIM_API_KEY")
	if apiKey == "" {
		t.Skip("CHECKHIM_API_KEY not set, skipping integration tests")
	}

	client := New(apiKey)

	t.Run("verify valid phone number", func(t *testing.T) {
		result, err := client.Verify(VerifyRequest{
			Number: "+244921204020", // Test number from the example
		})

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Log the result for debugging
		t.Logf("Result: Valid=%v, Carrier=%s", result.Valid, result.Carrier)
	})

	t.Run("verify with context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.VerifyWithContext(ctx, VerifyRequest{
			Number: "+5511984339000",
		})

		require.NoError(t, err)
		assert.NotNil(t, result)

		t.Logf("Result: Valid=%v, Carrier=%s", result.Valid, result.Carrier)
	})

	t.Run("verify invalid number format", func(t *testing.T) {
		result, err := client.Verify(VerifyRequest{
			Number: "invalid-number",
		})

		// This might return an error or a result with Valid=false
		// depending on the API implementation
		if err != nil {
			t.Logf("Error (expected): %v", err)
		} else {
			t.Logf("Result for invalid number: Valid=%v", result.Valid)
			assert.False(t, result.Valid)
		}
	})
}

// TestIntegration_ErrorHandling tests error scenarios with the real API
func TestIntegration_ErrorHandling(t *testing.T) {
	t.Run("invalid API key", func(t *testing.T) {
		client := New("invalid-api-key")

		result, err := client.Verify(VerifyRequest{
			Number: "+1234567890",
		})

		require.Error(t, err)
		assert.Nil(t, result)

		// Check if it's an API error
		if apiErr, ok := err.(*APIError); ok {
			assert.Equal(t, 401, apiErr.StatusCode)
			t.Logf("API Error: %v", apiErr)
		}
	})
}

// BenchmarkIntegration_Verify benchmarks the real API performance
func BenchmarkIntegration_Verify(b *testing.B) {
	apiKey := os.Getenv("CHECKHIM_API_KEY")
	if apiKey == "" {
		b.Skip("CHECKHIM_API_KEY not set, skipping integration benchmarks")
	}

	client := New(apiKey)
	req := VerifyRequest{Number: "+244921204020"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Verify(req)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
