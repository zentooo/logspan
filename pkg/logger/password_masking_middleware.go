package logger

import (
	"regexp"
	"strings"
)

// PasswordMaskingMiddleware creates a middleware that masks password-related information
// in log entries. It searches for password patterns in both message and fields.
type PasswordMaskingMiddleware struct {
	// MaskString is the string used to replace password values (default: "***")
	MaskString string

	// PasswordKeys are field keys that should be masked (case-insensitive)
	PasswordKeys []string

	// PasswordPatterns are regex patterns to detect passwords in messages
	PasswordPatterns []*regexp.Regexp
}

// NewPasswordMaskingMiddleware creates a new password masking middleware with default settings
func NewPasswordMaskingMiddleware() *PasswordMaskingMiddleware {
	return &PasswordMaskingMiddleware{
		MaskString: "***",
		PasswordKeys: []string{
			"password", "passwd", "pwd", "pass",
			"secret", "token", "key", "auth",
			"credential", "credentials", "api_key",
			"access_token", "refresh_token",
		},
		PasswordPatterns: []*regexp.Regexp{
			// Pattern for "password=value" or "password: value"
			regexp.MustCompile(`(?i)(password|passwd|pwd|pass|secret|token|key|auth|credential|api_key|access_token|refresh_token)[\s]*[=:]\s*[^\s]+`),
			// Pattern for JSON-like structures: "password":"value"
			regexp.MustCompile(`(?i)"(password|passwd|pwd|pass|secret|token|key|auth|credential|api_key|access_token|refresh_token)"\s*:\s*"[^"]+"`),
		},
	}
}

// WithMaskString sets a custom mask string
func (pmm *PasswordMaskingMiddleware) WithMaskString(maskString string) *PasswordMaskingMiddleware {
	pmm.MaskString = maskString
	return pmm
}

// WithPasswordKeys sets custom password field keys
func (pmm *PasswordMaskingMiddleware) WithPasswordKeys(keys []string) *PasswordMaskingMiddleware {
	pmm.PasswordKeys = keys
	return pmm
}

// AddPasswordKey adds a password field key to the existing list
func (pmm *PasswordMaskingMiddleware) AddPasswordKey(key string) *PasswordMaskingMiddleware {
	pmm.PasswordKeys = append(pmm.PasswordKeys, key)
	return pmm
}

// AddPasswordPattern adds a regex pattern for password detection in messages
func (pmm *PasswordMaskingMiddleware) AddPasswordPattern(pattern *regexp.Regexp) *PasswordMaskingMiddleware {
	pmm.PasswordPatterns = append(pmm.PasswordPatterns, pattern)
	return pmm
}

// Middleware returns the middleware function
func (pmm *PasswordMaskingMiddleware) Middleware() Middleware {
	return func(entry *LogEntry, next func(*LogEntry)) {
		// Mask passwords in message
		entry.Message = pmm.maskPasswordsInMessage(entry.Message)

		// Continue to next middleware
		next(entry)
	}
}

// maskPasswordsInMessage masks password patterns in the message string
func (pmm *PasswordMaskingMiddleware) maskPasswordsInMessage(message string) string {
	result := message

	for _, pattern := range pmm.PasswordPatterns {
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			// Find the separator (= or :) and replace everything after it
			if strings.Contains(match, "=") {
				parts := strings.SplitN(match, "=", 2)
				if len(parts) == 2 {
					return parts[0] + "=" + pmm.MaskString
				}
			} else if strings.Contains(match, ":") {
				parts := strings.SplitN(match, ":", 2)
				if len(parts) == 2 {
					// Handle JSON-like format: preserve quotes if they exist
					value := strings.TrimSpace(parts[1])
					if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
						return parts[0] + ":\"" + pmm.MaskString + "\""
					}
					return parts[0] + ": " + pmm.MaskString
				}
			}
			return pmm.MaskString
		})
	}

	return result
}

// isPasswordKey checks if a key is considered a password field (case-insensitive)
func (pmm *PasswordMaskingMiddleware) isPasswordKey(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, passwordKey := range pmm.PasswordKeys {
		if strings.ToLower(passwordKey) == lowerKey {
			return true
		}
	}
	return false
}
