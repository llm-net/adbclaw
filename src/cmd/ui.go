package cmd

import (
	"time"

	"github.com/llm-net/adb-claw/pkg/observe"
	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "UI inspection commands",
}

var uiTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Dump UI element tree with indexed elements",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		writer.Verbose("dumping UI tree")

		tree, err := observe.DumpUITree(client)
		if err != nil {
			writer.Fail("ui tree", "UI_DUMP_FAILED", err.Error(),
				"Try again — uiautomator dump can fail during animations", start)
			return nil
		}

		writer.Success("ui tree", map[string]interface{}{
			"elements": tree.Elements,
			"count":    len(tree.Elements),
		}, start)
		return nil
	},
}

var (
	uiFindText  string
	uiFindID    string
	uiFindIndex int
)

var uiFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find UI elements by text, resource-id, or index",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		writer.Verbose("finding UI elements")

		tree, err := observe.DumpUITree(client)
		if err != nil {
			writer.Fail("ui find", "UI_DUMP_FAILED", err.Error(),
				"Try again — uiautomator dump can fail during animations", start)
			return nil
		}

		// Search by index
		if cmd.Flags().Changed("index") {
			el, err := tree.FindByIndex(uiFindIndex)
			if err != nil {
				writer.Fail("ui find", "ELEMENT_NOT_FOUND", err.Error(),
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			writer.Success("ui find", map[string]interface{}{
				"elements": []observe.Element{*el},
				"count":    1,
			}, start)
			return nil
		}

		// Search by text
		if uiFindText != "" {
			results := tree.FindByText(uiFindText)
			if len(results) == 0 {
				writer.Fail("ui find", "ELEMENT_NOT_FOUND",
					"No element found with text '"+uiFindText+"'",
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			writer.Success("ui find", map[string]interface{}{
				"elements": results,
				"count":    len(results),
			}, start)
			return nil
		}

		// Search by resource-id
		if uiFindID != "" {
			results := tree.FindByID(uiFindID)
			if len(results) == 0 {
				writer.Fail("ui find", "ELEMENT_NOT_FOUND",
					"No element found with id '"+uiFindID+"'",
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			writer.Success("ui find", map[string]interface{}{
				"elements": results,
				"count":    len(results),
			}, start)
			return nil
		}

		writer.Fail("ui find", "MISSING_ARGS",
			"Specify --text, --id, or --index to search",
			"Example: adb-claw ui find --text Login", start)
		return nil
	},
}

func init() {
	uiFindCmd.Flags().StringVar(&uiFindText, "text", "", "Search by text content")
	uiFindCmd.Flags().StringVar(&uiFindID, "id", "", "Search by resource-id")
	uiFindCmd.Flags().IntVar(&uiFindIndex, "index", -1, "Get element by index number")

	uiCmd.AddCommand(uiTreeCmd)
	uiCmd.AddCommand(uiFindCmd)
	rootCmd.AddCommand(uiCmd)
}
