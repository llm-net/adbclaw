package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File transfer commands (push/pull)",
}

var filePushCmd = &cobra.Command{
	Use:   "push <local> <remote>",
	Short: "Push a local file to the device",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		local := args[0]
		remote := args[1]

		writer.Verbose("pushing %s → %s", local, remote)
		result, err := client.RawCommand("push", local, remote)
		if err != nil {
			writer.Fail("file push", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := strings.TrimSpace(result.Stdout + result.Stderr)
		if result.ExitCode != 0 {
			writer.Fail("file push", "PUSH_FAILED", output,
				"Check that the local file exists and the remote path is writable", start)
			return nil
		}

		writer.Success("file push", map[string]interface{}{
			"local":  local,
			"remote": remote,
			"detail": output,
		}, start)
		return nil
	},
}

var filePullCmd = &cobra.Command{
	Use:   "pull <remote> <local>",
	Short: "Pull a file from the device",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		remote := args[0]
		local := args[1]

		writer.Verbose("pulling %s → %s", remote, local)
		result, err := client.RawCommand("pull", remote, local)
		if err != nil {
			writer.Fail("file pull", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := strings.TrimSpace(result.Stdout + result.Stderr)
		if result.ExitCode != 0 {
			writer.Fail("file pull", "PULL_FAILED", output,
				"Check that the remote file exists", start)
			return nil
		}

		writer.Success("file pull", map[string]interface{}{
			"remote": remote,
			"local":  local,
			"detail": output,
		}, start)
		return nil
	},
}

func init() {
	fileCmd.AddCommand(filePushCmd)
	fileCmd.AddCommand(filePullCmd)
	rootCmd.AddCommand(fileCmd)
}
