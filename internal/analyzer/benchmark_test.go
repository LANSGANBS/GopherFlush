package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkAnalyzeSmallFile(b *testing.B) {
	analyzer := createTestAnalyzer()
	testFile := filepath.Join("testdata", "small.go")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.AnalyzeFile(testFile)
	}
}

func BenchmarkAnalyzeMediumFile(b *testing.B) {
	analyzer := createTestAnalyzer()
	testFile := filepath.Join("testdata", "medium.go")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.AnalyzeFile(testFile)
	}
}

func BenchmarkAnalyzeLargeFile(b *testing.B) {
	analyzer := createTestAnalyzer()
	testFile := filepath.Join("testdata", "large.go")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.AnalyzeFile(testFile)
	}
}

func BenchmarkPatternMatching(b *testing.B) {
	patterns := NewCPatterns()
	testLine := "static inline int calculate_sum(int a, int b) {"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = patterns.FuncPattern.FindStringSubmatch(testLine)
	}
}

func BenchmarkBracketMatching(b *testing.B) {
	analyzer := NewTextAnalyzer(LanguageC)
	lines := generateTestCode(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.findBracketEndEnhanced(lines, 0)
	}
}

func BenchmarkNormalizeText(b *testing.B) {
	content := generateLongString(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = normalizeTextFast(content)
	}
}

func BenchmarkCalculateHash(b *testing.B) {
	content := generateLongString(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateHashFast(content)
	}
}

func BenchmarkConcurrentAnalysis(b *testing.B) {
	analyzer := createTestAnalyzer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.Analyze("./testdata")
	}
}

func generateTestCode(lines int) []string {
	result := make([]string, lines)
	for i := 0; i < lines; i++ {
		if i == 0 {
			result[i] = "int main() {"
		} else if i == lines-1 {
			result[i] = "}"
		} else {
			result[i] = "    int x = 1;"
		}
	}
	return result
}

func generateLongString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte('a' + (i % 26))
	}
	return string(result)
}

func createTestAnalyzer() *Analyzer {
	_ = NewTestPatternRegistry()
	return &Analyzer{
		registry: nil,
		parser:   NewParser(),
	}
}

func NewTestPatternRegistry() *PatternRegistry {
	return NewPatternRegistry()
}

func TestMain(m *testing.M) {
	os.MkdirAll("testdata", 0755)

	smallCode := ""
	for i := 0; i < 100; i++ {
		smallCode += "func smallFunc" + string(rune('A'+i%26)) + "() {}\n"
	}
	os.WriteFile("testdata/small.go", []byte(smallCode), 0644)

	mediumCode := ""
	for i := 0; i < 500; i++ {
		mediumCode += "func mediumFunc" + string(rune('A'+i%26)) + string(rune('0'+i%10)) + "() {\n"
		mediumCode += "    x := " + string(rune('0'+i%10)) + "\n"
		mediumCode += "    return x\n"
		mediumCode += "}\n"
	}
	os.WriteFile("testdata/medium.go", []byte(mediumCode), 0644)

	largeCode := ""
	for i := 0; i < 2000; i++ {
		largeCode += "func largeFunc" + string(rune('A'+i%26)) + string(rune('0'+i%10)) + "() {\n"
		for j := 0; j < 5; j++ {
			largeCode += "    x" + string(rune('a'+j)) + " := " + string(rune('0'+j%10)) + "\n"
		}
		largeCode += "    return x\n"
		largeCode += "}\n"
	}
	os.WriteFile("testdata/large.go", []byte(largeCode), 0644)

	code := m.Run()

	os.RemoveAll("testdata")
	os.Exit(code)
}
