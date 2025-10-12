package main

import "testing"

func TestConfigOutput(t *testing.T) {
	tests := []struct {
		name string
		env  string
	}{
		{"local", "dev"},
		{"prod", "prod"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env == "" {
				t.Errorf("env should not empty")
			}
		})
	}

}
