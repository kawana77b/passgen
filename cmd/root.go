package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/atotto/clipboard"
	"github.com/kawana77b/passgen/internal/generator"
	"github.com/kawana77b/passgen/internal/strength"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var (
	length       int
	rule         string
	symbolSet    string
	count        int
	showStrength bool
	clip         bool
)

// resolved holds the parsed and validated values built in PreRunE.
var resolved struct {
	config generator.Config
}

var rootCmd = &cobra.Command{
	Use:   "passgen",
	Short: "Random password generator",
	Long: `passgen generates cryptographically secure random passwords using
the OS-provided CSPRNG.
By default it produces 10 passwords of 16 characters using letters, digits, and common symbols (!@#$%^&*).`,
	Example: `  passgen                          # 10 passwords, rule=full, symbol-set=common
  passgen -l 24                    # 24-character passwords
  passgen -n 5                     # 5 passwords
  passgen -r alnum                 # letters and digits only
  passgen -r full -s safe          # symbols restricted to -_.
  passgen -r full -s full -l 32    # all symbols, 32 characters
  passgen --show-strength          # show entropy level next to each password

  passgen uid v4                   # 10 UUIDs v4
  passgen uid v7 -n 5              # 5 UUIDs v7 (time-sortable)
  passgen uid nanoid               # 10 NanoIDs`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		r, err := parseRule(rule)
		if err != nil {
			return err
		}
		ss, err := parseSymbolSet(symbolSet)
		if err != nil {
			return err
		}
		if length < 1 {
			return fmt.Errorf("length must be at least 1")
		}
		if count < 1 {
			return fmt.Errorf("count must be at least 1")
		}
		resolved.config = generator.Config{
			Length:    length,
			Rule:      r,
			SymbolSet: ss,
			Count:     count,
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		passwords, err := generator.GenerateN(resolved.config)
		if err != nil {
			return err
		}

		sort.Slice(passwords, func(i, j int) bool {
			return strength.Entropy(passwords[i]) > strength.Entropy(passwords[j])
		})

		for _, p := range passwords {
			if showStrength {
				level := strength.Check(p)
				fmt.Printf("%-*s  [%s]\n", resolved.config.Length, p, level.String())
			} else {
				fmt.Println(p)
			}
		}

		if clip {
			if err := clipboard.WriteAll(passwords[0]); err != nil {
				return fmt.Errorf("clipboard: %w", err)
			}
			fmt.Fprintf(os.Stderr, "copied to clipboard: %s\n", passwords[0])
		}

		return nil
	},
}

func parseRule(s string) (generator.Rule, error) {
	switch s {
	case "lower":
		return generator.RuleLower, nil
	case "mixed":
		return generator.RuleMixed, nil
	case "alnum":
		return generator.RuleAlNum, nil
	case "full":
		return generator.RuleFull, nil
	default:
		return "", fmt.Errorf("unknown rule %q: choose from lower, mixed, alnum, full", s)
	}
}

func parseSymbolSet(s string) (generator.SymbolSet, error) {
	switch s {
	case "safe":
		return generator.SymbolSetSafe, nil
	case "common":
		return generator.SymbolSetCommon, nil
	case "full":
		return generator.SymbolSetFull, nil
	default:
		return "", fmt.Errorf("unknown symbol set %q: choose from safe, common, full", s)
	}
}

func Execute() {
	rootCmd.Version = Version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&length, "length", "l", 16, "Password length")
	rootCmd.Flags().StringVarP(&rule, "rule", "r", "full", "Generation rule: lower, mixed, alnum, full")
	rootCmd.Flags().StringVarP(&symbolSet, "symbol-set", "s", "common", "Symbol set when using rule=full: safe (-_.), common (!@#$%^&*), full (!@#$%^&*()-_=+[]{};:,./?')")
	rootCmd.Flags().IntVarP(&count, "count", "n", 10, "Number of passwords to generate")
	rootCmd.Flags().BoolVarP(&showStrength, "show-strength", "S", false, "Show entropy-based strength level next to each password")
	rootCmd.Flags().BoolVarP(&clip, "clip", "c", false, "Copy the strongest password to the clipboard")
}
