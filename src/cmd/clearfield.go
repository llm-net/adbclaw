package cmd

import (
	"time"

	"github.com/llm-net/adb-claw/pkg/input"
	"github.com/spf13/cobra"
)

var (
	clearFieldIndex int
	clearFieldID    string
	clearFieldText  string
)

var clearFieldCmd = &cobra.Command{
	Use:   "clear-field",
	Short: "Clear text in the focused input field",
	Long: `Clear the text in the currently focused input field.
Can optionally tap an element first to focus it:
  - By index: adb-claw clear-field --index 3
  - By resource ID: adb-claw clear-field --id "input_name"
  - By text: adb-claw clear-field --text "Username"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		var targetInfo map[string]interface{}

		// If element selector specified, tap to focus first
		switch {
		case cmd.Flags().Changed("index"):
			el, err := resolveElementByIndex(clearFieldIndex)
			if err != nil {
				writer.Fail("clear-field", "ELEMENT_NOT_FOUND", err.Error(),
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			targetInfo = elementInfo(el)
			writer.Verbose("tapping element [%d] to focus", clearFieldIndex)
			if err := input.Tap(client, el.Center.X, el.Center.Y); err != nil {
				writer.Fail("clear-field", "TAP_FAILED", err.Error(), "", start)
				return nil
			}

		case clearFieldID != "":
			el, err := resolveElementByID(clearFieldID)
			if err != nil {
				writer.Fail("clear-field", "ELEMENT_NOT_FOUND", err.Error(),
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			targetInfo = elementInfo(el)
			writer.Verbose("tapping element id=%q to focus", clearFieldID)
			if err := input.Tap(client, el.Center.X, el.Center.Y); err != nil {
				writer.Fail("clear-field", "TAP_FAILED", err.Error(), "", start)
				return nil
			}

		case clearFieldText != "":
			el, err := resolveElementByText(clearFieldText)
			if err != nil {
				writer.Fail("clear-field", "ELEMENT_NOT_FOUND", err.Error(),
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			targetInfo = elementInfo(el)
			writer.Verbose("tapping element text=%q to focus", clearFieldText)
			if err := input.Tap(client, el.Center.X, el.Center.Y); err != nil {
				writer.Fail("clear-field", "TAP_FAILED", err.Error(), "", start)
				return nil
			}
		}

		writer.Verbose("clearing field")
		method, err := input.ClearField(client)
		if err != nil {
			writer.Fail("clear-field", "CLEAR_FAILED", err.Error(), "", start)
			return nil
		}

		data := map[string]interface{}{
			"method": method,
		}
		if targetInfo != nil {
			data["element"] = targetInfo
		}
		writer.Success("clear-field", data, start)
		return nil
	},
}

func init() {
	clearFieldCmd.Flags().IntVar(&clearFieldIndex, "index", -1, "Focus element by UI tree index before clearing")
	clearFieldCmd.Flags().StringVar(&clearFieldID, "id", "", "Focus element by resource-id before clearing")
	clearFieldCmd.Flags().StringVar(&clearFieldText, "text", "", "Focus element by text content before clearing")

	rootCmd.AddCommand(clearFieldCmd)
}
