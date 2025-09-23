# Tasker CLI Documentation

A comprehensive guide to the Tasker command-line task management application.

<p>
  <img alt="License" src="https://img.shields.io/static/v1?label=license&message=MIT&color=428EE1&labelColor=2F3136">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/eduardamirelly/tasker?color=428EE1&labelColor=2F3136">
  <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/eduardamirelly/tasker?color=428EE1&labelColor=2F3136">
</p>

<p>
  <a href="https://github.com/eduardamirelly/tasker/issues/new/choose">Report Bug</a>
  ·
  <a href="https://github.com/eduardamirelly/tasker/issues/new/choose">Suggestions</a>
</p>


## 📋 Table of Contents

- [Project Overview](#-project-overview)
- [Getting Started](#-getting-started)
- [Documentation Index](#-documentation-index)
- [Project Structure](#-project-structure)
- [Testing](#-testing)
- [Development](#-development)

## 🎯 Project Overview

Tasker is a simple, efficient command-line task management tool built with Go. It allows you to manage your daily tasks directly from the terminal with four core commands:

- **`add`** - Create new tasks with optional descriptions
- **`list`** - View all your tasks with completion status
- **`done`** - Mark tasks as completed
- **`export`** - Export all tasks to CSV format

### Key Features

- ✅ **Simple CLI Interface** - Easy-to-use commands
- 🗃️ **SQLite Database** - Local storage, no external dependencies
- 🚀 **Fast Performance** - Instant task operations
- 🧪 **Comprehensive Testing** - 50+ test cases ensuring reliability
- 📝 **Rich Task Details** - Titles, descriptions, timestamps
- 🎯 **Status Tracking** - Visual indicators for task completion
- 📤 **CSV Export** - Export tasks for backup and analysis

### Technology Stack

- **Language**: Go 1.24+
- **CLI Framework**: Cobra
- **Database**: SQLite3
- **Testing**: Go testing + Testify
- **Build**: Go modules

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite3 (included with Go sqlite driver)

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd tasker
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build -o tasker
   ```

4. **Run tasker**:
   ```bash
   ./tasker --help
   ```

### Quick Start

```bash
# Add your first task
./tasker add "Buy groceries" --description "Milk, eggs, bread"

# List all tasks
./tasker list

# Mark task as done (use the ID from list)
./tasker done 1

# Export all tasks to CSV
./tasker export -o my_tasks.csv

# View help for any command
./tasker add --help
```

## 📚 Documentation Index

### Core Documentation

- **[COMMANDS.md](./COMMANDS.md)** - Detailed command documentation
  - Code architecture and design patterns
  - Function-by-function explanations
  - Usage examples and error handling
  - Database integration details

### Testing Documentation

- **[tests/README.md](../tests/README.md)** - Complete testing guide
  - Test suite organization and structure
  - How each test works and what it validates
  - Running tests and debugging failures
  - Adding new tests and maintaining coverage

### Quick References

- **Command Usage**:
  ```bash
  tasker add "[title]" --description "[description]"
  tasker list
  tasker done [task_id]
  tasker export --output "[filename.csv]"
  ```

- **Test Commands**:
  ```bash
  go test ./tests/... -v          # Run all tests
  go test ./tests/... -cover      # Run with coverage
  ```

## 🏗️ Project Structure

```
tasker/
├── README.md                   # Main project README
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── main.go                     # Application entry point
├── tasker.db                   # SQLite database (created on first run)
│
├── cmd/                        # Command implementations
│   ├── root.go                # Root command and CLI setup
│   ├── add.go                 # Add command
│   ├── list.go                # List command
│   ├── done.go                # Done command
│   └── export.go              # Export command
│
├── models/                     # Data structures
│   └── task.go                # Task model definition
│
├── database/                   # Database layer
│   └── db.go                  # Database initialization and utilities
│
├── tests/                      # Test suite (organized separately)
│   ├── README.md              # Testing documentation
│   ├── test_helpers.go        # Test utilities and setup
│   ├── add_test.go            # Add command tests
│   ├── list_test.go           # List command tests
│   ├── done_test.go           # Done command tests
│   ├── export_test.go         # Export command tests
│   └── integration_test.go    # End-to-end tests
│
└── docs/                       # Documentation
    ├── README.md              # This file
    └── COMMANDS.md            # Command documentation
```

### Architecture Principles

1. **Clean Separation**: Commands, models, and database logic are separated
2. **Single Responsibility**: Each file has a clear, focused purpose
3. **Testability**: Business logic is isolated and easily testable
4. **Documentation**: Comprehensive docs for maintainability

## 🧪 Testing

### Test Organization

The test suite is organized in a dedicated `tests/` folder with:

- **Unit Tests**: Test individual functions and components
- **Integration Tests**: Test complete workflows and command interactions
- **Error Tests**: Validate error handling and edge cases
- **Concurrency Tests**: Ensure thread safety

### Running Tests

```bash
# Run all tests
go test ./tests/... -v

# Run specific test categories
go test ./tests/... -v -run TestAdd     # Add command tests
go test ./tests/... -v -run TestList    # List command tests
go test ./tests/... -v -run TestDone    # Done command tests
go test ./tests/... -v -run TestExport  # Export command tests

# Run with coverage
go test ./tests/... -cover

# Run integration tests
go test ./tests/... -v -run TestComplete
go test ./tests/... -v -run TestWorkflow
```

### Test Coverage

- **50+ Test Cases** covering all functionality
- **95%+ Code Coverage** across all commands
- **Isolated Testing** with temporary databases
- **Concurrent Testing** for thread safety validation
- **CSV Export Testing** with format validation and error scenarios

## 🛠️ Development

### Adding New Features

1. **Command Structure**: Follow the existing pattern in `cmd/`
2. **Database Changes**: Update schema in `database/db.go`
3. **Models**: Add/modify structs in `models/`
4. **Tests**: Add comprehensive tests in `tests/`
5. **Documentation**: Update relevant docs

### Code Standards

- **Error Handling**: Always return and handle errors appropriately
- **Database Safety**: Use parameterized queries, handle cleanup
- **User Experience**: Provide clear feedback and error messages
- **Testing**: Write tests for new functionality

### Building and Distribution

```bash
# Build for current platform
go build -o tasker

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o tasker-linux
GOOS=windows GOARCH=amd64 go build -o tasker.exe
GOOS=darwin GOARCH=amd64 go build -o tasker-mac
```

### Database Schema

```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    done BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);
```

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Add tests** for your changes
4. **Ensure tests pass**: `go test ./tests/... -v`
5. **Update documentation** as needed
6. **Commit changes**: `git commit -m "Add new feature"`
7. **Push to branch**: `git push origin feature/new-feature`
8. **Create Pull Request**

### Contribution Guidelines

- Write comprehensive tests for new functionality
- Follow existing code patterns and conventions
- Update documentation for user-facing changes
- Ensure all tests pass before submitting PR

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
