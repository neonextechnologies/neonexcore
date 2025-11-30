# How to Contribute

Welcome! We're excited that you want to contribute to Neonex Core. This guide will help you get started.

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:
- **Go 1.21+** installed
- **Git** configured with your GitHub account
- **PostgreSQL** or **MySQL** (optional, for testing)
- **Docker** (optional, for containerized testing)

### Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/neonexcore.git
cd neonexcore

# Add upstream remote
git remote add upstream https://github.com/neonexcore/neonexcore.git

# Verify remotes
git remote -v
```

---

## Development Workflow

### 1. Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name

# Or for bugfix
git checkout -b fix/issue-description
```

### Branch Naming Conventions

- **feature/** - New features (`feature/add-redis-cache`)
- **fix/** - Bug fixes (`fix/user-validation-error`)
- **docs/** - Documentation (`docs/update-readme`)
- **refactor/** - Code refactoring (`refactor/service-layer`)
- **test/** - Test additions (`test/user-repository`)
- **chore/** - Maintenance tasks (`chore/update-dependencies`)

### 2. Make Changes

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# Or use hot reload
neonex serve

# Make your changes
# ...
```

### 3. Write Tests

```go
// Example: modules/user/repository_test.go
func TestUserRepository_Create(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()
    
    repo := &UserRepository{db: db}
    
    // Test
    user, err := repo.Create(context.Background(), CreateUserData{
        Email:    "test@example.com",
        Password: "password123",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### 4. Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./modules/user/...

# Verbose output
go test -v ./...

# With race detector
go test -race ./...
```

### 5. Check Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for security issues
gosec ./...

# Check dependencies
go mod tidy
go mod verify
```

---

## Coding Standards

### Code Style

```go
// ✅ Good: Clear naming
func (s *UserService) CreateUser(ctx context.Context, email string) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }
    
    user := &User{Email: email}
    return s.repo.Create(ctx, user)
}

// ❌ Bad: Unclear naming
func (s *UserService) cu(c context.Context, e string) (*User, error) {
    if e == "" {
        return nil, errors.New("invalid")
    }
    
    u := &User{Email: e}
    return s.repo.Create(c, u)
}
```

### Error Handling

```go
// ✅ Good: Wrapped errors with context
if err := s.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// ❌ Bad: Lost error context
if err := s.repo.Create(ctx, user); err != nil {
    return nil, err
}
```

### Comments

```go
// ✅ Good: Clear documentation
// CreateUser creates a new user with the provided email.
// Returns ErrInvalidEmail if email is empty or invalid.
// Returns ErrDuplicateEmail if email already exists.
func (s *UserService) CreateUser(ctx context.Context, email string) (*User, error) {
    // ...
}

// ❌ Bad: No documentation
func (s *UserService) CreateUser(ctx context.Context, email string) (*User, error) {
    // ...
}
```

### Package Structure

```
modules/
  user/
    user.go          // Types and interfaces
    model.go         // Database models
    repository.go    // Data access
    service.go       // Business logic
    controller.go    // HTTP handlers
    routes.go        // Route definitions
    di.go           // Dependency injection
    module.json     // Module metadata
```

---

## Commit Guidelines

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code formatting (no logic change)
- **refactor**: Code refactoring
- **test**: Adding/updating tests
- **chore**: Maintenance tasks

### Examples

```bash
# Feature
git commit -m "feat(user): add email verification"

# Bug fix
git commit -m "fix(auth): resolve JWT token expiration issue"

# Documentation
git commit -m "docs(readme): update installation instructions"

# Multiple lines
git commit -m "feat(user): add password reset functionality

- Add password reset endpoint
- Implement email token generation
- Add password validation

Closes #123"
```

---

## Pull Request Process

### 1. Push Changes

```bash
# Push to your fork
git push origin feature/your-feature-name
```

### 2. Create Pull Request

1. Go to your fork on GitHub
2. Click "Pull Request"
3. Select your branch
4. Fill in the PR template:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
```

### 3. Code Review

- **Respond to feedback** promptly
- **Make requested changes** in new commits
- **Update tests** if needed
- **Resolve conflicts** with main branch

```bash
# Sync with upstream
git fetch upstream
git rebase upstream/main

# Resolve conflicts if any
# Then force push
git push -f origin feature/your-feature-name
```

### 4. Merge

Once approved:
- **Squash commits** if requested
- **Wait for CI** to pass
- **Maintainer will merge** your PR

---

## Testing Guidelines

### Test Coverage

Aim for:
- **90%+** for repositories
- **85%+** for services
- **75%+** for controllers

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure

```go
func TestName(t *testing.T) {
    // Arrange - Setup
    db := setupTestDB(t)
    repo := NewRepository(db)
    
    // Act - Execute
    result, err := repo.Method(params)
    
    // Assert - Verify
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## Documentation

### Code Documentation

```go
// Package user provides user management functionality.
package user

// User represents a user in the system.
type User struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Create creates a new user with the provided data.
// Returns ErrInvalidEmail if email is invalid.
// Returns ErrDuplicateEmail if email already exists.
func (r *Repository) Create(ctx context.Context, data CreateUserData) (*User, error) {
    // Implementation
}
```

### Markdown Documentation

When updating docs:
- Use **clear headings**
- Include **code examples**
- Add **tables** for reference
- Use **warnings** for important notes
- Link to **related docs**

---

## Issue Reporting

### Bug Reports

Include:
- **Description** - What happened?
- **Expected** - What should happen?
- **Steps to Reproduce** - How to trigger the bug?
- **Environment** - Go version, OS, database
- **Logs** - Relevant error messages

```markdown
**Description**
User creation fails with validation error

**Expected Behavior**
User should be created successfully

**Steps to Reproduce**
1. POST /api/users with valid data
2. See error response

**Environment**
- Go: 1.21
- OS: Ubuntu 22.04
- Database: PostgreSQL 15

**Logs**
```
Error: invalid email format
```
```

### Feature Requests

Include:
- **Problem** - What problem does this solve?
- **Solution** - Proposed solution
- **Alternatives** - Other options considered
- **Use Case** - Real-world example

---

## Community Guidelines

### Communication

- Be **respectful** and **inclusive**
- **Help others** learn and grow
- **Ask questions** when unclear
- **Share knowledge** and experiences

### Getting Help

- Check **documentation** first
- Search **existing issues**
- Ask in **GitHub Discussions**
- Join our **Discord server**

### Recognition

Contributors are recognized in:
- **CHANGELOG.md** - Feature/fix credits
- **README.md** - Contributor list
- **Release notes** - Major contributions

---

## Development Tips

### Debugging

```bash
# Enable debug mode
export DEBUG=true

# Use Delve debugger
dlv debug main.go

# Add breakpoint
(dlv) break main.main
(dlv) continue
```

### Performance Testing

```bash
# Benchmark tests
go test -bench=. -benchmem

# CPU profiling
go test -cpuprofile=cpu.out
go tool pprof cpu.out

# Memory profiling
go test -memprofile=mem.out
go tool pprof mem.out
```

### Database Migrations

```go
// Create migration
neonex migrate create add_users_table

// Run migrations
neonex migrate up

// Rollback
neonex migrate down
```

---

## Resources

- [**Code of Conduct**](code-of-conduct.md) - Community standards
- [**Development Setup**](development-setup.md) - Detailed setup guide
- [**Architecture**](../introduction/architecture.md) - System design
- [**Style Guide**](https://golang.org/doc/effective_go) - Go best practices

---

## Questions?

- **GitHub Issues** - Bug reports, feature requests
- **GitHub Discussions** - Questions, ideas
- **Discord** - Real-time chat
- **Email** - team@neonexcore.dev

---

Thank you for contributing to Neonex Core! 🚀
