# Contributing to Glance

First off, thank you for considering contributing to Glance! It's people like you that make Glance such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include screenshots and animated GIFs if possible**
* **Include your environment details** (OS, Go version, Browser, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and expected behavior**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Follow the Go and JavaScript style guides
* End all files with a newline
* Avoid platform-dependent code
* Include appropriate test cases
* Update relevant documentation

## Development Setup

### Prerequisites
- Go 1.19+
- Node.js 18+
- SQLite 3
- Git

### Local Setup

```bash
# Clone the repository
git clone https://github.com/glance-project/glance.git
cd glance

# Create a feature branch
git checkout -b feature/your-feature-name

# Install dependencies
go mod download
npm install

# Initialize database
sqlite3 glance.db < scripts/init.sql

# Run tests
go test ./...
npm test

# Start development server
go run ./cmd/main.go
```

## Style Guides

### Go Code Style

* Follow the [Effective Go](https://golang.org/doc/effective_go) guide
* Use `gofmt` for formatting
* Use `golint` for linting
* Variable names should be camelCase
* Exported identifiers should be PascalCase
* Write tests for all public functions

```go
package api

// Handler handles API requests
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### JavaScript Code Style

* Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
* Use ES6 modules
* Use camelCase for variables and functions
* Use PascalCase for classes
* Include JSDoc comments for public methods

```javascript
/**
 * Initialize the widget
 * @param {Element} container - The container element
 * @param {Object} config - Configuration options
 */
export class Widget {
  constructor(container, config = {}) {
    this.container = container;
    this.config = config;
  }

  initialize() {
    // Implementation
  }
}
```

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

```
Add rate limiting middleware

- Implement token bucket algorithm
- Add middleware to protect API endpoints
- Include comprehensive tests
- Update documentation

Fixes #123
```

### Documentation

* Use Markdown for documentation
* Include code examples where appropriate
* Update relevant documentation when making changes
* Keep README.md up to date

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with verbose output
go test -v ./...
```

### Writing Tests

* Name test functions as `TestFunctionName`
* Use table-driven tests for multiple cases
* Include comments explaining complex test logic
* Aim for >80% code coverage

```go
func TestRateLimiter(t *testing.T) {
    tests := []struct {
        name    string
        input   int
        want    bool
    }{
        {"Allow valid", 5, true},
        {"Deny over limit", 65, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := rateLimiter.Allow(tt.input); got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Build and Deployment

### Building

```bash
# Development build
go build -o bin/glance ./cmd/main.go

# Production build
go build -ldflags="-s -w" -o bin/glance-prod ./cmd/main.go
```

### Linting

```bash
# Go
golint ./...
go vet ./...

# JavaScript
npm run lint
```

## Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Write** or update tests
5. **Run** linting and tests
6. **Commit** your changes with clear messages
7. **Push** to your fork
8. **Create** a Pull Request

## Approval Process

Pull requests are reviewed by maintainers. We look for:
* Code follows style guidelines
* Tests are included and passing
* Documentation is updated
* No unrelated changes
* Clear commit history

## Additional Notes

### Issue and Pull Request Labels

* `bug` - Something isn't working
* `enhancement` - New feature or request
* `documentation` - Improvements or additions to documentation
* `good first issue` - Good for newcomers
* `help wanted` - Extra attention is needed
* `question` - Further information is requested

## Recognition

Contributors will be recognized in:
* Git commit history
* GitHub contributors page
* Release notes
* Project documentation

## Questions?

Feel free to create a GitHub discussion or reach out to the maintainers.

## License

By contributing to Glance, you agree that your contributions will be licensed under its MIT License.

## Acknowledgments

This contributing guide was inspired by the [Atom Contributing Guide](https://github.com/atom/atom/blob/master/CONTRIBUTING.md).
