package parser

import (
	"errors"
	"testing"
)

func TestCSVParser_Table(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr error
	}{
		{"empty", "", 0, ErrEmptyInput},
		{"one-row", "a,b\n1,2\n", 1, nil},
		{"short-row-fills-empty", "a,b\n1\n", 1, nil},
	}

	p := CSVParser{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := p.Parse([]byte(tc.input))
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if len(got) != tc.wantLen {
				t.Fatalf("want %d, got %d", tc.wantLen, len(got))
			}
		})
	}
}

func TestJSONParser_Table(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr error // this is parser from parsing
	}{
		{"empty", "", 0, ErrEmptyInput},
		{"array", `[{"a":"1"},{"a":"2","b":3}]`, 2, nil},
		{"wrapped", `{"data":[{"x":"y"}]}`, 1, nil},
		{"invalid", `{"not_data":true}`, 0, ErrInvalidFormat},
	}

	p := JSONParser{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := p.Parse([]byte(tc.input))
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Len() != tc.wantLen {
				t.Fatalf("want %d, got %d", tc.wantLen, got.Len())
			}
		})
	}
}
