package cmd

import (
	"time"

	"github.com/llm-net/adb-claw/pkg/device"
	"github.com/spf13/cobra"
)

var screenCmd = &cobra.Command{
	Use:   "screen",
	Short: "Screen management commands",
}

var screenStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show screen state (on/off, locked, rotation)",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		status, err := device.GetScreenStatus(client)
		if err != nil {
			writer.Fail("screen status", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		writer.Success("screen status", status, start)
		return nil
	},
}

var screenOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn the screen on (wake up)",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		if err := device.ScreenOn(client); err != nil {
			writer.Fail("screen on", "SCREEN_FAILED", err.Error(), "", start)
			return nil
		}

		writer.Success("screen on", map[string]interface{}{
			"action": "wakeup",
		}, start)
		return nil
	},
}

var screenOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn the screen off (sleep)",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		if err := device.ScreenOff(client); err != nil {
			writer.Fail("screen off", "SCREEN_FAILED", err.Error(), "", start)
			return nil
		}

		writer.Success("screen off", map[string]interface{}{
			"action": "sleep",
		}, start)
		return nil
	},
}

var screenUnlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Wake and unlock screen (no-password devices only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		if err := device.ScreenUnlock(client); err != nil {
			writer.Fail("screen unlock", "UNLOCK_FAILED", err.Error(),
				"This only works on devices without password/PIN lock", start)
			return nil
		}

		writer.Success("screen unlock", map[string]interface{}{
			"action": "wakeup_and_swipe",
		}, start)
		return nil
	},
}

var screenRotationCmd = &cobra.Command{
	Use:   "rotation <auto|0|1|2|3>",
	Short: "Set screen rotation",
	Long: `Set screen rotation mode:
  auto — Enable auto-rotation
  0    — Portrait (default)
  1    — Landscape (90° clockwise)
  2    — Reverse portrait
  3    — Reverse landscape`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		mode := args[0]

		if err := device.SetRotation(client, mode); err != nil {
			writer.Fail("screen rotation", "ROTATION_FAILED", err.Error(), "", start)
			return nil
		}

		writer.Success("screen rotation", map[string]interface{}{
			"rotation": mode,
		}, start)
		return nil
	},
}

func init() {
	screenCmd.AddCommand(screenStatusCmd)
	screenCmd.AddCommand(screenOnCmd)
	screenCmd.AddCommand(screenOffCmd)
	screenCmd.AddCommand(screenUnlockCmd)
	screenCmd.AddCommand(screenRotationCmd)
	rootCmd.AddCommand(screenCmd)
}
