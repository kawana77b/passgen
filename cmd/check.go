package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/kawana77b/passgen/internal/strength"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <password>",
	Short: "Check the entropy and strength level of a password",
	Example: `  passgen check "myPassword123!"
  passgen check "correct-horse-battery"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		password := args[0]
		entropy := strength.Entropy(password)
		level := strength.Check(password)

		out := struct {
			Password string  `json:"password"`
			Length   int     `json:"length"`
			Entropy  float64 `json:"entropy_bits"`
			Level    string  `json:"level"`
		}{
			Password: password,
			Length:   len(password),
			Entropy:  round2(entropy),
			Level:    level.String(),
		}

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, buf.String())
		return nil
	},
}

// round2 rounds f to 2 decimal places.
func round2(f float64) float64 {
	return math.Round(f*100) / 100
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
