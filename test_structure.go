package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/checkhim/go-sdk"
)

func main() {
	client := checkhim.New("test-api-key")

	// Test request structure
	req := checkhim.VerifyRequest{
		Number: "+5511984339000",
	}

	// Print what would be sent (for debugging)
	reqJSON, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Request would be:\n%s\n", string(reqJSON))

	// Test response structure
	resp := checkhim.VerifyResponse{
		Carrier: "UNITEL",
		Valid:   true,
	}

	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nExpected response format:\n%s\n", string(respJSON))

	fmt.Printf("\nUsage example:\n")
	fmt.Printf("client := checkhim.New(\"your-api-key\")\n")
	fmt.Printf("result, err := client.Verify(checkhim.VerifyRequest{Number: \"+5511984339000\"})\n")
	fmt.Printf("if err != nil { panic(err) }\n")
	fmt.Printf("println(result.Valid, result.Carrier)\n")
}
