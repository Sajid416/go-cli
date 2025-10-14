package parser

import "testing"

func BenchmarkParseCSV(b *testing.B) {
	data := "a,b,c\n1,2,3\n4,5,6\n7,8,9\n"
	p := CSVParser{}
	b.ResetTimer() //feature
	for i := 0; i < b.N; i++ {
		_, _ = p.Parse([]byte(data))
	}
}
