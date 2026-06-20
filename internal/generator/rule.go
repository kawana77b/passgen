package generator

import "fmt"

type Rule string

const (
	RuleLower Rule = "lower" // lowercase letters only
	RuleMixed Rule = "mixed" // upper and lowercase letters
	RuleAlNum Rule = "alnum" // letters and digits
	RuleFull  Rule = "full"  // letters, digits, and symbols
)

const (
	lowerChars = "abcdefghijklmnopqrstuvwxyz"
	upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars = "0123456789"
)

type RuleInfo struct {
	Name        Rule
	Charset     string // display label shown in list-rules
	Description string

	baseGroups []string // character groups guaranteed to appear in every password
	useSymbols bool     // whether a symbol group is appended (requires SymbolSet)
}

var Rules = []RuleInfo{
	{RuleLower, "a-z", "Lowercase letters only", []string{lowerChars}, false},
	{RuleMixed, "a-z, A-Z", "Upper and lowercase letters", []string{lowerChars, upperChars}, false},
	{RuleAlNum, "a-z, A-Z, 0-9", "Letters and digits", []string{lowerChars, upperChars, digitChars}, false},
	{RuleFull, "a-z, A-Z, 0-9, symbols", "Letters, digits, and symbols (default)", []string{lowerChars, upperChars, digitChars}, true},
}

// LookupRule returns the RuleInfo for the given rule, or an error if the rule is unknown.
func LookupRule(r Rule) (RuleInfo, error) {
	for _, info := range Rules {
		if info.Name == r {
			return info, nil
		}
	}
	return RuleInfo{}, fmt.Errorf("unknown rule %q: choose from lower, mixed, alnum, full", r)
}

// resolveGroups returns the character groups for this rule.
// ss is only consulted when useSymbols is true.
func (ri RuleInfo) resolveGroups(ss SymbolSetInfo) []string {
	if !ri.useSymbols {
		return ri.baseGroups
	}

	groups := make([]string, len(ri.baseGroups)+1)
	copy(groups, ri.baseGroups)
	groups[len(ri.baseGroups)] = ss.Charset

	return groups
}
