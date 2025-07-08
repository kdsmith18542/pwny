# Contributing to MSF-Go

Thank you for your interest in contributing to MSF-Go! This guide will help you get started with contributing to the project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Requests](#pull-requests)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)
- [Documentation](#documentation)

## 👥 Code of Conduct

This project adheres to the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you are expected to uphold this code.

## 🚀 Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
   ```bash
   git clone https://github.com/yourusername/msfgo.git
   cd msfgo
   ```
3. Set up the development environment:
   ```bash
   go mod tidy
   ```

## 🔧 Development Workflow

1. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the code style guidelines

3. Run tests:
   ```bash
   go test ./...
   ```

4. Commit your changes with a descriptive message:
   ```bash
   git commit -m "Add: Brief description of changes"
   ```

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Open a pull request against the `main` branch

## 🎨 Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Keep lines under 120 characters
- Use meaningful variable and function names
- Add comments for exported functions and types
- Write unit tests for new functionality

### Linting

Run the linter before submitting a PR:
```bash
golangci-lint run
```

## 🧪 Testing

- Write tests for new features and bug fixes
- Run all tests:
  ```bash
  go test ./...
  ```
- Run tests with coverage:
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

## 🔄 Pull Requests

1. Keep PRs focused on a single feature or bug fix
2. Update documentation as needed
3. Include tests for new functionality
4. Ensure all tests pass
5. Update the [CHANGELOG.md](CHANGELOG.md)
6. Request reviews from maintainers

## 🐛 Reporting Issues

When reporting issues, please include:

1. A clear title and description
2. Steps to reproduce the issue
3. Expected vs actual behavior
4. Version information (Go version, OS, etc.)
5. Any relevant logs or error messages

## 💡 Feature Requests

We welcome feature requests! Please:

1. Check if the feature already exists
2. Explain why this feature would be valuable
3. Include any relevant use cases

## 📚 Documentation

Good documentation is crucial. When adding new features:

1. Update relevant `.md` files
2. Add GoDoc comments for public APIs
3. Include usage examples
4. Update the [PROGRESS.md](docs/PROGRESS.md) if applicable

## 🤝 Community

Join our community on [Discord/Slack] to discuss ideas, ask questions, and collaborate with other contributors.

## 🙏 Thank You!

Your contributions make MSF-Go better for everyone. Thank you for your time and effort!
