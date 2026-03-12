package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <uri>",
	Short: "Open a URI via Android intent (deep link, URL, etc.)",
	Long: `Open a URI using Android's ACTION_VIEW intent.
Examples:
  adbclaw open https://www.google.com
  adbclaw open myapp://path/to/screen
  adbclaw open "market://details?id=com.example"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		uri := args[0]

		writer.Verbose("opening URI: %s", uri)
		result, err := client.Shell("am", "start",
			"-a", "android.intent.action.VIEW",
			"-d", uri)
		if err != nil {
			writer.Fail("open", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := result.Stdout + result.Stderr
		// Check for am start error patterns (not URI content)
		if strings.Contains(output, "Error type") ||
			strings.Contains(output, "Error:") ||
			strings.Contains(output, "Exception") ||
			result.ExitCode != 0 {
			writer.Fail("open", "OPEN_FAILED",
				strings.TrimSpace(output),
				"Check the URI format and ensure a handler is installed", start)
			return nil
		}

		data := map[string]interface{}{
			"uri": uri,
		}

		// Try to extract component info from am start output
		// "Starting: Intent { act=... cmp=com.example/.Activity }"
		for _, line := range strings.Split(result.Stdout, "\n") {
			if strings.Contains(line, "cmp=") {
				idx := strings.Index(line, "cmp=")
				rest := line[idx+4:]
				end := strings.IndexAny(rest, " }")
				if end > 0 {
					component := rest[:end]
					pkg, activity := parseComponent(component)
					if pkg != "" {
						data["package"] = pkg
					}
					if activity != "" {
						data["activity"] = activity
					}
				}
			}
		}

		writer.Success("open", data, start)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
