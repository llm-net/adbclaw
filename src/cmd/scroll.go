package cmd

import (
	"time"

	"github.com/llm-net/adb-claw/pkg/input"
	"github.com/spf13/cobra"
)

var (
	scrollIndex    int
	scrollPages    int
	scrollDistance  int
	scrollDuration int
)

var scrollCmd = &cobra.Command{
	Use:   "scroll <direction>",
	Short: "Scroll the screen or a scrollable element",
	Long: `Scroll in a direction: up, down, left, right.
Examples:
  adb-claw scroll down                  # Scroll down one screen
  adb-claw scroll up --pages 3          # Scroll up 3 screens
  adb-claw scroll down --index 5        # Scroll within element at index 5
  adb-claw scroll left --distance 500   # Scroll left by 500 pixels`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		direction := args[0]

		pages := scrollPages
		if pages < 1 {
			pages = 1
		}

		// Pre-compute swipe coordinates once (element/screen size won't change between pages)
		var x1, y1, x2, y2 int
		var targetInfo map[string]interface{}

		if cmd.Flags().Changed("index") {
			// Scroll within a specific element — resolve once
			el, err := resolveElementByIndex(scrollIndex)
			if err != nil {
				writer.Fail("scroll", "ELEMENT_NOT_FOUND", err.Error(),
					"Use 'adb-claw ui tree' to see available elements", start)
				return nil
			}
			targetInfo = elementInfo(el)

			dist := scrollDistance
			if dist == 0 {
				switch direction {
				case "up", "down":
					dist = (el.Bounds.Bottom - el.Bounds.Top) * 60 / 100
				case "left", "right":
					dist = (el.Bounds.Right - el.Bounds.Left) * 60 / 100
				}
			}

			x1, y1, x2, y2, err = input.ScrollInBounds(
				el.Bounds.Left, el.Bounds.Top, el.Bounds.Right, el.Bounds.Bottom,
				dist, direction)
			if err != nil {
				writer.Fail("scroll", "SCROLL_FAILED", err.Error(), "", start)
				return nil
			}
		} else {
			// Full screen scroll — get screen size once
			screenW, screenH, err := input.GetScreenSize(client)
			if err != nil {
				writer.Fail("scroll", "SCREEN_SIZE_FAILED", err.Error(),
					"Ensure device is connected", start)
				return nil
			}

			dist := scrollDistance
			if dist == 0 {
				switch direction {
				case "up", "down":
					dist = screenH * 60 / 100
				case "left", "right":
					dist = screenW * 60 / 100
				}
			}

			x1, y1, x2, y2, err = input.ScrollDirection(screenW, screenH, dist, direction)
			if err != nil {
				writer.Fail("scroll", "INVALID_DIRECTION", err.Error(), "", start)
				return nil
			}
		}

		// Execute scrolls
		var totalScrolled int
		for i := 0; i < pages; i++ {
			writer.Verbose("scroll %s page %d/%d: (%d,%d) → (%d,%d)", direction, i+1, pages, x1, y1, x2, y2)
			if err := input.Swipe(client, x1, y1, x2, y2, scrollDuration); err != nil {
				writer.Fail("scroll", "SWIPE_FAILED", err.Error(), "", start)
				return nil
			}

			// Calculate distance scrolled
			dy := y1 - y2
			dx := x1 - x2
			if dy < 0 {
				dy = -dy
			}
			if dx < 0 {
				dx = -dx
			}
			if dy > dx {
				totalScrolled += dy
			} else {
				totalScrolled += dx
			}

			// Pause between pages (except last)
			if i < pages-1 {
				time.Sleep(300 * time.Millisecond)
			}
		}

		data := map[string]interface{}{
			"direction":       direction,
			"pages":           pages,
			"distance_pixels": totalScrolled,
			"method":          "adb_swipe",
		}
		if targetInfo != nil {
			data["element"] = targetInfo
		}
		writer.Success("scroll", data, start)
		return nil
	},
}

func init() {
	scrollCmd.Flags().IntVar(&scrollIndex, "index", -1, "Scroll within element by UI tree index")
	scrollCmd.Flags().IntVar(&scrollPages, "pages", 1, "Number of pages to scroll")
	scrollCmd.Flags().IntVar(&scrollDistance, "distance", 0, "Scroll distance in pixels (0 = auto)")
	scrollCmd.Flags().IntVar(&scrollDuration, "duration", 300, "Swipe duration in ms")

	rootCmd.AddCommand(scrollCmd)
}
