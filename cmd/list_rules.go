package cmd

import (
	"fmt"

	"github.com/kawana77b/passgen/internal/generator"
	"github.com/spf13/cobra"
)

var listRulesCmd = &cobra.Command{
	Use:   "list-rules",
	Short: "List available generation rules and symbol sets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rules (-r / --rule):")
		for _, r := range generator.Rules {
			fmt.Printf("  %-8s  %-24s  %s\n", r.Name, r.Charset, r.Description)
		}

		fmt.Println()
		fmt.Println("Symbol sets (-s / --symbol-set, used with --rule=full):")
		for _, s := range generator.SymbolSets {
			fmt.Printf("  %s  %-28s  %s\n", s.Name, s.Charset, s.Desc)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRulesCmd)
}
