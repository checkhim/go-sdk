# CheckHim Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/checkhim/go-sdk.svg)](https://pkg.go.dev/github.com/checkhim/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/checkhim/go-sdk)](https://goreportcard.com/report/github.com/checkhim/go-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/checkhim/go-sdk/badge.svg?branch=main)](https://coveralls.io/github/checkhim/go-sdk?branch=main)

The official Go SDK for the CheckHim phone number verification API. Verify phone numbers quickly and reliably with our global verification service.

## Features

- **Simple API** - Easy-to-use interface with minimal setup
- **Global Coverage** - Verify phone numbers from around the world  
- **High Performance** - Optimized for speed and reliability
- **Secure** - Built-in security best practices
- **Error Handling** - Comprehensive error handling and recovery
- **Rich Response Data** - Get detailed information about verified numbers
- **Context Support** - Full context.Context support for cancellation and timeouts
- **Well Tested** - Extensive test coverage with examples

## Installation

```bash
go get github.com/checkhim/go-sdk
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/checkhim/go-sdk"
)

func main() {
    // Create a new client with your API key
    client := checkhim.New("your-api-key-here")
    
    // Verify a phone number
    result, err := client.Verify(checkhim.VerifyRequest{
        Number: "+5511984339000",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the results
    fmt.Printf("Valid: %v, Carrier: %s\n", result.Valid, result.Carrier)
}
```

## Usage Examples

### Basic Verification

```go
client := checkhim.New("your-api-key")

result, err := client.Verify(checkhim.VerifyRequest{
    Number: "+1234567890",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Valid: %v\n", result.Valid)
fmt.Printf("Carrier: %s\n", result.Carrier)
```

### With Context and Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.VerifyWithContext(ctx, checkhim.VerifyRequest{
    Number: "+5511984339000",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Valid: %v, Carrier: %s\n", result.Valid, result.Carrier)
```

### Custom Configuration

```go
client := checkhim.New("your-api-key", checkhim.Config{
    BaseURL: "https://api.checkhim.tech", // Custom API endpoint
    Timeout: 30 * time.Second,            // Custom timeout
})

result, err := client.Verify(checkhim.VerifyRequest{
    Number: "+244921204020",
})
```

### Error Handling

```go
result, err := client.Verify(checkhim.VerifyRequest{
    Number: "+invalid-number",
})
if err != nil {
    // Check if it's an API error
    if apiErr, ok := err.(*checkhim.APIError); ok {
        fmt.Printf("API Error - Status: %d, Message: %s, Code: %s\n",
            apiErr.StatusCode, apiErr.Message, apiErr.Code)
        
        // Handle specific error codes
        switch apiErr.Code {
        case "unauthorized":
            log.Fatal("Invalid API key")
        case "rate_limit_exceeded":
            log.Fatal("Rate limit exceeded")
        default:
            log.Printf("API error: %v", apiErr)
        }
    } else {
        log.Printf("Network error: %v", err)
    }
    return
}

// Process successful result
fmt.Printf("Valid: %v\n", result.Valid)
```

### Batch Verification

```go
phoneNumbers := []string{
    "+1234567890",
    "+5511984339000", 
    "+244921204020",
}

for _, number := range phoneNumbers {
    result, err := client.Verify(checkhim.VerifyRequest{Number: number})
    if err != nil {
        log.Printf("Error verifying %s: %v", number, err)
        continue
    }
    
    fmt.Printf("%s: Valid=%v, Carrier=%s\n", 
        number, result.Valid, result.Carrier)
}
```

## API Reference

### Client

#### `New(apiKey string, configs ...Config) *Client`

Creates a new CheckHim client.

**Parameters:**
- `apiKey` - Your CheckHim API key (required)
- `configs` - Optional configuration (see Config struct)

#### `Verify(req VerifyRequest) (*VerifyResponse, error)`

Verifies a phone number.

#### `VerifyWithContext(ctx context.Context, req VerifyRequest) (*VerifyResponse, error)`

Verifies a phone number with context support.

### Types

#### `VerifyRequest`

```go
type VerifyRequest struct {
    Number string `json:"number"` // Phone number with country code (e.g., "+1234567890")
}
```

#### `VerifyResponse`

```go
type VerifyResponse struct {
    Carrier string `json:"carrier"` // Mobile carrier name (e.g., "UNITEL")
    Valid   bool   `json:"valid"`   // Whether the number is valid
}
```

#### `Config`

```go
type Config struct {
    BaseURL    string        // Custom API base URL
    Timeout    time.Duration // HTTP request timeout
    HTTPClient *http.Client  // Custom HTTP client
}
```

#### `APIError`

```go
type APIError struct {
    StatusCode int                    // HTTP status code
    Message    string                 // Error message
    Code       string                 // Error code
    Details    map[string]interface{} // Additional error details
}
```

## Error Codes

Common error codes returned by the API:

| Code | Description |
|------|-------------|
| `unauthorized` | Invalid or missing API key |
| `invalid_request` | Malformed request or missing required fields |
| `rate_limit_exceeded` | Too many requests, please slow down |
| `invalid_number` | Phone number format is invalid |
| `insufficient_credits` | Account has insufficient credits |
| `service_unavailable` | Temporary service unavailability |

## Configuration

### Environment Variables

You can set default configuration using environment variables:

```bash
export CHECKHIM_API_KEY="your-api-key"
export CHECKHIM_BASE_URL="https://api.checkhim.tech"
export CHECKHIM_TIMEOUT="30s"
```

### Custom HTTP Client

```go
customHTTPClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := checkhim.New("your-api-key", checkhim.Config{
    HTTPClient: customHTTPClient,
})
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

## Examples

Check out the [examples](examples/) directory for more comprehensive usage examples:

- [Basic Usage](examples/main.go) - Complete example with various use cases

To run the examples:

```bash
cd examples
go mod tidy
go run main.go
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/go-sdk.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes and add tests
5. Run tests: `go test ./...`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Add tests for new functionality
- Update documentation as needed

## Support

- Email: support@checkhim.tech
- Website: [https://checkhim.tech](https://checkhim.tech)
- Documentation: [https://docs.checkhim.tech](https://docs.checkhim.tech)
- Issues: [GitHub Issues](https://github.com/checkhim/go-sdk/issues)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and updates.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

If you discover a security vulnerability, please send an e-mail to security@checkhim.tech. All security vulnerabilities will be promptly addressed.

---

Made with ❤️ by the CheckHim team
