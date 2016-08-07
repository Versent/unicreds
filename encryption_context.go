package unicreds

import (
	"fmt"
	"strings"
)

// EncryptionContextValue key value with helper methods for flag parser
type EncryptionContextValue map[string]*string

// NewEncryptionContextValue create a new encryption context
func NewEncryptionContextValue() *EncryptionContextValue {
	m := make(EncryptionContextValue)
	return &m
}

// Set converts a flag value into an encryption context key value
func (h *EncryptionContextValue) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected KEY:VALUE got '%s'", value)
	}
	(*h)[parts[0]] = &parts[1]
	return nil
}

func (h *EncryptionContextValue) String() string {
	return ""
}

// IsCumulative flag this value as cumulative
func (h *EncryptionContextValue) IsCumulative() bool {
	return true
}
