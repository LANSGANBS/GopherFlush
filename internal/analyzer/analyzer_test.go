package analyzer

import (
	"testing"
)

func TestDetectorPatterns(t *testing.T) {
	t.Run("CPatterns", func(t *testing.T) {
		patterns := NewCPatterns()
		if patterns.FuncPattern == nil {
			t.Error("FuncPattern should not be nil")
		}
		if patterns.GlobalPattern == nil {
			t.Error("GlobalPattern should not be nil")
		}
	})

	t.Run("PythonPatterns", func(t *testing.T) {
		patterns := NewPythonPatterns()
		if patterns.FuncPattern == nil {
			t.Error("FuncPattern should not be nil")
		}
		if patterns.GlobalPattern == nil {
			t.Error("GlobalPattern should not be nil")
		}
	})

	t.Run("JavaPatterns", func(t *testing.T) {
		patterns := NewJavaPatterns()
		if patterns.FuncPattern == nil {
			t.Error("FuncPattern should not be nil")
		}
		if patterns.GlobalPattern == nil {
			t.Error("GlobalPattern should not be nil")
		}
	})
}

func TestPatternRegistry(t *testing.T) {
	registry := NewPatternRegistry()

	tests := []struct {
		lang     Language
		expected bool
	}{
		{LanguageC, true},
		{LanguageCPP, true},
		{LanguagePython, true},
		{LanguageJava, true},
		{LanguageUnknown, false},
	}

	for _, tt := range tests {
		patterns := registry.Get(tt.lang)
		if tt.expected && patterns == nil {
			t.Errorf("Expected patterns for language %v, got nil", tt.lang)
		}
		if !tt.expected && patterns != nil {
			t.Errorf("Expected nil for language %v, got patterns", tt.lang)
		}
	}
}

func TestCFunctionDetection(t *testing.T) {
	analyzer := NewTextAnalyzer(LanguageC)

	tests := []struct {
		name     string
		code     []string
		expected int
	}{
		{
			name: "simple function",
			code: []string{
				"int main() {",
				"    return 0;",
				"}",
			},
			expected: 1,
		},
		{
			name: "function with modifiers",
			code: []string{
				"static inline int helper(void) {",
				"    return 42;",
				"}",
			},
			expected: 1,
		},
		{
			name: "function with attributes",
			code: []string{
				"__attribute__((deprecated)) void old_func() {",
				"}",
			},
			expected: 0,
		},
		{
			name: "nested functions",
			code: []string{
				"void outer() {",
				"    int x = 1;",
				"    if (x) {",
				"        // nested block",
				"    }",
				"}",
			},
			expected: 1,
		},
		{
			name: "skip preprocessor",
			code: []string{
				"#include <stdio.h>",
				"#define MAX 100",
				"int func() {",
				"    return MAX;",
				"}",
			},
			expected: 1,
		},
		{
			name: "skip comments",
			code: []string{
				"// This is a comment",
				"int func() { /* inline comment */",
				"    return 0;",
				"}",
			},
			expected: 1,
		},
		{
			name: "function with string containing braces",
			code: []string{
				`void print_json() {`,
				`    printf("{\"key\": \"value\"}");`,
				`}`,
			},
			expected: 2,
		},
		{
			name: "multi-line function declaration",
			code: []string{
				"int very_long_function(",
				"    int arg1,",
				"    int arg2) {",
				"    return arg1 + arg2;",
				"}",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions := analyzer.detectCFunctions(tt.code)
			if len(functions) != tt.expected {
				t.Errorf("Expected %d functions, got %d", tt.expected, len(functions))
				for _, f := range functions {
					t.Logf("Found function: %s at line %d-%d", f.Name, f.StartLine, f.EndLine)
				}
			}
		})
	}
}

func TestPythonFunctionDetection(t *testing.T) {
	analyzer := NewTextAnalyzer(LanguagePython)

	tests := []struct {
		name     string
		code     []string
		expected int
	}{
		{
			name: "simple function",
			code: []string{
				"def hello():",
				"    print('hello')",
			},
			expected: 1,
		},
		{
			name: "function with decorator",
			code: []string{
				"@staticmethod",
				"@cache",
				"def cached_func():",
				"    return expensive_calc()",
			},
			expected: 1,
		},
		{
			name: "decorator with arguments",
			code: []string{
				"@route('/api/users')",
				"def get_users():",
				"    return []",
			},
			expected: 1,
		},
		{
			name: "nested function",
			code: []string{
				"def outer():",
				"    def inner():",
				"        pass",
				"    return inner",
			},
			expected: 2,
		},
		{
			name: "class with methods",
			code: []string{
				"class MyClass:",
				"    def __init__(self):",
				"        pass",
				"",
				"    def method(self):",
				"        pass",
			},
			expected: 2,
		},
		{
			name: "function with complex body",
			code: []string{
				"def complex_func():",
				"    if True:",
				"        for i in range(10):",
				"            if i > 5:",
				"                break",
				"    return None",
				"",
				"def next_func():",
				"    pass",
			},
			expected: 2,
		},
		{
			name: "async function",
			code: []string{
				"async def fetch_data():",
				"    await some_io()",
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions := analyzer.detectPythonFunctions(tt.code)
			if len(functions) != tt.expected {
				t.Errorf("Expected %d functions, got %d", tt.expected, len(functions))
				for _, f := range functions {
					t.Logf("Found function: %s at line %d-%d", f.Name, f.StartLine, f.EndLine)
				}
			}
		})
	}
}

