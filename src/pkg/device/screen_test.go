package device

import (
	"fmt"
	"testing"

	"github.com/llm-net/adbclaw/pkg/adb"
)

type mockCommander struct {
	calls   [][]string
	results []*adb.Result
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

func TestGetScreenStatus(t *testing.T) {
	m := &mockCommander{
		results: []*adb.Result{
			// dumpsys power
			{Stdout: "Display Power: state=ON\nmWakefulness=Awake\n"},
			// dumpsys window (lock check)
			{Stdout: "mShowingLockscreen=false\n"},
			// dumpsys window displays (rotation)
			{Stdout: "mCurrentRotation=0\n"},
		},
	}

	status, err := GetScreenStatus(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Display != "on" {
		t.Errorf("Display = %q, want %q", status.Display, "on")
	}
	if status.Locked {
		t.Error("Locked = true, want false")
	}
	if status.Rotation != 0 {
		t.Errorf("Rotation = %d, want 0", status.Rotation)
	}
}

func TestGetScreenStatusLocked(t *testing.T) {
	m := &mockCommander{
		results: []*adb.Result{
			{Stdout: "mWakefulness=Asleep\n"},
			{Stdout: "mShowingLockscreen=true\n"},
			{Stdout: "mCurrentRotation=1\n"},
		},
	}

	status, err := GetScreenStatus(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Display != "off" {
		t.Errorf("Display = %q, want %q", status.Display, "off")
	}
	if !status.Locked {
		t.Error("Locked = false, want true")
	}
	if status.Rotation != 1 {
		t.Errorf("Rotation = %d, want 1", status.Rotation)
	}
}

func TestSetRotation(t *testing.T) {
	t.Run("auto", func(t *testing.T) {
		m := &mockCommander{
			results: []*adb.Result{
				{Stdout: ""},
			},
		}
		err := SetRotation(m, "auto")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(m.calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(m.calls))
		}
		// Should set accelerometer_rotation to 1
		args := m.calls[0]
		if args[len(args)-1] != "1" {
			t.Errorf("expected last arg '1', got %q", args[len(args)-1])
		}
	})

	t.Run("fixed", func(t *testing.T) {
		m := &mockCommander{
			results: []*adb.Result{
				{Stdout: ""}, // disable auto
				{Stdout: ""}, // set user_rotation
			},
		}
		err := SetRotation(m, "2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(m.calls) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(m.calls))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		m := &mockCommander{}
		err := SetRotation(m, "5")
		if err == nil {
			t.Error("expected error for rotation 5")
		}
	})
}
