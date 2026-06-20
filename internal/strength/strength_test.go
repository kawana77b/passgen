package strength

import (
	"math"
	"testing"
)

func TestEntropy(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantPool int
		// effective length = unique + (total-unique)*0.5
		wantEffective float64
	}{
		{"all unique lowercase", "abcdef", 26, 6},
		{"all unique lower+upper", "abcDEF", 52, 6},
		{"all unique lower+digit", "abc123", 36, 6},
		{"all unique lower+upper+digit", "abcDEF123", 62, 9},
		{"all unique all pools", "abcDEF123!", 94, 10},
		// repetition penalty: unique=1, repeated=15, effective=1+7.5=8.5
		{"all same char", "aaaaaaaaaaaaaaaa", 26, 8.5},
		// unique=5 (a,b,c,d,e), repeated=5, effective=5+2.5=7.5
		{"some repeats", "aabbccddee", 26, 7.5},
		{"empty", "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Entropy(tt.password)
			var want float64
			if tt.wantPool > 0 {
				want = tt.wantEffective * math.Log2(float64(tt.wantPool))
			}
			if math.Abs(got-want) > 1e-9 {
				t.Errorf("Entropy(%q) = %.4f, want %.4f (pool=%d, effective=%.1f)",
					tt.password, got, want, tt.wantPool, tt.wantEffective)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     Level
	}{
		// Weak: < 40 bits
		{"short lowercase", "abc", Weak},
		// pool=26, effective=8 → 8*4.700=37.6
		{"8 unique lowercase", "abcdefgh", Weak},
		// pool=26, unique=1, effective=8.5 → 8.5*4.700=39.9 (< 40)
		{"all same char 16", "aaaaaaaaaaaaaaaa", Weak},

		// Fair: 40–59 bits
		// pool=26, effective=9 → 9*4.700=42.3
		{"9 unique lowercase", "abcdefghi", Fair},
		// pool=62, effective=8 → 8*5.954=47.6
		{"lower+upper+digit 8", "abcDEF12", Fair},

		// Strong: 60–79 bits
		// pool=94, effective=10 → 10*6.555=65.5
		{"all pools 10 unique", "abcDEF1234!", Strong},
		// pool=62, effective=12 → 12*5.954=71.4
		{"lower+upper+digit 12 unique", "abcDEF123456", Strong},

		// Very Strong: >= 80 bits
		// pool=94, effective=13 → 13*6.555=85.2
		{"all pools 13 unique", "abcDEF123456!", VeryStrong},
		// pool=58(lower+symbol), unique=11, effective=11+5=16 → 16*5.858=93.7
		{"long passphrase with repeats", "correct-horse-battery", VeryStrong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Check(tt.password)
			if got != tt.want {
				t.Errorf("Check(%q) = %s (entropy=%.2f bits), want %s",
					tt.password, got, Entropy(tt.password), tt.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{Weak, "weak"},
		{Fair, "fair"},
		{Strong, "strong"},
		{VeryStrong, "very strong"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}
