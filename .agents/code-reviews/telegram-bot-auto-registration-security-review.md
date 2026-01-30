# Code Review: Telegram Bot Auto-Registration and Command System

**Date**: 2026-01-30
**Reviewer**: AI Code Review Agent
**Scope**: Telegram bot auto-registration feature implementation

---

## Stats

- **Files Modified**: 7
- **Files Added**: 2
- **Files Deleted**: 0
- **New lines**: 123
- **Deleted lines**: 93

---

## Summary

The implementation adds auto-registration functionality for Telegram bot users and registers bot commands with Telegram API. The code follows existing patterns and conventions well. However, there are **2 critical security issues** and **1 high-severity issue** that must be addressed before deployment.

---

## Critical Issues

### Issue 1: Predictable Password Pattern - Security Vulnerability

**severity**: critical
**file**: bots/telegram/bot/registration.go
**line**: 39
**issue**: Deterministic password pattern `tg_secure_%d` is predictable and vulnerable to brute force attacks
**detail**: The password is generated using only the Telegram user ID (`fmt.Sprintf("tg_secure_%d", msg.From.ID)`), which is a sequential integer. An attacker who knows or can guess a Telegram user ID can easily construct the password and gain unauthorized access to the account. Telegram user IDs are often exposed in public groups, channels, and can be enumerated. This creates a critical security vulnerability where any user's account can be compromised.

**Example attack scenario**:
1. Attacker finds user's Telegram ID (e.g., 123456789) from a public group
2. Attacker constructs password: `tg_secure_123456789`
3. Attacker uses `/link` command or web login with email `tg_123456789@telegram.bot` and the constructed password
4. Attacker gains full access to victim's account

**suggestion**: 
```go
// Use cryptographically secure random password with HMAC for verification
import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
)

// Store a secret key in environment variable (DO NOT hardcode)
var telegramPasswordSecret = []byte(os.Getenv("TELEGRAM_PASSWORD_SECRET"))

// Generate password using HMAC with secret key
func generateTelegramPassword(telegramID int64) (string, error) {
    if len(telegramPasswordSecret) == 0 {
        return "", fmt.Errorf("TELEGRAM_PASSWORD_SECRET not configured")
    }
    
    h := hmac.New(sha256.New, telegramPasswordSecret)
    h.Write([]byte(fmt.Sprintf("%d", telegramID)))
    hash := h.Sum(nil)
    
    // Add random salt for additional security
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    combined := append(hash, salt...)
    return base64.URLEncoding.EncodeToString(combined), nil
}
```

Alternatively, if you want to keep deterministic passwords for simplicity, use HMAC with a secret key:
```go
// Simpler deterministic approach with HMAC
func generateTelegramPassword(telegramID int64) string {
    h := hmac.New(sha256.New, telegramPasswordSecret)
    h.Write([]byte(fmt.Sprintf("%d", telegramID)))
    return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
```

This requires adding `TELEGRAM_PASSWORD_SECRET` to environment variables and `.env.example`.

---

### Issue 2: Password Logged in Error Messages

**severity**: critical
**file**: bots/telegram/bot/registration.go
**line**: 82
**issue**: Password may be exposed in error logs through error message details
**detail**: When both registration and login fail, the error is logged with `reg_err=%v, login_err=%v`. If the gRPC error messages contain the password (which some frameworks do in validation errors), it will be logged in plaintext. This violates security best practices and could expose passwords in log files.

**suggestion**:
```go
// Sanitize error messages before logging
log.Printf("[ERROR] Failed to register/login Telegram user %d: registration_failed=%t, login_failed=%t", 
    msg.From.ID, 
    err != nil, 
    loginErr != nil)

// Or log error types without details
if err != nil {
    log.Printf("[ERROR] Registration failed for Telegram user %d: %T", msg.From.ID, err)
}
if loginErr != nil {
    log.Printf("[ERROR] Login failed for Telegram user %d: %T", msg.From.ID, loginErr)
}
```

---

## High Severity Issues

### Issue 3: Context Timeout May Be Too Short for Network Issues

**severity**: high
**file**: bots/telegram/bot/registration.go
**line**: 24
**issue**: 5-second timeout may be insufficient for registration + login fallback flow
**detail**: The function uses a single 5-second context timeout for potentially two sequential gRPC calls (Register + Login). If the registration call takes 3-4 seconds due to network latency or database load, the login fallback may fail due to context deadline exceeded, even though the user exists. This creates a poor user experience where existing users cannot log in during high load.

**suggestion**:
```go
// Use separate contexts for each operation
ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel1()

resp, err := h.clients.User.Register(ctx1, &userpb.RegisterRequest{
    Email:    email,
    Password: password,
    Name:     name,
})

// ... handle registration success ...

// Create new context for login attempt
ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel2()

loginResp, loginErr := h.clients.User.Login(ctx2, &userpb.LoginRequest{
    Email:    email,
    Password: password,
})
```

Or increase the timeout to 10 seconds:
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

---

## Medium Severity Issues

### Issue 4: Command Descriptions Exceed Telegram's Recommended Length

**severity**: medium
**file**: bots/telegram/bot/bot.go
**lines**: 44-69
**issue**: Bilingual command descriptions may be truncated in some Telegram clients
**detail**: While Telegram's API allows up to 256 characters for command descriptions, the recommended length is 3-256 characters, and many clients display only the first 60-80 characters. The bilingual descriptions (e.g., "Start bot and create account | Начать работу и создать аккаунт") are 60+ characters and may be cut off in mobile clients, showing only the English part.

