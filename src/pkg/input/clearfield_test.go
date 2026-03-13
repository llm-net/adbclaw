package input

import (
	"fmt"
	"testing"

	"github.com/llm-net/adb-claw/pkg/adb"
)

// mockCommander records shell calls and returns preset responses.
type mockCommander struct {
	calls   [][]string
	results []*adb.Result // returned in order
	idx     int
}

func (m *mockCommander) Shell(args ...string) (*adb.Result, error) {
	m.calls = append(m.calls, args)
	if m.idx >= len(m.results) {
		return &adb.Result{}, nil
	}
	r := m.results[m.idx]
	m.idx++
	return r, nil
}

func (m *mockCommander) ExecOut(args ...string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCommander) RawCommand(args ...string) (*adb.Result, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestGetSDKLevel(t *testing.T) {
	m := &mockCommander{
		results: []*adb.Result{{Stdout: "34\n"}},
	}
	level, err := GetSDKLevel(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != 34 {
		t.Errorf("GetSDKLevel() = %d, want 34", level)
	}
}

func TestClearFieldSDK31Plus(t *testing.T) {
	m := &mockCommander{
		results: []*adb.Result{
			{Stdout: "31\n"},  // getprop
			{Stdout: ""},     // keycombination
			{Stdout: ""},     // keyevent DEL
		},
	}
	method, err := ClearField(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "keycombination" {
		t.Errorf("method = %q, want %q", method, "keycombination")
	}
	if len(m.calls) != 3 {
		t.Fatalf("expected 3 shell calls, got %d", len(m.calls))
	}
	// Check keycombination call
	if m.calls[1][0] != "input" || m.calls[1][1] != "keycombination" {
		t.Errorf("expected keycombination call, got %v", m.calls[1])
	}
}

func TestClearFieldOldSDK(t *testing.T) {
	m := &mockCommander{
		results: []*adb.Result{
			{Stdout: "30\n"},  // getprop
			{Stdout: ""},     // keyevent MOVE_END
			{Stdout: ""},     // keyevent DEL DEL ...
		},
	}
	method, err := ClearField(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "repeated_del" {
		t.Errorf("method = %q, want %q", method, "repeated_del")
	}
}
