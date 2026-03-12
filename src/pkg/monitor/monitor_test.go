package monitor

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

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TextEntry
		wantErr bool
	}{
		{
			name:  "valid text entry",
			input: `{"text":"Hello world","class":"android.widget.TextView"}`,
			want:  TextEntry{Text: "Hello world", Class: "android.widget.TextView"},
		},
		{
			name:  "unicode text",
			input: `{"text":"你好世界","class":"android.widget.TextView"}`,
			want:  TextEntry{Text: "你好世界", Class: "android.widget.TextView"},
		},
		{
			name:  "empty text",
			input: `{"text":"","class":""}`,
			want:  TextEntry{Text: "", Class: ""},
		},
		{
			name:    "invalid json",
			input:   `not json`,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Text != tt.want.Text {
				t.Errorf("Text = %q, want %q", got.Text, tt.want.Text)
			}
			if got.Class != tt.want.Class {
				t.Errorf("Class = %q, want %q", got.Class, tt.want.Class)
			}
		})
	}
}
