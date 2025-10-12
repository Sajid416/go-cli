package parser

import "testing"

// Fuzz to ensure no panics and parser gracefully handles arbitrary input.
func FuzzJSONParser_NoPanic(f *testing.F) {
	f.Add(`[]`)
	f.Add(`[{"a":"1"}]`)
	f.Add(`{"data":[{"a":"1"}]}`)
	f.Add(`{"data":[]}`)
	f.Add(`{"x":1}`)

	p := JSONParser{}
	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic: %v on input %q", r, s)
			}
		}()
		_, _ = p.Parse([]byte(s))
	})
}
