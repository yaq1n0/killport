package output

import (
	"testing"
)

func TestTruncateName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"shorter than max", "node", 20, "node"},
		{"exact length", "abcde", 5, "abcde"},
		{"longer than max", "very-long-process-name", 10, "very-lo..."},
		{"empty string", "", 10, ""},
		{"single char max", "hello", 4, "h..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateName(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateName(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestIsListen(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{"LISTEN", true},
		{"LISTENING", true},
		{"ESTABLISHED", false},
		{"CLOSE_WAIT", false},
		{"", false},
		{"listen", false},
		{"Listen", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := isListen(tt.status)
			if got != tt.want {
				t.Errorf("isListen(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