func TestJavaFunctionDetection(t *testing.T) {
	analyzer := NewTextAnalyzer(LanguageJava)

	tests := []struct {
		name     string
		code     []string
		expected int
	}{
		{
			name: "simple method",
			code: []string{
				"public void doSomething() {",
				"    System.out.println(\"hello\");",
				"}",
			},
			expected: 1,
		},
		{
			name: "method with generics",
			code: []string{
				"public <T> List<T> getList() {",
				"    return new ArrayList<>();",
				"}",
			},
			expected: 1,
		},
		{
			name: "method with annotations",
			code: []string{
				"@Override",
				"@Transactional",
				"public void save() {",
				"}",
			},
			expected: 1,
		},
		{
			name: "static method",
			code: []string{
				"public static void main(String[] args) {",
				"    System.out.println(\"Hello\");",
				"}",
			},
			expected: 1,
		},
		{
			name: "method with throws",
			code: []string{
				"public void risky() throws IOException {",
				"}",
			},
			expected: 1,
		},
		{
			name: "inner class",
			code: []string{
				"public class Outer {",
				"    private int field;",
				"",
				"    public void method() {",
				"    }",
				"",
				"    class Inner {",
				"        void innerMethod() {",
				"        }",
				"    }",
				"}",
			},
			expected: 2,
		},
		{
			name: "skip class declaration",
			code: []string{
				"public class MyClass {",
				"    private int x;",
				"}",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions := analyzer.detectJavaFunctions(tt.code)
			if len(functions) != tt.expected {
				t.Errorf("Expected %d functions, got %d", tt.expected, len(functions))
				for _, f := range functions {
					t.Logf("Found function: %s at line %d-%d", f.Name, f.StartLine, f.EndLine)
				}
			}
		})
	}
}

func TestDetectorUtils(t *testing.T) {
	t.Run("CalculateIndent", func(t *testing.T) {
		tests := []struct {
			line     string
			expected int
		}{
			{"    hello", 4},
			{"\thello", 4},
			{"  \thello", 6},
			{"hello", 0},
			{"", 0},
		}

		for _, tt := range tests {
			result := CalculateIndent(tt.line)
			if result != tt.expected {
				t.Errorf("CalculateIndent(%q) = %d, expected %d", tt.line, result, tt.expected)
			}
		}
	})

	t.Run("IsKeyword", func(t *testing.T) {
		keywords := []string{"if", "for", "while", "return", "class", "def"}
		nonKeywords := []string{"myFunc", "variable", "MyClass"}

		for _, kw := range keywords {
			if !IsKeyword(kw) {
				t.Errorf("IsKeyword(%q) should be true", kw)
			}
		}

		for _, nk := range nonKeywords {
			if IsKeyword(nk) {
				t.Errorf("IsKeyword(%q) should be false", nk)
			}
		}
	})

	t.Run("IsValidIdentifier", func(t *testing.T) {
		valid := []string{"myVar", "_private", "CamelCase", "snake_case", "var123"}
		invalid := []string{"", "123var", "my-var", "my var", "my.var"}

		for _, id := range valid {
			if !IsValidIdentifier(id) {
				t.Errorf("IsValidIdentifier(%q) should be true", id)
			}
		}

		for _, id := range invalid {
			if IsValidIdentifier(id) {
				t.Errorf("IsValidIdentifier(%q) should be false", id)
			}
		}
	})

	t.Run("StripComments", func(t *testing.T) {
		tests := []struct {
			line     string
			lang     string
			expected string
		}{
			{"int x = 1; // comment", "c", "int x = 1; "},
			{"int x = 1;", "c", "int x = 1;"},
			{"x = 1 # comment", "python", "x = 1 "},
			{"x = 1", "python", "x = 1"},
		}

		for _, tt := range tests {
			result := StripComments(tt.line, tt.lang)
			if result != tt.expected {
				t.Errorf("StripComments(%q, %q) = %q, expected %q", tt.line, tt.lang, result, tt.expected)
			}
		}
	})
}

func TestBracketMatching(t *testing.T) {
	analyzer := NewTextAnalyzer(LanguageC)

	t.Run("simple brackets", func(t *testing.T) {
		lines := []string{
			"int main() {",
			"    return 0;",
			"}",
		}
		end := analyzer.findBracketEndEnhanced(lines, 0)
		if end != 3 {
			t.Errorf("Expected end at line 3, got %d", end)
		}
	})

	t.Run("nested brackets", func(t *testing.T) {
		lines := []string{
			"void func() {",
			"    if (true) {",
			"        // nested",
			"    }",
			"}",
		}
		end := analyzer.findBracketEndEnhanced(lines, 0)
		if end != 5 {
			t.Errorf("Expected end at line 5, got %d", end)
		}
	})

	t.Run("brackets in strings", func(t *testing.T) {
		lines := []string{
			`void func() {`,
			`    char* s = "{nested}";`,
			`}`,
		}
		end := analyzer.findBracketEndEnhanced(lines, 0)
		if end != 3 {
			t.Errorf("Expected end at line 3, got %d", end)
		}
	})

	t.Run("multi-line comment", func(t *testing.T) {
		lines := []string{
			"void func() {",
			"    /* comment with { bracket */",
			"}",
		}
		end := analyzer.findBracketEndEnhanced(lines, 0)
		if end != 3 {
			t.Errorf("Expected end at line 3, got %d", end)
		}
	})
}

func TestTextAnalyzerWithRegistry(t *testing.T) {
	registry := NewPatternRegistry()
	analyzer := NewTextAnalyzerWithRegistry(LanguagePython, registry)

	if analyzer.GetLanguage() != LanguagePython {
		t.Error("Language should be Python")
	}

	if analyzer.GetPatterns() == nil {
		t.Error("Patterns should not be nil")
	}
}
