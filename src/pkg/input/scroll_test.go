package input

import (
	"testing"
)

func TestScrollDirection(t *testing.T) {
	// Screen: 1080x2400
	tests := []struct {
		direction      string
		wantY1GreaterY2 bool // true if y1 > y2 (swipe from bottom to top)
		wantX1GreaterX2 bool // true if x1 > x2 (swipe from right to left)
	}{
		// "scroll down" = see content below = swipe bottom→top = y1 > y2
		{"down", true, false},
		// "scroll up" = see content above = swipe top→bottom = y1 < y2
		{"up", false, false},
		// "scroll right" = see content right = swipe right→left = x1 > x2
		{"right", false, true},
		// "scroll left" = see content left = swipe left→right = x1 < x2
		{"left", false, false},
	}

	for _, tt := range tests {
		x1, y1, x2, y2, err := ScrollDirection(1080, 2400, 1440, tt.direction)
		if err != nil {
			t.Fatalf("ScrollDirection(%q) error: %v", tt.direction, err)
		}

		if tt.direction == "down" || tt.direction == "up" {
			if x1 != x2 {
				t.Errorf("ScrollDirection(%q): x1=%d != x2=%d, vertical scroll should have same x", tt.direction, x1, x2)
			}
			if (y1 > y2) != tt.wantY1GreaterY2 {
				t.Errorf("ScrollDirection(%q): y1=%d, y2=%d, wantY1>Y2=%v", tt.direction, y1, y2, tt.wantY1GreaterY2)
			}
		}

		if tt.direction == "left" || tt.direction == "right" {
			if y1 != y2 {
				t.Errorf("ScrollDirection(%q): y1=%d != y2=%d, horizontal scroll should have same y", tt.direction, y1, y2)
			}
			if (x1 > x2) != tt.wantX1GreaterX2 {
				t.Errorf("ScrollDirection(%q): x1=%d, x2=%d, wantX1>X2=%v", tt.direction, x1, x2, tt.wantX1GreaterX2)
			}
		}

		_ = x1
		_ = y1
		_ = x2
		_ = y2
	}
}

func TestScrollDirectionInvalid(t *testing.T) {
	_, _, _, _, err := ScrollDirection(1080, 2400, 1440, "diagonal")
	if err == nil {
		t.Error("ScrollDirection(\"diagonal\") should return error")
	}
}

func TestScrollDirectionCentered(t *testing.T) {
	// Verify the scroll is centered on the screen
	x1, y1, x2, y2, err := ScrollDirection(1080, 2400, 1000, "down")
	if err != nil {
		t.Fatal(err)
	}
	// Center should be at 540, 1200
	midX := (x1 + x2) / 2
	midY := (y1 + y2) / 2
	if midX != 540 {
		t.Errorf("midX = %d, want 540", midX)
	}
	if midY != 1200 {
		t.Errorf("midY = %d, want 1200", midY)
	}
}

func TestScrollInBounds(t *testing.T) {
	// Element at left=100, top=400, right=900, bottom=1600
	// Center should be (500, 1000), height=1200, distance=720 (60%)
	x1, y1, x2, y2, err := ScrollInBounds(100, 400, 900, 1600, 0, "down")
	if err != nil {
		t.Fatal(err)
	}

	// Swipe should be centered at element center (500, 1000)
	midX := (x1 + x2) / 2
	midY := (y1 + y2) / 2
	if midX != 500 {
		t.Errorf("midX = %d, want 500", midX)
	}
	if midY != 1000 {
		t.Errorf("midY = %d, want 1000", midY)
	}

	// "scroll down" should have y1 > y2 (swipe bottom to top)
	if y1 <= y2 {
		t.Errorf("scroll down: y1=%d should be > y2=%d", y1, y2)
	}

	// All coordinates should be within element bounds
	if y1 < 400 || y1 > 1600 {
		t.Errorf("y1=%d out of bounds [400, 1600]", y1)
	}
	if y2 < 400 || y2 > 1600 {
		t.Errorf("y2=%d out of bounds [400, 1600]", y2)
	}
}

func TestScrollInBoundsOffCenter(t *testing.T) {
	// Element NOT centered on screen: left=600, top=800, right=1000, bottom=2000
	// Center: (800, 1400), height=1200, distance=720 (60%)
	x1, y1, x2, y2, err := ScrollInBounds(600, 800, 1000, 2000, 0, "down")
	if err != nil {
		t.Fatal(err)
	}

	midX := (x1 + x2) / 2
	midY := (y1 + y2) / 2
	if midX != 800 {
		t.Errorf("midX = %d, want 800", midX)
	}
	if midY != 1400 {
		t.Errorf("midY = %d, want 1400", midY)
	}
}
