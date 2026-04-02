# Contributing to TaskFix

First off, thank you for considering contributing to TaskFix! It's people like you that make TaskFix such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

---

## Code of Conduct

This project and everyone participating in it is governed by respect, professionalism, and inclusivity. By participating, you are expected to uphold this standard. Please be respectful and constructive in all interactions.

### Our Standards

**Examples of behavior that contributes to a positive environment:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Examples of unacceptable behavior:**
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate

---

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the [issue tracker](https://github.com/abdomassoun/taskfix/issues) to avoid duplicates. When creating a bug report, include as many details as possible using the bug report template.

**Good bug reports include:**
- A clear, descriptive title
- Exact steps to reproduce the problem
- The behavior you observed and what you expected
- Screenshots if applicable
- Your environment (OS, Go version, TaskFix version)
- Any relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the proposed feature
- Explain why this enhancement would be useful
- List any similar features in other tools (if applicable)
- Include mockups or examples if relevant

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `documentation` - Documentation improvements

### Pull Requests

We actively welcome your pull requests! Here's how to submit one:

1. Fork the repo and create your branch from `main`
2. Make your changes and ensure tests pass
3. Update documentation if needed
4. Write a clear commit message
5. Submit the pull request

---

## Development Setup

### Prerequisites

Make sure you have the following installed:
- **Go**: Version 1.21 or higher ([install guide](https://golang.org/doc/install))
- **Git**: For version control
- **Make**: For running build commands (optional but recommended)

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/taskfix.git
   cd taskfix
   ```

3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/abdomassoun/taskfix.git
   ```

4. **Install dependencies:**
   ```bash
   go mod download
   ```

5. **Build the project:**
   ```bash
   make build
   # Or manually:
   go build -o taskfix .
   ```

6. **Set up your API key for testing:**
   ```bash
   cat > ~/.tfixrc << 'EOF'
   {
     "api_key": "sk-or-v1-your-test-api-key",
     "model": "openai/gpt-4o-mini"
   }
   EOF
   chmod 600 ~/.tfixrc
   ```

7. **Run tests:**
   ```bash
   make test
   # Or manually:
   go test ./...
   ```

### Verify Installation

```bash
./taskfix "test bug description"
```

If everything is set up correctly, you should see formatted output.

---

## Development Workflow

### Creating a New Branch

Always create a new branch for your work:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

**Branch naming conventions:**
- `feature/feature-name` - New features
- `fix/bug-description` - Bug fixes
- `docs/what-changed` - Documentation updates
- `refactor/what-changed` - Code refactoring
- `test/what-added` - Test additions

### Making Changes

1. **Write clean, readable code** following Go best practices
2. **Add tests** for new functionality
3. **Update documentation** if you change behavior
4. **Run tests frequently** to catch issues early

### Testing Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/ai/

# Run tests verbosely
go test -v ./...

# Run linter
make lint
# or
go vet ./...
```

### Keeping Your Fork Updated

Regularly sync your fork with the upstream repository:

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

---

## Coding Standards

### Go Style Guidelines

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guide
- Use `gofmt` to format your code (run automatically with `make fmt`)
- Follow Go naming conventions (exported vs unexported)
- Keep functions small and focused (single responsibility)
- Add comments for exported functions and complex logic

### Code Organization

```
internal/
├── ai/         # AI provider integrations
├── rules/      # Rule processing logic
├── prompt/     # Prompt building
└── output/     # Output formatting

cmd/            # CLI commands and flags
configs/        # Configuration presets
```

### Error Handling

```go
// Good: Return errors, don't panic
func ProcessTask(input string) (string, error) {
    if input == "" {
        return "", fmt.Errorf("input cannot be empty")
    }
    // ...
}

// Bad: Panic on errors
func ProcessTask(input string) string {
    if input == "" {
        panic("input cannot be empty")
    }
    // ...
}
```

### Comments

```go
// Good: Explain why, not what
// Cache the API client to avoid rate limiting
var cachedClient *Client

// Bad: State the obvious
// This variable holds the client
var cachedClient *Client
```

### Testing

- Write table-driven tests when testing multiple scenarios
- Use meaningful test names that describe the scenario
- Test both success and failure cases
- Mock external dependencies (API calls, file I/O)

Example:
```go
func TestFormatTask(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "bug in login",
            expected: "Title: Fix login bug",
            wantErr:  false,
        },
        {
            name:     "empty input",
            input:    "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FormatTask(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FormatTask() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FormatTask() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

---

## Commit Guidelines

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semi-colons, etc.)
- `refactor`: Code refactoring without changing functionality
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates, etc.

**Examples:**

```bash
feat(ai): add support for Anthropic Claude models

- Add Claude API client implementation
- Update config to support provider selection
- Add tests for Claude integration

Closes #123
```

```bash
fix(config): handle missing config file gracefully

Previously, the application would crash if no config file was found.
Now it falls back to environment variables and CLI flags.

Fixes #456
```

```bash
docs(readme): update installation instructions

- Add Windows installation steps
- Clarify API key setup process
- Fix broken links
```

### Writing Good Commit Messages

**Do:**
- Use the imperative mood ("Add feature" not "Added feature")
- Start with a lowercase letter after the type/scope
- Keep the subject line under 50 characters
- Explain **what** and **why**, not **how**
- Reference issues and pull requests

**Don't:**
- End the subject line with a period
- Include implementation details in the subject
- Make commits too large (break them up)

---

## Pull Request Process

### Before Submitting

1. **Update your branch** with the latest from `main`:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all tests** and ensure they pass:
   ```bash
   make test
   ```

3. **Run the linter**:
   ```bash
   make lint
   ```

4. **Update documentation** if needed

5. **Build the project** to ensure it compiles:
   ```bash
   make build
   ```

### Creating the Pull Request

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Open a pull request** on GitHub

3. **Fill out the PR template** completely:
   - Describe your changes
   - Link related issues
   - Mark checklist items
   - Add screenshots if UI-related

4. **Respond to feedback** from reviewers promptly

### PR Title Guidelines

Use the same format as commit messages:
- `feat: Add dark mode support`
- `fix: Resolve config loading issue`
- `docs: Update contributing guidelines`

### Review Process

- At least one maintainer must approve your PR
- All CI checks must pass
- Address all review comments or explain why you disagree
- Keep the discussion focused and constructive

### After Approval

Once approved, a maintainer will merge your PR. You can then:

1. **Delete your feature branch**:
   ```bash
   git branch -d feature/your-feature-name
   git push origin --delete feature/your-feature-name
   ```

2. **Update your local main**:
   ```bash
   git checkout main
   git pull upstream main
   ```

---

## Testing Guidelines

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/ai/

# With coverage
go test -cover ./...

# With race detection
go test -race ./...

# Verbose output
go test -v ./...
```

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests for multiple scenarios
- Mock external dependencies (API calls, file I/O)
- Test both success and error paths
- Keep tests fast and independent

### Test Coverage

We aim for at least 70% code coverage. Check coverage with:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Documentation

### Code Documentation

- Add comments to all exported functions, types, and constants
- Use GoDoc conventions
- Explain complex algorithms or business logic
- Keep comments up-to-date with code changes

### User Documentation

When adding features or changing behavior:

1. Update the main README.md
2. Update INSTALL.md if installation changes
3. Add examples to demonstrate usage
4. Update configs/README.md if config changes

### Documentation Style

- Use clear, concise language
- Include code examples
- Use proper markdown formatting
- Add screenshots for UI changes
- Link to related documentation

---

## Getting Help

If you have questions or need help:

- 💬 [GitHub Discussions](https://github.com/abdomassoun/taskfix/discussions) - For general questions
- 🐛 [Issue Tracker](https://github.com/abdomassoun/taskfix/issues) - For bugs and feature requests
- 📖 [Documentation](https://github.com/abdomassoun/taskfix/wiki) - For guides and references

---

## Recognition

Contributors will be recognized in:
- The project README (if significant contribution)
- Release notes
- Git history

---

## License

By contributing to TaskFix, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to TaskFix! 🎉**
