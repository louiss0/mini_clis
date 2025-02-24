# Go Utility Projects

A collection of useful Go applications and utilities including task management, form handling, password generation, and more.

## Project Structure

This workspace contains several independent Go projects:

- **task-list**: CLI task manager for daily task organization
- **form**: Form application with text input capabilities
- **counter**: Simple counter application
- **pass-gen**: Feature-rich password generator
- **shared**: Common utilities and error handling

## Task List CLI

A command-line task manager for organizing your daily tasks.

### Features
- Add new tasks
- Delete existing tasks
- Edit task details
- List all tasks

### Usage
```bash
# Install
go install ./task-list

# Add a task
task-list add "Complete project documentation"

# List all tasks
task-list list

# Edit a task
task-list edit <task-id> "Updated task description"

# Delete a task
task-list delete <task-id>
```

## Form Application

A Go application for handling form input with text processing capabilities.

### Features
- Text input handling
- Form validation
- Input processing

### Usage
```bash
# Install
go install ./form

# Run the form application
form
```

## Counter Application

A simple counter implementation in Go.

### Features
- Increment/decrement functionality
- Counter state management

### Usage
```bash
# Install
go install ./counter

# Run the counter application
counter
```

## Password Generator

A versatile password generator with multiple generation methods.

### Features
- Encode: Convert input text to secure passwords
- Numeric: Generate numeric passwords
- Words: Create memorable passwords using words
- Leetspeak: Transform text into leetspeak passwords

### Usage
```bash
# Install
go install ./pass-gen

# Generate encoded password
pass-gen encode "your-text"

# Generate numeric password
pass-gen numeric

# Generate word-based password
pass-gen words

# Generate leetspeak password
pass-gen leetspeak "your-text"
```

## Shared Utilities

Common utilities and error handling used across projects.

### Features
- Custom error types
- Shared helper functions
- Common interfaces

### Usage
```go
import "path/to/shared"

// Example error handling
if err := someFunction(); err != nil {
    return shared.NewError("operation failed", err)
}
```

## Installation

To install all projects:

```bash
# Clone the repository
git clone <repository-url>

# Change to project directory
cd <project-directory>

# Install all dependencies
go mod download

# Build all projects
go build ./...
```

Each subproject can be built and installed independently using `go install` in its respective directory.

