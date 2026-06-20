package generator

import "fmt"

type SymbolSet string

const (
	SymbolSetSafe   SymbolSet = "safe"   // safe symbols unlikely to be rejected by websites
	SymbolSetCommon SymbolSet = "common" // common symbols (default)
	SymbolSetFull   SymbolSet = "full"   // all supported symbols
)

type SymbolSetInfo struct {
	Name    SymbolSet
	Charset string // display label and actual character pool
	Desc    string
}

var SymbolSets = []SymbolSetInfo{
	{SymbolSetSafe, `-_.`, "Safe — unlikely to be rejected by websites"},
	{SymbolSetCommon, `!@#$%^&*`, "Common (default)"},
	{SymbolSetFull, `!@#$%^&*()-_=+[]{};:,./?'`, "Full — all supported symbols"},
}

// LookupSymbolSet returns the SymbolSetInfo for the given symbol set, or an error if the symbol set is unknown.
func LookupSymbolSet(ss SymbolSet) (SymbolSetInfo, error) {
	for _, info := range SymbolSets {
		if info.Name == ss {
			return info, nil
		}
	}
	return SymbolSetInfo{}, fmt.Errorf("unknown symbol set %q: choose from safe, common, full", ss)
}
