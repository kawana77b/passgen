package uid

import (
	"regexp"
	"testing"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// TestNewUUID checks that v4 and v7 produce valid UUID strings.
func TestNewUUID(t *testing.T) {
	for _, kind := range []Kind{KindV4, KindV7} {
		got, err := New(kind)
		if err != nil {
			t.Fatalf("kind=%s: unexpected error: %v", kind, err)
		}
		if !uuidPattern.MatchString(got) {
			t.Errorf("kind=%s: %q is not a valid UUID", kind, got)
		}
	}
}

// TestNewNanoID checks that nanoid produces a non-empty string.
func TestNewNanoID(t *testing.T) {
	got, err := New(KindNanoID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Error("nanoid returned an empty string")
	}
}

// TestNewNCount checks that NewN returns exactly the requested number of IDs.
func TestNewNCount(t *testing.T) {
	const count = 10
	for _, kind := range []Kind{KindV4, KindV7, KindNanoID} {
		results, err := NewN(kind, count)
		if err != nil {
			t.Fatalf("kind=%s: unexpected error: %v", kind, err)
		}
		if len(results) != count {
			t.Errorf("kind=%s: expected %d results, got %d", kind, count, len(results))
		}
	}
}
