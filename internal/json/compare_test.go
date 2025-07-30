package json

import (
	"testing"
)

func BenchmarkCompareJSONFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareJSONFiles("/Users/thomas/go/poc/json/tree.json", "/Users/thomas/go/poc/json/tree2.json")
	}
}
