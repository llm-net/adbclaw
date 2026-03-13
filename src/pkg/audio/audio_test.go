package audio

import (
	"testing"
)

func TestDEXEmbed(t *testing.T) {
	if len(dexData) == 0 {
		t.Fatal("embedded DEX data is empty")
	}
	// DEX files start with "dex\n"
	if string(dexData[:4]) != "dex\n" {
		t.Fatalf("embedded DEX has wrong magic: got %q, want %q", string(dexData[:4]), "dex\n")
	}
}