**suggestion**: Consider using Telegram's command scope feature to register different commands for different languages, or keep descriptions shorter:
```go
// Option 1: Shorter descriptions
{
    Command:     "start",
    Description: "Start bot and register",
},

// Option 2: Use BotCommandScope to register language-specific commands
// This requires detecting user language and registering commands per user
// See: https://core.telegram.org/bots/api#setmycommands
```

For now, the current implementation is acceptable but may need refinement based on user feedback.

---

### Issue 5: Missing Input Validation for Telegram User Data

**severity**: medium
**file**: bots/telegram/bot/registration.go
**lines**: 30-37
**issue**: No validation of Telegram user name length or special characters
**detail**: The code constructs a name from `FirstName + " " + LastName` or falls back to `UserName` without validating length or sanitizing special characters. Telegram allows names up to 64 characters, but the backend may have different limits. Additionally, names could contain special characters that might cause issues in the database or UI.

**suggestion**:
```go
// Validate and sanitize name
name := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
if name == "" {
    name = msg.From.UserName
}
if name == "" {
    name = fmt.Sprintf("User%d", msg.From.ID)
}

// Limit name length to prevent database issues
const maxNameLength = 100
if len(name) > maxNameLength {
    name = name[:maxNameLength]
}

// Optional: Sanitize special characters if needed
// name = sanitizeName(name)
```

---

## Low Severity Issues

### Issue 6: Test Coverage Could Be Improved

**severity**: low
**file**: bots/telegram/bot/bot_test.go
**lines**: 1-40
**issue**: Test only validates command structure, not actual registration logic
**detail**: The new test `TestRegisterCommandsStructure` only validates the command data structure but doesn't test the actual `registerCommands()` method or the registration flow. While this is acceptable for a unit test (since it requires a real Telegram bot token), integration tests should be added.

**suggestion**: Add integration test documentation or a test that mocks the Telegram API:
```go
// Add to bot_test.go
func TestRegisterCommandsWithMockAPI(t *testing.T) {
    // This would require creating a mock for tgbotapi.BotAPI
    // For now, document that manual testing is required
    t.Skip("Requires manual testing with real Telegram bot token")
}
```

Also add a comment in the test file explaining the testing strategy.

---

### Issue 7: Magic Numbers in Code

**severity**: low
**file**: bots/telegram/bot/registration.go
**line**: 24
**issue**: Timeout value (5 seconds) is hardcoded
**detail**: The 5-second timeout is a magic number that should be a named constant for maintainability.

**suggestion**:
```go
const (
    registrationTimeout = 5 * time.Second
)

// In function:
ctx, cancel := context.WithTimeout(context.Background(), registrationTimeout)
```

---

### Issue 8: Inconsistent Error Message Format

**severity**: low
**file**: bots/telegram/bot/registration.go
**line**: 83
**issue**: User-facing error message doesn't match the pattern of other messages
**detail**: The error message "❌ Failed to create account. Please try again later." is hardcoded in the function, while other messages use constants from `messages.go` (e.g., `MsgServiceError`). This creates inconsistency and makes internationalization harder.

**suggestion**:
```go
// Add to messages.go:
const (
    // ...
    MsgRegistrationFailed = "❌ Failed to create account. Please try again later."
)

// In registration.go:
h.sendMessage(msg.Chat.ID, MsgRegistrationFailed, nil)
```

---

## Positive Observations

✅ **Good practices observed**:

1. **Thread Safety**: Proper use of mutex locks in `handleStart` to prevent race conditions during concurrent registrations
2. **Error Handling**: Comprehensive error checking for nil pointers and response validation
3. **Logging**: Consistent use of structured logging with `[INFO]`, `[ERROR]`, `[WARN]` prefixes
4. **Code Organization**: Clean separation of concerns between bot initialization, command handling, and registration
5. **Graceful Degradation**: Bot continues to start even if command registration fails
6. **Backward Compatibility**: Existing `/link` command functionality is preserved
7. **Test Coverage**: Existing tests continue to pass, new test added for command structure
8. **Documentation**: Clear comments explaining the registration flow

---

## Recommendations

### Immediate Actions (Before Deployment)

1. **CRITICAL**: Fix the predictable password vulnerability (Issue 1)
2. **CRITICAL**: Sanitize error logs to prevent password exposure (Issue 2)
3. **HIGH**: Increase context timeout or use separate contexts (Issue 3)

### Short-term Improvements

4. Add `TELEGRAM_PASSWORD_SECRET` to environment configuration
5. Add input validation for user names (Issue 5)
6. Extract magic numbers to constants (Issue 7)
7. Move error messages to constants (Issue 8)

### Long-term Enhancements

8. Add integration tests with mocked Telegram API
9. Consider implementing command scopes for language-specific commands
10. Add monitoring/alerting for failed registrations
11. Implement rate limiting for registration attempts to prevent abuse

---

## Conclusion

The implementation is well-structured and follows existing code patterns. However, **the critical security vulnerability with predictable passwords must be fixed before deployment**. The password generation mechanism needs to use either:

1. HMAC with a secret key (deterministic but secure), or
2. Cryptographically secure random passwords stored in the database

Without this fix, all Telegram-originated accounts are vulnerable to unauthorized access.

After addressing the critical and high-severity issues, the code will be production-ready.

---

## Approval Status

**Status**: ❌ **BLOCKED - Critical Security Issues**

**Required fixes before merge**:
- [ ] Issue 1: Implement secure password generation
- [ ] Issue 2: Sanitize error logs
- [ ] Issue 3: Fix context timeout handling

**Recommended fixes**:
- [ ] Issue 5: Add input validation
- [ ] Issue 7: Extract constants
- [ ] Issue 8: Use message constants

Once critical issues are resolved, re-review and approve.
