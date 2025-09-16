# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of CheckHim Go SDK
- Phone number verification functionality
- Context support for request cancellation and timeouts
- Comprehensive error handling with custom error types
- Configurable HTTP client and timeouts
- Rich response data including country, carrier, and line type information
- Full test coverage with examples and benchmarks

### Features
- `checkhim.New()` - Create new client with API key
- `client.Verify()` - Verify phone numbers
- `client.VerifyWithContext()` - Verify with context support
- Custom configuration options via `Config` struct
- Detailed `VerifyResponse` with validation results
- `APIError` type for structured error handling

## [1.0.0] - 2025-09-01

### Added
- Initial stable release
- Full API compatibility with CheckHim verification service
- Production-ready SDK with comprehensive documentation
- Examples and usage guides
- Contributing guidelines for open source contributors

### Security
- Secure API key handling
- HTTPS by default
- No sensitive data logging

---

## Release Types

- **Added** for new features
- **Changed** for changes in existing functionality  
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** for vulnerability fixes
