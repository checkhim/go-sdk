# Security Policy

## Supported Versions

We actively support the following versions of the CheckHim Go SDK:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | Yes             |
| 0.x.x   | No (Beta)       |

## Reporting a Vulnerability

The CheckHim team takes security vulnerabilities seriously. We appreciate your efforts to responsibly disclose any issues you may find.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please send an email to: **security@checkhim.tech**

### What to Include

Please include the following information in your report:

- **Description**: A clear description of the vulnerability
- **Impact**: Potential impact and attack scenarios
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Environment**: Go version, SDK version, and operating system
- **Proof of Concept**: Code examples or screenshots (if applicable)

### Response Timeline

- **Acknowledgment**: We will acknowledge receipt within 24 hours
- **Initial Assessment**: We will provide an initial assessment within 72 hours
- **Updates**: We will provide regular updates every 7 days
- **Resolution**: We aim to resolve critical issues within 30 days

### Security Best Practices

When using the CheckHim Go SDK:

#### API Key Security

```go
// Good: Use environment variables
apiKey := os.Getenv("CHECKHIM_API_KEY")
client := checkhim.New(apiKey)

// Bad: Don't hardcode API keys
client := checkhim.New("your-actual-api-key-here")
```

#### Timeout Configuration

```go
// Good: Set reasonable timeouts
client := checkhim.New(apiKey, checkhim.Config{
    Timeout: 30 * time.Second,
})

// Good: Use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
result, err := client.VerifyWithContext(ctx, req)
```

#### Input Validation

```go
// Good: Validate input before sending
func validatePhoneNumber(number string) error {
    if number == "" {
        return errors.New("phone number cannot be empty")
    }
    if !strings.HasPrefix(number, "+") {
        return errors.New("phone number must include country code")
    }
    return nil
}
```

#### Error Handling

```go
// Good: Handle API errors appropriately
result, err := client.Verify(req)
if err != nil {
    if apiErr, ok := err.(*checkhim.APIError); ok {
        // Log error without exposing sensitive data
        log.Printf("API error: status=%d, code=%s", apiErr.StatusCode, apiErr.Code)
        return
    }
    log.Printf("Network error: %v", err)
    return
}
```

### Security Considerations

#### Data Privacy

- Phone numbers are considered sensitive personal data
- Ensure compliance with local privacy regulations (GDPR, CCPA, etc.)
- Implement appropriate data handling and retention policies
- Consider data minimization principles

#### Network Security

- Always use HTTPS endpoints (default in the SDK)
- Consider using VPN or private networks for sensitive operations
- Implement retry mechanisms with exponential backoff
- Monitor for unusual API usage patterns

#### Authentication

- Rotate API keys regularly
- Use different API keys for different environments
- Implement proper key management practices
- Monitor API key usage and set up alerts for anomalies

### Vulnerability Disclosure Policy

We follow responsible disclosure practices:

1. **Private Disclosure**: Security issues are first reported privately
2. **Investigation**: We investigate and develop fixes
3. **Coordination**: We coordinate with reporters on disclosure timeline
4. **Public Disclosure**: Details are made public after fixes are released
5. **Credit**: We provide credit to security researchers (with permission)

### Hall of Fame

We recognize security researchers who help improve our security:

<!-- Security researchers will be listed here -->

### Contact Information

- **Security Team**: security@checkhim.tech
- **General Support**: support@checkhim.tech
- **Website**: https://checkhim.tech

### Legal

This security policy is subject to our [Terms of Service](https://checkhim.tech/terms) and [Privacy Policy](https://checkhim.tech/privacy).

---

Thank you for helping keep CheckHim and the Go community secure! 🛡️
