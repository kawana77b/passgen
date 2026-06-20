package cmd

import (
	"fmt"

	"github.com/kawana77b/passgen/internal/uid"
	"github.com/spf13/cobra"
)

var uidCount int

var uidCmd = &cobra.Command{
	Use:   "uid <v4|v7|nanoid>",
	Short: "Generate UUIDs or NanoIDs",
	Example: `  passgen uid v4
  passgen uid v7 -n 5
  passgen uid nanoid`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"v4", "v7", "nanoid"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if uidCount < 1 {
			return fmt.Errorf("count must be at least 1")
		}

		kind := uid.Kind(args[0])
		ids, err := uid.NewN(kind, uidCount)
		if err != nil {
			return err
		}
		for _, id := range ids {
			fmt.Println(id)
		}
		return nil
	},
}

func init() {
	uidCmd.Flags().IntVarP(&uidCount, "count", "n", 10, "Number of IDs to generate")
	rootCmd.AddCommand(uidCmd)
}
