package observe

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/llm-net/adb-claw/pkg/adb"

	"golang.org/x/image/draw"
)

// ScreenshotResult holds the raw PNG data and optional base64 encoding.
type ScreenshotResult struct {
	Format string `json:"format"`
	Base64 string `json:"base64,omitempty"`
	Path   string `json:"path,omitempty"`
	Size   int    `json:"size_bytes"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// TakeScreenshot captures the device screen via "adb exec-out screencap -p".
// Returns raw PNG bytes. If maxWidth > 0, the image is downscaled proportionally.
func TakeScreenshot(cmd adb.Commander, maxWidth int) ([]byte, error) {
	data, err := cmd.ExecOut("screencap", "-p")
	if err != nil {
		return nil, fmt.Errorf("screencap failed: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("screencap returned empty data")
	}
	// Validate PNG header
	if len(data) < 8 || string(data[1:4]) != "PNG" {
		return nil, fmt.Errorf("screencap returned invalid PNG data (%d bytes)", len(data))
	}
	if maxWidth > 0 {
		data, err = downscalePNG(data, maxWidth)
		if err != nil {
			return nil, fmt.Errorf("screenshot downscale failed: %w", err)
		}
	}
	return data, nil
}

// ScreenshotAsBase64 captures a screenshot and returns a ScreenshotResult with base64 encoding.
// If maxWidth > 0, the image is downscaled proportionally.
func ScreenshotAsBase64(cmd adb.Commander, maxWidth int) (*ScreenshotResult, error) {
	data, err := TakeScreenshot(cmd, maxWidth)
	if err != nil {
		return nil, err
	}
	return &ScreenshotResult{
		Format: "png",
		Base64: base64.StdEncoding.EncodeToString(data),
		Size:   len(data),
	}, nil
}

// downscalePNG decodes a PNG, scales it to fit within maxWidth, and re-encodes.
// If the image is already smaller than maxWidth, it is returned as-is.
func downscalePNG(data []byte, maxWidth int) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	srcBounds := src.Bounds()
	origW := srcBounds.Dx()
	origH := srcBounds.Dy()

	if origW <= maxWidth {
		return data, nil
	}

	newW := maxWidth
	newH := origH * maxWidth / origW

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
