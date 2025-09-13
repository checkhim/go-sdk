# Contributing to CheckHim Go SDK

Thank you for your interest in contributing to the CheckHim Go SDK! We welcome contributions from developers around the world.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to conduct@checkhim.tech.

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. Check if the issue already exists in our [issue tracker](https://github.com/checkhim/go-sdk/issues)
2. Make sure you're using the latest version of the SDK
3. Provide a clear and detailed description of the problem
4. Include relevant code examples and error messages

### Suggesting Features

We welcome feature suggestions! Please:

1. Check if the feature has already been requested
2. Clearly describe the feature and its use case
3. Explain why it would be valuable to other users
4. Consider providing a proof of concept or example implementation

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Update documentation** as necessary
5. **Ensure tests pass** by running `go test ./...`
6. **Create a pull request** with a clear title and description

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Setup

```bash
# Fork and clone the repository
git clone https://github.com/your-username/go-sdk.git
cd go-sdk

# Install dependencies
go mod tidy

# Run tests to make sure everything works
go test ./...
```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Write clear, self-documenting code with appropriate comments

### Code Organization

- Keep functions focused and small
- Use meaningful variable and function names
- Group related functionality together
- Follow Go package naming conventions

### Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Add benchmarks for performance-critical code

Example test structure:

```go
func TestClient_Verify(t *testing.T) {
    t.Run("successful verification", func(t *testing.T) {
        // Test implementation
    })
    
    t.Run("error case", func(t *testing.T) {
        // Test implementation
    })
}
```

### Documentation

- Add godoc comments for all public functions and types
- Include usage examples in comments
- Update README.md for significant changes
- Add examples to the examples directory when appropriate

Example documentation:

```go
// Verify verifies a phone number using the CheckHim API.
// It returns detailed information about the phone number including
// validity, country, carrier, and line type.
//
// Example usage:
//   result, err := client.Verify(checkhim.VerifyRequest{
//       Number: "+1234567890",
//   })
func (c *Client) Verify(req VerifyRequest) (*VerifyResponse, error) {
    // Implementation
}
```

## Commit Guidelines

### Commit Message Format

Use the following format for commit messages:

```
type(scope): subject

body

footer
```

### Types

- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(client): add context support for verification requests

Add VerifyWithContext method to allow cancellation and timeouts
for verification requests. This improves reliability in production
environments where request timeouts are critical.

Closes #123
```

```
fix(client): handle malformed JSON responses gracefully

Previously, malformed JSON responses would cause panics. Now they
are handled gracefully with appropriate error messages.

Fixes #456
```

## Testing Guidelines

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem
```

### Test Types

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test interaction with the API (using mock servers)
3. **Example Tests**: Demonstrate usage patterns
4. **Benchmark Tests**: Measure performance

### Mock Servers

Use `httptest.NewServer` for testing HTTP interactions:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Mock response
    json.NewEncoder(w).Encode(VerifyResponse{Valid: true})
}))
defer server.Close()

client := New("test-key", Config{BaseURL: server.URL})
```

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- `MAJOR`: Incompatible API changes
- `MINOR`: New functionality (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Release Checklist

- [ ] Update version in relevant files
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Update documentation
- [ ] Create GitHub release with release notes

## Getting Help

If you need help with development:

1. Check existing issues and discussions
2. Read the documentation thoroughly
3. Join our Discord server: [link]
4. Email the maintainers: maintainers@checkhim.tech

## Recognition

Contributors will be recognized in:

- CONTRIBUTORS.md file
- GitHub contributors page
- Release notes (for significant contributions)

Thank you for contributing to CheckHim Go SDK! 🎉
