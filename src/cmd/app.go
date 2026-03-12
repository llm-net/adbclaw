package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "App management commands",
}

var appListAll bool

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		shellArgs := []string{"pm", "list", "packages"}
		if !appListAll {
			shellArgs = append(shellArgs, "-3") // third-party only
		}
		result, err := client.Shell(shellArgs...)
		if err != nil {
			writer.Fail("app list", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		var packages []string
		for _, line := range strings.Split(result.Stdout, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "package:") {
				packages = append(packages, strings.TrimPrefix(line, "package:"))
			}
		}

		writer.Success("app list", map[string]interface{}{
			"packages": packages,
			"count":    len(packages),
			"all":      appListAll,
		}, start)
		return nil
	},
}

var appCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current foreground app/activity",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		result, err := client.Shell("dumpsys", "window", "displays")
		if err != nil {
			writer.Fail("app current", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		// Parse mCurrentFocus or mFocusedApp
		var focusedWindow, focusedApp string
		for _, line := range strings.Split(result.Stdout, "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "mCurrentFocus") || strings.Contains(line, "mFocusedWindow") {
				focusedWindow = extractWindowName(line)
			}
			if strings.Contains(line, "mFocusedApp") {
				focusedApp = extractWindowName(line)
			}
		}

		// Extract package and activity
		pkg, activity := parseComponent(focusedWindow)
		if pkg == "" {
			pkg, activity = parseComponent(focusedApp)
		}

		writer.Success("app current", map[string]interface{}{
			"package":  pkg,
			"activity": activity,
			"window":   focusedWindow,
		}, start)
		return nil
	},
}

var appLaunchCmd = &cobra.Command{
	Use:   "launch <package>",
	Short: "Launch an app by package name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		pkg := args[0]

		// Check if it contains an activity (pkg/activity format)
		if strings.Contains(pkg, "/") {
			result, err := client.Shell("am", "start", "-n", pkg)
			if err != nil {
				writer.Fail("app launch", "ADB_ERROR", err.Error(), "", start)
				return nil
			}
			if result.ExitCode != 0 || strings.Contains(result.Stdout, "Error type") || strings.Contains(result.Stderr, "Exception") {
				msg := strings.TrimSpace(result.Stdout + result.Stderr)
				writer.Fail("app launch", "LAUNCH_FAILED", msg,
					"Check the package/activity name", start)
				return nil
			}
			writer.Success("app launch", map[string]interface{}{
				"package":  pkg,
				"method":   "am_start",
			}, start)
			return nil
		}

		// Use monkey to launch the default activity
		result, err := client.Shell("monkey", "-p", pkg,
			"-c", "android.intent.category.LAUNCHER", "1")
		if err != nil {
			writer.Fail("app launch", "ADB_ERROR", err.Error(), "", start)
			return nil
		}
		if strings.Contains(result.Stdout, "No activities found") {
			writer.Fail("app launch", "LAUNCH_FAILED",
				"No launcher activity found for "+pkg,
				"Check the package name: adbclaw app list", start)
			return nil
		}

		writer.Success("app launch", map[string]interface{}{
			"package": pkg,
			"method":  "monkey",
		}, start)
		return nil
	},
}

var appStopCmd = &cobra.Command{
	Use:   "stop <package>",
	Short: "Force stop an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		pkg := args[0]

		_, err := client.Shell("am", "force-stop", pkg)
		if err != nil {
			writer.Fail("app stop", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		writer.Success("app stop", map[string]interface{}{
			"package": pkg,
		}, start)
		return nil
	},
}

var appInstallReplace bool

var appInstallCmd = &cobra.Command{
	Use:   "install <apk_path>",
	Short: "Install an APK",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		apkPath := args[0]

		installArgs := []string{"install"}
		if appInstallReplace {
			installArgs = append(installArgs, "-r")
		}
		installArgs = append(installArgs, apkPath)

		writer.Verbose("installing %s", apkPath)
		result, err := client.RawCommand(installArgs...)
		if err != nil {
			writer.Fail("app install", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := strings.TrimSpace(result.Stdout + result.Stderr)
		if !strings.Contains(output, "Success") {
			writer.Fail("app install", "INSTALL_FAILED", output,
				"Check the APK path and device compatibility", start)
			return nil
		}

		writer.Success("app install", map[string]interface{}{
			"apk":     apkPath,
			"replace": appInstallReplace,
		}, start)
		return nil
	},
}

var appUninstallCmd = &cobra.Command{
	Use:   "uninstall <package>",
	Short: "Uninstall an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		pkg := args[0]

		writer.Verbose("uninstalling %s", pkg)
		result, err := client.RawCommand("uninstall", pkg)
		if err != nil {
			writer.Fail("app uninstall", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := strings.TrimSpace(result.Stdout + result.Stderr)
		if !strings.Contains(output, "Success") {
			writer.Fail("app uninstall", "UNINSTALL_FAILED", output,
				"Check the package name: adbclaw app list", start)
			return nil
		}

		writer.Success("app uninstall", map[string]interface{}{
			"package": pkg,
		}, start)
		return nil
	},
}

var appClearCmd = &cobra.Command{
	Use:   "clear <package>",
	Short: "Clear app data and cache",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		pkg := args[0]

		writer.Verbose("clearing data for %s", pkg)
		result, err := client.Shell("pm", "clear", pkg)
		if err != nil {
			writer.Fail("app clear", "ADB_ERROR", err.Error(), "", start)
			return nil
		}

		output := strings.TrimSpace(result.Stdout)
		if !strings.Contains(output, "Success") {
			writer.Fail("app clear", "CLEAR_FAILED", output,
				"Check the package name: adbclaw app list", start)
			return nil
		}

		writer.Success("app clear", map[string]interface{}{
			"package": pkg,
		}, start)
		return nil
	},
}

func init() {
	appListCmd.Flags().BoolVar(&appListAll, "all", false, "Include system apps")
	appInstallCmd.Flags().BoolVar(&appInstallReplace, "replace", false, "Replace existing app (-r)")

	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appCurrentCmd)
	appCmd.AddCommand(appLaunchCmd)
	appCmd.AddCommand(appStopCmd)
	appCmd.AddCommand(appInstallCmd)
	appCmd.AddCommand(appUninstallCmd)
	appCmd.AddCommand(appClearCmd)
	rootCmd.AddCommand(appCmd)
}

// extractWindowName extracts the window/component name from a dumpsys line.
func extractWindowName(line string) string {
	// Lines look like: "mCurrentFocus=Window{abc123 u0 com.example/.MainActivity}"
	// or "mFocusedApp=AppWindowToken{...} ... com.example/.MainActivity"
	if idx := strings.LastIndex(line, " "); idx >= 0 {
		name := strings.TrimRight(line[idx+1:], "}")
		return name
	}
	return strings.TrimSpace(line)
}

// parseComponent extracts package and activity from "com.example/.MainActivity" format.
func parseComponent(component string) (pkg, activity string) {
	component = strings.TrimSpace(component)
	if component == "" {
		return "", ""
	}
	parts := strings.SplitN(component, "/", 2)
	pkg = parts[0]
	if len(parts) > 1 {
		activity = parts[1]
		// Expand shorthand: ".MainActivity" → "com.example.MainActivity"
		if strings.HasPrefix(activity, ".") {
			activity = pkg + activity
		}
	}
	return pkg, activity
}
