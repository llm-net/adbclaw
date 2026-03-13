package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed skill.json
var skillJSON string

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Output the skill description (for AI agent integration)",
	Long:  "Prints the skill.json file that describes adb-claw's capabilities for AI agents.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(skillJSON)
	},
}

func init() {
	rootCmd.AddCommand(skillCmd)
}
