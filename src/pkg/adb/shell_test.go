package adb

import (
	"testing"
	"time"
)

func TestNewClientDefaults(t *testing.T) {
	c := NewClient("", 0)
	if c.ADBPath != "adb" {
		t.Errorf("ADBPath = %q, want %q", c.ADBPath, "adb")
	}
	if c.Serial != "" {
		t.Errorf("Serial = %q, want empty", c.Serial)
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", c.Timeout, 30*time.Second)
	}
}

func TestNewClientCustom(t *testing.T) {
	c := NewClient("abc123", 5*time.Second)
	if c.Serial != "abc123" {
		t.Errorf("Serial = %q, want %q", c.Serial, "abc123")
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want %v", c.Timeout, 5*time.Second)
	}
}

func TestBaseArgsNoSerial(t *testing.T) {
	c := NewClient("", 0)
	args := c.BaseArgs()
	if args != nil {
		t.Errorf("BaseArgs() = %v, want nil", args)
	}
}

func TestBaseArgsWithSerial(t *testing.T) {
	c := NewClient("device123", 0)
	args := c.BaseArgs()
	if len(args) != 2 || args[0] != "-s" || args[1] != "device123" {
		t.Errorf("BaseArgs() = %v, want [-s device123]", args)
	}
}
