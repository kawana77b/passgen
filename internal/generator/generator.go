package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type Config struct {
	Length    int
	Rule      Rule
	SymbolSet SymbolSet
	Count     int
}

func randInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func randChar(s string) (byte, error) {
	i, err := randInt(len(s))
	if err != nil {
		return 0, err
	}
	return s[i], nil
}

// Generate generates a single password based on the provided configuration.
func Generate(cfg Config) (string, error) {
	ruleInfo, err := LookupRule(cfg.Rule)
	if err != nil {
		return "", err
	}
	ssInfo, err := LookupSymbolSet(cfg.SymbolSet)
	if err != nil {
		return "", err
	}

	groups := ruleInfo.resolveGroups(ssInfo)

	if cfg.Length < len(groups) {
		return "", fmt.Errorf("length must be at least %d for the selected rule", len(groups))
	}

	charset := strings.Join(groups, "")
	buf := make([]byte, cfg.Length)

	// Shuffle the first len(groups) indices so each group gets a guaranteed position.
	indices := make([]int, cfg.Length)
	for i := range indices {
		indices[i] = i
	}
	for i := len(groups) - 1; i >= 0; i-- {
		j, err := randInt(i + 1)
		if err != nil {
			return "", err
		}
		indices[i], indices[j] = indices[j], indices[i]
	}
	for i, group := range groups {
		c, err := randChar(group)
		if err != nil {
			return "", err
		}
		buf[indices[i]] = c
	}

	// Fill remaining positions from the full charset.
	usedPositions := make(map[int]bool, len(groups))
	for i := range groups {
		usedPositions[indices[i]] = true
	}
	for i := range buf {
		if usedPositions[i] {
			continue
		}
		c, err := randChar(charset)
		if err != nil {
			return "", err
		}
		buf[i] = c
	}

	return string(buf), nil
}

// GenerateN generates multiple passwords based on the provided configuration.
func GenerateN(cfg Config) ([]string, error) {
	results := make([]string, cfg.Count)
	for i := range results {
		p, err := Generate(cfg)
		if err != nil {
			return nil, err
		}
		results[i] = p
	}
	return results, nil
}
