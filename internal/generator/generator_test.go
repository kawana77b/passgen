package generator

import (
	"strings"
	"testing"
)

// TestGenerateLength checks that the generated password has the requested length.
func TestGenerateLength(t *testing.T) {
	for _, length := range []int{8, 16, 32} {
		cfg := Config{Length: length, Rule: RuleAlNum, SymbolSet: SymbolSetCommon, Count: 1}
		got, err := Generate(cfg)
		if err != nil {
			t.Fatalf("length=%d: unexpected error: %v", length, err)
		}
		if len(got) != length {
			t.Errorf("length=%d: got len %d", length, len(got))
		}
	}
}

// TestGenerateRuleFull checks that a full-rule password contains at least one
// character from each required group (lower, upper, digit, symbol).
func TestGenerateRuleFull(t *testing.T) {
	cfg := Config{Length: 16, Rule: RuleFull, SymbolSet: SymbolSetCommon, Count: 1}
	got, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	groups := map[string]string{
		"lower":  lowerChars,
		"upper":  upperChars,
		"digit":  digitChars,
		"symbol": `!@#$%^&*`,
	}
	for name, chars := range groups {
		if !strings.ContainsAny(got, chars) {
			t.Errorf("password %q missing %s character", got, name)
		}
	}
}

// TestGenerateNCount checks that GenerateN returns exactly the requested number of passwords.
func TestGenerateNCount(t *testing.T) {
	const count = 10
	cfg := Config{Length: 16, Rule: RuleFull, SymbolSet: SymbolSetCommon, Count: count}
	results, err := GenerateN(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != count {
		t.Errorf("expected %d passwords, got %d", count, len(results))
	}
}
