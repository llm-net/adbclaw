package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/llm-net/adbclaw/pkg/adb"
	"github.com/llm-net/adbclaw/pkg/output"
	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time: -ldflags "-X github.com/llm-net/adbclaw/cmd.Version=v0.1.0"
var Version = "dev"

var (
	flagSerial  string
	flagOutput  string
	flagTimeout int
	flagVerbose bool

	writer *output.Writer
	client *adb.Client
)

var rootCmd = &cobra.Command{
	Use:     "adbclaw",
	Short:   "Android device control CLI for AI agents",
	Long:    "adbclaw is a CLI tool for controlling Android devices via ADB, designed for AI agent automation.",
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		writer = output.NewWriter(flagOutput, flagVerbose)
		client = adb.NewClient(flagSerial, time.Duration(flagTimeout)*time.Millisecond)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagSerial, "serial", "s", "", "Target device serial")
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "json", "Output format: json | text | quiet")
	rootCmd.PersistentFlags().IntVar(&flagTimeout, "timeout", 30000, "Command timeout in milliseconds")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Enable debug output to stderr")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if writer != nil && writer.HasFailed {
		os.Exit(1)
	}
}
