package input

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/llm-net/adb-claw/pkg/adb"
)

var wmSizeRe = regexp.MustCompile(`(\d+)x(\d+)`)

// GetScreenSize returns the physical screen width and height.
func GetScreenSize(cmd adb.Commander) (width, height int, err error) {
	result, err := cmd.Shell("wm", "size")
	if err != nil {
		return 0, 0, fmt.Errorf("wm size failed: %w", err)
	}
	// Parse "Physical size: 1080x2400"
	for _, line := range strings.Split(result.Stdout, "\n") {
		if strings.Contains(line, "Physical size") {
			m := wmSizeRe.FindStringSubmatch(line)
			if len(m) == 3 {
				w, _ := strconv.Atoi(m[1])
				h, _ := strconv.Atoi(m[2])
				return w, h, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("could not parse screen size from: %s", strings.TrimSpace(result.Stdout))
}

// ScrollDirection calculates swipe coordinates for a scroll direction.
// Returns (x1, y1, x2, y2) for the swipe.
func ScrollDirection(screenW, screenH, distance int, direction string) (x1, y1, x2, y2 int, err error) {
	centerX := screenW / 2
	centerY := screenH / 2

	switch strings.ToLower(direction) {
	case "down":
		// "scroll down" = see content below = swipe bottom→top (finger drags up)
		y1 = centerY + distance/2
		y2 = centerY - distance/2
		x1, x2 = centerX, centerX
	case "up":
		// "scroll up" = see content above = swipe top→bottom (finger drags down)
		y1 = centerY - distance/2
		y2 = centerY + distance/2
		x1, x2 = centerX, centerX
	case "right":
		// "scroll right" = see content to the right = swipe right→left (finger drags left)
		x1 = centerX + distance/2
		x2 = centerX - distance/2
		y1, y2 = centerY, centerY
	case "left":
		// "scroll left" = see content to the left = swipe left→right (finger drags right)
		x1 = centerX - distance/2
		x2 = centerX + distance/2
		y1, y2 = centerY, centerY
	default:
		return 0, 0, 0, 0, fmt.Errorf("invalid direction %q, use: up, down, left, right", direction)
	}
	return x1, y1, x2, y2, nil
}

// ScrollInBounds calculates swipe coordinates within element bounds.
func ScrollInBounds(left, top, right, bottom, distance int, direction string) (x1, y1, x2, y2 int, err error) {
	centerX := (left + right) / 2
	centerY := (top + bottom) / 2
	boundsH := bottom - top
	boundsW := right - left

	if distance == 0 {
		// Default to 60% of the element dimension
		switch strings.ToLower(direction) {
		case "up", "down":
			distance = boundsH * 60 / 100
		case "left", "right":
			distance = boundsW * 60 / 100
		}
	}

	// Pass centerX*2, centerY*2 as "virtual screen dimensions" so that
	// ScrollDirection computes center = (centerX*2)/2 = centerX, centering the swipe on the element.
	return ScrollDirection(centerX*2, centerY*2, distance, direction)
}
