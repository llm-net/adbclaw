package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Response is the standard JSON envelope for all command outputs.
type Response struct {
	OK         bool        `json:"ok"`
	Command    string      `json:"command"`
	Data       interface{} `json:"data,omitempty"`
	Error      *ErrorInfo  `json:"error,omitempty"`
	DurationMs int64       `json:"duration_ms"`
	Timestamp  string      `json:"timestamp"`
}

// ErrorInfo describes a command error.
type ErrorInfo struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// Writer handles formatted output to a destination.
type Writer struct {
	out       io.Writer
	format    string // "json", "text", "quiet"
	verbose   bool
	HasFailed bool // set to true when Fail() is called; callers can check for exit code
}

// NewWriter creates a Writer with the given format.
func NewWriter(format string, verbose bool) *Writer {
	return &Writer{
		out:     os.Stdout,
		format:  format,
		verbose: verbose,
	}
}

// Success writes a successful response.
func (w *Writer) Success(command string, data interface{}, start time.Time) {
	resp := Response{
		OK:         true,
		Command:    command,
		Data:       data,
		DurationMs: time.Since(start).Milliseconds(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}
	w.write(resp)
}

// Fail writes an error response and marks the writer as failed.
func (w *Writer) Fail(command string, code, message, suggestion string, start time.Time) {
	w.HasFailed = true
	resp := Response{
		OK:      false,
		Command: command,
		Error: &ErrorInfo{
			Code:       code,
			Message:    message,
			Suggestion: suggestion,
		},
		DurationMs: time.Since(start).Milliseconds(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}
	w.write(resp)
}

// Verbose prints a debug message to stderr if verbose mode is on.
func (w *Writer) Verbose(format string, args ...interface{}) {
	if w.verbose {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}

func (w *Writer) write(resp Response) {
	switch w.format {
	case "quiet":
		if !resp.OK && resp.Error != nil {
			fmt.Fprintln(os.Stderr, resp.Error.Message)
		}
	case "text":
		if resp.OK {
			if resp.Data != nil {
				w.writeTextData(resp.Data)
			}
		} else if resp.Error != nil {
			fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", resp.Error.Code, resp.Error.Message)
			if resp.Error.Suggestion != "" {
				fmt.Fprintf(os.Stderr, "Hint: %s\n", resp.Error.Suggestion)
			}
		}
	default: // json
		enc := json.NewEncoder(w.out)
		enc.SetIndent("", "  ")
		enc.Encode(resp)
	}
}

// writeTextData writes data in a compact, human-readable format.
// Falls back to JSON for complex nested structures.
func (w *Writer) writeTextData(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			switch inner := val.(type) {
			case string:
				if len(inner) > 200 {
					// Truncate long strings (e.g. base64 screenshots)
					fmt.Fprintf(w.out, "%s: <%d chars>\n", key, len(inner))
				} else {
					fmt.Fprintf(w.out, "%s: %s\n", key, inner)
				}
			case []interface{}:
				fmt.Fprintf(w.out, "%s: (%d items)\n", key, len(inner))
				for _, item := range inner {
					b, _ := json.Marshal(item)
					fmt.Fprintf(w.out, "  %s\n", string(b))
				}
			default:
				fmt.Fprintf(w.out, "%s: %v\n", key, val)
			}
		}
	default:
		// Fallback to JSON for complex types
		b, _ := json.MarshalIndent(data, "", "  ")
		fmt.Fprintln(w.out, string(b))
	}
}
