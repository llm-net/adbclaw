package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell <command>",
	Short: "Run a raw adb shell command",
	Long: `Execute a raw command on the device via adb shell.
The command output is captured and returned in the JSON envelope.
Example:
  adbclaw shell "ls /sdcard/"
  adbclaw shell "getprop ro.build.version.release"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		// Join all args as a single shell command string
		command := strings.Join(args, " ")

		writer.Verbose("shell: %s", command)
		result, err := client.Shell(command)
		if err != nil {
			writer.Fail("shell", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		writer.Success("shell", map[string]interface{}{
			"stdout":    strings.TrimRight(result.Stdout, "\n"),
			"stderr":    strings.TrimRight(result.Stderr, "\n"),
			"exit_code": result.ExitCode,
		}, start)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
