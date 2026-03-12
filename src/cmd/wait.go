package cmd

import (
	"strings"
	"time"

	"github.com/llm-net/adbclaw/pkg/observe"
	"github.com/spf13/cobra"
)

var (
	waitText     string
	waitID       string
	waitActivity string
	waitGone     bool
	waitTimeout  int
	waitInterval int
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for a UI element or activity to appear/disappear",
	Long: `Wait for a condition to be met on the device screen.
Examples:
  adbclaw wait --text "Login"                 # Wait for text to appear
  adbclaw wait --id "btn_submit"              # Wait for element by ID
  adbclaw wait --text "Loading" --gone        # Wait for text to disappear
  adbclaw wait --activity ".MainActivity"     # Wait for activity
  adbclaw wait --text "Done" --timeout 20000  # Custom timeout (20s)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		if waitText == "" && waitID == "" && waitActivity == "" {
			writer.Fail("wait", "MISSING_ARGS",
				"Specify --text, --id, or --activity",
				"Example: adbclaw wait --text \"Login\"", start)
			return nil
		}

		timeout := time.Duration(waitTimeout) * time.Millisecond
		interval := time.Duration(waitInterval) * time.Millisecond

		deadline := time.Now().Add(timeout)
		attempts := 0

		for time.Now().Before(deadline) {
			attempts++

			if waitActivity != "" {
				// Activity mode: check dumpsys window
				found, activity, err := checkActivity(waitActivity)
				if err != nil {
					writer.Verbose("activity check error (attempt %d): %v", attempts, err)
				} else if found && !waitGone {
					writer.Success("wait", map[string]interface{}{
						"condition": "activity",
						"activity":  activity,
						"gone":      false,
						"attempts":  attempts,
					}, start)
					return nil
				} else if !found && waitGone {
					writer.Success("wait", map[string]interface{}{
						"condition": "activity",
						"activity":  waitActivity,
						"gone":      true,
						"attempts":  attempts,
					}, start)
					return nil
				}
			} else {
				// Text/ID mode: check UI tree
				tree, err := observe.DumpUITree(client)
				if err != nil {
					writer.Verbose("ui dump error (attempt %d): %v", attempts, err)
				} else {
					var found bool
					var matchedElement map[string]interface{}

					if waitText != "" {
						results := tree.FindByText(waitText)
						if len(results) > 0 {
							found = true
							matchedElement = elementInfo(&results[0])
						}
					} else if waitID != "" {
						results := tree.FindByID(waitID)
						if len(results) > 0 {
							found = true
							matchedElement = elementInfo(&results[0])
						}
					}

					if found && !waitGone {
						data := map[string]interface{}{
							"condition": "element",
							"gone":      false,
							"attempts":  attempts,
						}
						if matchedElement != nil {
							data["element"] = matchedElement
						}
						writer.Success("wait", data, start)
						return nil
					} else if !found && waitGone {
						data := map[string]interface{}{
							"condition": "element",
							"gone":      true,
							"attempts":  attempts,
						}
						if waitText != "" {
							data["text"] = waitText
						}
						if waitID != "" {
							data["id"] = waitID
						}
						writer.Success("wait", data, start)
						return nil
					}
				}
			}

			time.Sleep(interval)
		}

		// Timeout
		condition := "element"
		if waitActivity != "" {
			condition = "activity"
		}
		action := "appear"
		if waitGone {
			action = "disappear"
		}
		writer.Fail("wait", "WAIT_TIMEOUT",
			condition+" did not "+action+" within "+timeout.String(),
			"Try increasing --timeout or check the condition", start)
		return nil
	},
}

func init() {
	waitCmd.Flags().StringVar(&waitText, "text", "", "Wait for element with this text")
	waitCmd.Flags().StringVar(&waitID, "id", "", "Wait for element with this resource-id")
	waitCmd.Flags().StringVar(&waitActivity, "activity", "", "Wait for this activity to be in foreground")
	waitCmd.Flags().BoolVar(&waitGone, "gone", false, "Wait for element/activity to disappear")
	waitCmd.Flags().IntVar(&waitTimeout, "timeout", 10000, "Timeout in milliseconds")
	waitCmd.Flags().IntVar(&waitInterval, "interval", 800, "Poll interval in milliseconds")

	rootCmd.AddCommand(waitCmd)
}

// checkActivity checks if the given activity is currently in the foreground.
func checkActivity(activity string) (bool, string, error) {
	result, err := client.Shell("dumpsys", "window", "displays")
	if err != nil {
		return false, "", err
	}

	for _, line := range strings.Split(result.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "mCurrentFocus") || strings.Contains(line, "mFocusedWindow") || strings.Contains(line, "mFocusedApp") {
			if strings.Contains(line, activity) {
				windowName := extractWindowName(line)
				return true, windowName, nil
			}
		}
	}
	return false, "", nil
}
