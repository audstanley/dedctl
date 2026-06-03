# Phase 11: Multi-Hash Password Authentication

## Overview

Replace the single SHA-256 password hashing with three configurable options per user: **sha512**, **bcrypt**, and **plain**. The default is **sha512** for stronger security. Existing users with empty `password_type` will need to regenerate hashes.

## Design

### Password Types

| Type | Description | Use Case |
|------|-------------|----------|
| `sha512` | SHA-512 hex digest | Recommended default, strong and fast |
| `bcrypt` | bcrypt with cost 10 | Recommended for production, slow/adaptive |
| `plain` | Raw string comparison | Only for testing, NOT recommended |

### Config File Format

```yaml
users:
  - username: admin
    password_hash: "<sha512 hex digest>"
    password_type: sha512
    is_admin: true
  - username: operator
    password_hash: "$2b$12$..."
    password_type: bcrypt
    is_admin: false
  - username: viewer
    password_hash: "secret123"
    password_type: plain
    is_admin: false
```

### Hash Command

```bash
# Generate a SHA-512 hash
steamctl hash "mypassword" sha512
# Output: a1b2c3...

# Generate a bcrypt hash
steamctl hash "mypassword" bcrypt
# Output: $2b$12$xyz...

# Generate a plain password (just echoes back)
steamctl hash "mypassword" plain
# Output: mypassword
```

## Files to Change

### 1. `internal/config/config.go`
- Add `PasswordType string` field to `UserConfig` struct
- Add mapstructure tag: `mapstructure:"password_type"`

### 2. `internal/service/auth.go`
- Add `PasswordType string` field to `UserInfo` struct
- Add imports: `crypto/sha512`, `golang.org/x/crypto/bcrypt`
- Replace `Login` method with `verifyPassword` helper:
  - If `password_type == ""` or `== "sha512"`: use SHA-512 hex digest
  - If `password_type == "bcrypt"`: use `bcrypt.CompareHashAndPassword`
  - If `password_type == "plain"`: direct string comparison
- Default (empty type): sha512

### 3. `internal/app/server.go`
- Wire `PasswordType` from config into `UserInfo` when building users slice

### 4. `go.mod` + `go.sum`
- Add `golang.org/x/crypto` as direct dependency
- Run `go mod tidy`

### 5. `cmd/steamctl/hash.go` (new file)
- New cobra command: `steamctl hash <password> <type>`
- Output plain text hash to stdout
- Validate type (sha512, bcrypt, plain)
- Error on invalid type

### 6. `internal/service/auth_test.go`
- Replace existing `hashPassword` helper with `hashPasswordSHA512`
- Add `hashPasswordBcrypt` helper
- Update existing tests to use explicit `password_type`
- Add test: `TestLoginSuccessBcrypt`
- Add test: `TestLoginSuccessSHA512`
- Add test: `TestLoginSuccessPlain`
- Add test: `TestLoginWrongPasswordSHA512`
- Add test: `TestLoginWrongPasswordBcrypt`
- Add test: `TestLoginWrongPasswordPlain`
- Add test: `TestLoginInvalidPasswordType`

### 7. `configs/config.yaml`
- Add `password_type` to existing users
- Add comments explaining each type
- Regenerate admin/operator hashes using `steamctl hash`
- Keep admin as sha512, add bcrypt example comment

### 8. `PROJECT_SUMMARY.md`
- Update Phase 10 status to committed
- Add Phase 11 entry with completion status
- Update Known Issues (remove SHA-256 reference)
- Add `steamctl hash` to "How to Run" section
- Update technology stack section

## Implementation Order

1. `config.go` — Add struct field (no deps)
2. `auth.go` — Refactor Login (needs bcrypt import)
3. `server.go` — Wire PasswordType through
4. `hash.go` — New command file
5. `go.mod` — Add bcrypt dependency
6. `auth_test.go` — Expand test suite
7. `config.yaml` — Update example users
8. `PROJECT_SUMMARY.md` — Update documentation

## Testing Strategy

- Unit tests cover all three hash types
- `TestTokenIsValidAgainstUtils` still passes (JWT is unchanged)
- `TestTokenFailsWithWrongSecret` still passes (JWT is unchanged)
- New integration test: login with each type, verify correct token returned

## Migration Notes for Users

Existing users need to regenerate their password hashes:

```bash
# Generate new SHA-512 hash
steamctl hash "admin123" sha512

# Generate new bcrypt hash (recommended for production)
steamctl hash "admin123" bcrypt

# Update config.yaml with the new hash and password_type
```
