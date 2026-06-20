package strength

import "math"

type Level int

const (
	Weak       Level = 1
	Fair       Level = 2
	Strong     Level = 3
	VeryStrong Level = 4
)

func (l Level) String() string {
	switch l {
	case Weak:
		return "weak"
	case Fair:
		return "fair"
	case Strong:
		return "strong"
	case VeryStrong:
		return "very strong"
	default:
		return "unknown"
	}
}

// Character pools. Pool sizes must match the actual character sets.
var pools = []struct {
	chars string
	size  int
}{
	{"abcdefghijklmnopqrstuvwxyz", 26},
	{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", 26},
	{"0123456789", 10},
	{"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", 32}, // printable ASCII symbols (RFC 4086)
}

// Entropy returns the estimated entropy of the password in bits.
//
// Formula: H = L_eff × log2(N)
//
//	N     = pool size (sum of sizes of each character pool present in the password)
//	L_eff = effective length, which penalises repeated characters:
//	        L_eff = unique_count + (total_length - unique_count) × 0.5
//
// Repeated characters reduce L_eff because each additional occurrence of the
// same character contributes less unpredictability than a novel one. A password
// consisting entirely of one character (e.g. "aaaa…") therefore scores much
// lower than one with the same length but all distinct characters.
//
// The base formula and pool sizes follow RFC 4086 "Randomness Requirements for
// Security" (Section 8.1): https://www.rfc-editor.org/rfc/rfc4086#section-8.1
func Entropy(password string) float64 {
	runes := []rune(password)
	if len(runes) == 0 {
		return 0
	}

	poolSize := 0
	for _, p := range pools {
		for _, c := range runes {
			if containsRune(p.chars, c) {
				poolSize += p.size
				break
			}
		}
	}
	if poolSize == 0 {
		return 0
	}

	unique := uniqueCount(runes)
	repeated := len(runes) - unique
	effectiveLength := float64(unique) + float64(repeated)*0.5

	return effectiveLength * math.Log2(float64(poolSize))
}

// Check evaluates the strength of the given password and returns a Level from Weak to VeryStrong.
//
// Strength thresholds are based on entropy in bits, derived from the guidance in
// NIST SP 800-63B (Section 5.1.1) and the general industry convention for
// minimum acceptable entropy in authentication contexts:
//
//	< 40 bits  → Weak        (insufficient for most uses)
//	40–59 bits → Fair        (acceptable for low-risk contexts)
//	60–79 bits → Strong      (recommended for general use)
//	≥ 80 bits  → Very Strong (suitable for high-security contexts)
//
// See:
//   - NIST SP 800-63B: https://pages.nist.gov/800-63-3/sp800-63b.html
//   - RFC 4086:        https://www.rfc-editor.org/rfc/rfc4086
func Check(password string) Level {
	e := Entropy(password)
	switch {
	case e >= 80:
		return VeryStrong
	case e >= 60:
		return Strong
	case e >= 40:
		return Fair
	default:
		return Weak
	}
}

func uniqueCount(runes []rune) int {
	seen := make(map[rune]struct{}, len(runes))
	for _, r := range runes {
		seen[r] = struct{}{}
	}
	return len(seen)
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
