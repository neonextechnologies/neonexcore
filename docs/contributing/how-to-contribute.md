# How to Contribute

Guidelines for contributing to Neonex Core.

---

## Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/neonexcore.git
   cd neonexcore
   ```

2. **Create a branch**
   ```bash
   git checkout -b feature/your-feature
   ```

3. **Make changes**
   - Write clean, documented code
   - Follow Go conventions
   - Add tests for new features

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Commit changes**
   ```bash
   git commit -m \"feat: add new feature\"
   ```

6. **Push and create PR**
   ```bash
   git push origin feature/your-feature
   ```

---

## Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `refactor`: Code refactoring
- `test`: Tests
- `chore`: Maintenance

**Example:**
```
feat: add user authentication module

Implement JWT-based authentication with:
- Login endpoint
- Token validation middleware
- User session management

Closes #123
```

---

## Code Style

### Go Conventions

```go
// ✅ Good
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}

// ❌ Bad
func create_user(req *CreateUserRequest) *User {
    // Implementation
}
```

### Documentation

```go
// CreateUser creates a new user account.
// Returns ErrInvalidEmail if email format is invalid.
// Returns ErrDuplicateEmail if email already exists.
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

---

## Testing Requirements

- Write unit tests for new code
- Maintain minimum 80% coverage
- Test error cases
- Include integration tests where appropriate

```go
func TestCreateUser(t *testing.T) {
    // Test implementation
}
```

---

## Pull Request Process

1. **Update documentation** if needed
2. **Add tests** for new features
3. **Update CHANGELOG.md**
4. **Request review** from maintainers
5. **Address feedback**
6. **Merge** when approved

---

## Code Review Guidelines

### For Contributors
- Respond to feedback promptly
- Keep PRs focused and small
- Update based on reviewer comments

### For Reviewers
- Be constructive and respectful
- Focus on code quality
- Approve when standards are met

---

## Next Steps

- [**Code of Conduct**](code-of-conduct.md)
- [**Development Setup**](development-setup.md)

---

**Questions?** Check [Support](../resources/support.md)
