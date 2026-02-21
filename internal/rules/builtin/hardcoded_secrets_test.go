package builtin

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestHardcodedSecretsRule_Name(t *testing.T) {
	rule := NewHardcodedSecretsRule()
	if rule.Name() != "hardcoded-secrets" {
		t.Errorf("Expected name 'hardcoded-secrets', got '%s'", rule.Name())
	}
}

func TestHardcodedSecretsRule_Description(t *testing.T) {
	rule := NewHardcodedSecretsRule()
	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestHardcodedSecretsRule_Enabled(t *testing.T) {
	rule := NewHardcodedSecretsRule()
	if !rule.Enabled() {
		t.Error("Rule should be enabled by default")
	}
}

func TestHardcodedSecretsRule_Check_URLs(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "http URL",
			code: `package test
var url = "http://example.com/api"`,
			expected: 1,
		},
		{
			name: "https URL",
			code: `package test
var url = "https://api.example.com/v1"`,
			expected: 1,
		},
		{
			name: "no URL",
			code: `package test
var name = "hello world"`,
			expected: 0,
		},
		{
			name: "URL in function",
			code: `package test
func main() {
	url := "https://example.com"
	_ = url
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			violations := rule.Check(file, fset, "test.go")
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
				for _, v := range violations {
					t.Logf("Violation: %s", v.Message)
				}
			}
		})
	}
}

func TestHardcodedSecretsRule_Check_APIKeys(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "api_key variable",
			code: `package test
var apiKey = "sk-1234567890abcdef"`,
			expected: 1,
		},
		{
			name: "password variable",
			code: `package test
var password = "mysecretpassword123"`,
			expected: 1,
		},
		{
			name: "token variable",
			code: `package test
var token = "ghp_xxxxxxxxxxxxxxxxxxxx"`,
			expected: 1,
		},
		{
			name: "secret variable",
			code: `package test
var secret = "my_super_secret_value"`,
			expected: 1,
		},
		{
			name: "placeholder value",
			code: `package test
var apiKey = "your_api_key_here"`,
			expected: 0,
		},
		{
			name: "example value",
			code: `package test
var token = "example_token_123"`,
			expected: 0,
		},
		{
			name: "short value",
			code: `package test
var key = "short"`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			violations := rule.Check(file, fset, "test.go")
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
				for _, v := range violations {
					t.Logf("Violation: %s", v.Message)
				}
			}
		})
	}
}

func TestHardcodedSecretsRule_Check_DatabaseConnections(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "mysql connection",
			code: `package test
var dsn = "mysql://user:password@localhost:3306/db"`,
			expected: 1,
		},
		{
			name: "postgresql connection",
			code: `package test
var dsn = "postgresql://user:pass@localhost/db"`,
			expected: 1,
		},
		{
			name: "mongodb connection",
			code: `package test
var dsn = "mongodb://localhost:27017"`,
			expected: 1,
		},
		{
			name: "redis connection",
			code: `package test
var dsn = "redis://localhost:6379"`,
			expected: 1,
		},
		{
			name: "jdbc connection",
			code: `package test
var dsn = "jdbc:mysql://localhost/db"`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			violations := rule.Check(file, fset, "test.go")
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
				for _, v := range violations {
					t.Logf("Violation: %s", v.Message)
				}
			}
		})
	}
}

func TestHardcodedSecretsRule_Check_MultipleViolations(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	code := `package test

var (
	apiKey = "sk-1234567890abcdef"
	dbURL  = "mysql://user:password@localhost/db"
)

func main() {
	url := "https://api.example.com"
	token := "ghp_xxxxxxxxxxxxxxxxxxxx"
	_ = url
	_ = token
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	violations := rule.Check(file, fset, "test.go")
	if len(violations) < 3 {
		t.Errorf("Expected at least 3 violations, got %d", len(violations))
		for _, v := range violations {
			t.Logf("Violation: %s at line %d", v.Message, v.Line)
		}
	}
}

func TestHardcodedSecretsRule_IsURL(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		value    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"HTTP://EXAMPLE.COM", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"not a url", false},
	}

	for _, tt := range tests {
		result := rule.isURL(tt.value)
		if result != tt.expected {
			t.Errorf("isURL(%q) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestHardcodedSecretsRule_IsAPIToken(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		varName  string
		value    string
		expected bool
	}{
		{"apiKey", "sk-1234567890abcdef", true},
		{"api_key", "real_api_key_value", true},
		{"token", "ghp_xxxxxxxxxxxxxxxxxxxx", true},
		{"password", "mysecretpassword", true},
		{"secret", "mysecretvalue", true},
		{"apiKey", "your_key_here", false},
		{"apiKey", "example_key", false},
		{"name", "john", false},
		{"short", "abc", false},
	}

	for _, tt := range tests {
		result := rule.isAPIToken(tt.varName, tt.value)
		if result != tt.expected {
			t.Errorf("isAPIToken(%q, %q) = %v, expected %v", tt.varName, tt.value, result, tt.expected)
		}
	}
}

func TestHardcodedSecretsRule_IsDatabaseConnection(t *testing.T) {
	rule := NewHardcodedSecretsRule()

	tests := []struct {
		value    string
		expected bool
	}{
		{"mysql://user:pass@localhost/db", true},
		{"postgresql://localhost/db", true},
		{"postgres://localhost/db", true},
		{"mongodb://localhost:27017", true},
		{"redis://localhost:6379", true},
		{"jdbc:mysql://localhost/db", true},
		{"user:password@localhost", true},
		{"not a connection string", false},
		{"localhost:3306", false},
	}

	for _, tt := range tests {
		result := rule.isDatabaseConnection(tt.value)
		if result != tt.expected {
			t.Errorf("isDatabaseConnection(%q) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}
