package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/checkhim/go-sdk"
)

func main() {
	// Create a new CheckHim client with your API key
	client := checkhim.New("your-api-key-here")

	// Example 1: Basic phone number verification
	fmt.Println("=== Basic Verification ===")
	result, err := client.Verify(checkhim.VerifyRequest{
		Number: "+5511984339000",
	})
	if err != nil {
		log.Printf("Error verifying phone number: %v", err)
	} else {
		fmt.Printf("Valid: %v\n", result.Valid)
		fmt.Printf("Carrier: %s\n", result.Carrier)
	}

	fmt.Println()

	// Example 2: Verification with context and timeout
	fmt.Println("=== Verification with Context ===")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result2, err := client.VerifyWithContext(ctx, checkhim.VerifyRequest{
		Number: "+1234567890",
	})
	if err != nil {
		log.Printf("Error verifying phone number with context: %v", err)
	} else {
		fmt.Printf("Valid: %v, Carrier: %s\n", result2.Valid, result2.Carrier)
	}

	fmt.Println()

	// Example 3: Custom configuration
	fmt.Println("=== Custom Configuration ===")
	customClient := checkhim.New("your-api-key-here", checkhim.Config{
		BaseURL: "https://api.checkhim.tech", // Custom base URL
		Timeout: 30 * time.Second,            // Custom timeout
	})

	result3, err := customClient.Verify(checkhim.VerifyRequest{
		Number: "+244921204020",
	})
	if err != nil {
		log.Printf("Error with custom client: %v", err)
	} else {
		fmt.Printf("Valid: %v, Carrier: %s\n", result3.Valid, result3.Carrier)
	}

	fmt.Println()

	// Example 4: Error handling
	fmt.Println("=== Error Handling ===")
	_, err = client.Verify(checkhim.VerifyRequest{
		Number: "", // Empty number will cause an error
	})
	if err != nil {
		// Check if it's an API error
		if apiErr, ok := err.(*checkhim.APIError); ok {
			fmt.Printf("API Error - Status: %d, Message: %s, Code: %s\n",
				apiErr.StatusCode, apiErr.Message, apiErr.Code)
		} else {
			fmt.Printf("Other error: %v\n", err)
		}
	}

	fmt.Println()

	// Example 5: Batch verification
	fmt.Println("=== Batch Verification ===")
	phoneNumbers := []string{
		"+1234567890",
		"+5511984339000",
		"+244921204020",
		"+4912345678",
	}

	for _, number := range phoneNumbers {
		result, err := client.Verify(checkhim.VerifyRequest{Number: number})
		if err != nil {
			fmt.Printf("%-15s: Error - %v\n", number, err)
		} else {
			fmt.Printf("%-15s: Valid=%v, Carrier=%s\n", number, result.Valid, result.Carrier)
		}
	}
}
