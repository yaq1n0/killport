package process

import (
	"testing"
)

func TestParseLsofLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantPID  string
		wantPort string
		wantName string
		wantStat string
	}{
		{
			name:     "typical LISTEN line",
			line:     "node       1234  user   13u  IPv4 0x1234  0t0  TCP *:3000 (LISTEN)",
			wantPID:  "1234",
			wantPort: "3000",
			wantName: "node",
			wantStat: "LISTEN",
		},
		{
			name:     "ESTABLISHED connection",
			line:     "curl       5678  user   5u   IPv4 0x5678  0t0  TCP 127.0.0.1:52345->127.0.0.1:3000 (ESTABLISHED)",
			wantPID:  "5678",
			wantPort: "3000",
			wantName: "curl",
			wantStat: "ESTABLISHED",
		},
		{
			name:     "process name with escaped spaces",
			line:     "Google\\x20Chrome  9999  user   100u  IPv4 0xabc  0t0  TCP *:8080 (LISTEN)",
			wantPID:  "9999",
			wantPort: "8080",
			wantName: "Google Chrome",
			wantStat: "LISTEN",
		},
		{
			name:     "IPv6 address",
			line:     "node       1234  user   13u  IPv6 0x1234  0t0  TCP [::1]:3000 (LISTEN)",
			wantPID:  "1234",
			wantPort: "3000",
			wantName: "node",
			wantStat: "LISTEN",
		},
		{
			name:    "too few fields",
			line:    "node 1234 user",
			wantNil: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantNil: true,
		},
		{
			name:    "header line",
			line:    "COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME",
			wantNil: true,
		},
		{
			name:    "no port in address",
			line:     "node       1234  user   13u  IPv4 0x1234  0t0  TCP localhost (LISTEN)",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLsofLine(tt.line)
			if tt.wantNil {
				if got != nil {
					t.Errorf("parseLsofLine(%q) = %+v, want nil", tt.line, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("parseLsofLine(%q) = nil, want non-nil", tt.line)
			}
			if got.PID != tt.wantPID {
				t.Errorf("PID = %q, want %q", got.PID, tt.wantPID)
			}
			if got.Port != tt.wantPort {
				t.Errorf("Port = %q, want %q", got.Port, tt.wantPort)
			}
			if got.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", got.Name, tt.wantName)
			}
			if got.Status != tt.wantStat {
				t.Errorf("Status = %q, want %q", got.Status, tt.wantStat)
			}
		})
	}
}
