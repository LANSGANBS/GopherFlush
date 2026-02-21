package builtin

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestInaccurateConstantRule(t *testing.T) {
	t.Run("detect inaccurate PI", func(t *testing.T) {
		code := `package test

const PI = 3.14
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d", len(violations))
		}
	})

	t.Run("accept accurate PI", func(t *testing.T) {
		code := `package test

const PI = 3.14159265358979323846
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("detect inaccurate E", func(t *testing.T) {
		code := `package test

const E = 2.7
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d", len(violations))
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		code := `package test

const pi = 3.14
const Pi = 3.14
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 2 {
			t.Errorf("Expected 2 violations, got %d", len(violations))
		}
	})

	t.Run("unknown constant", func(t *testing.T) {
		code := `package test

const MY_VALUE = 123.456
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Expected 0 violations for unknown constant, got %d", len(violations))
		}
	})
}

func TestInaccurateConstantRuleWithCustomConstants(t *testing.T) {
	t.Run("custom constant detection", func(t *testing.T) {
		code := `package test

const GRAVITY = 9.8
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRuleWithConstants([]ConstantDefinition{
			{Name: "GRAVITY", Value: 9.80665},
		})
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d", len(violations))
		}
	})

	t.Run("add constant dynamically", func(t *testing.T) {
		code := `package test

const SPEED_OF_LIGHT = 299792458
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		rule.AddConstant("SPEED_OF_LIGHT", 299792458.0)

		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Expected 0 violations, got %d", len(violations))
		}
	})
}

func TestInaccurateConstantRuleTolerance(t *testing.T) {
	t.Run("custom tolerance", func(t *testing.T) {
		code := `package test

const PI = 3.14
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		rule.SetTolerance(0.001)

		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation with tight tolerance, got %d", len(violations))
		}
	})

	t.Run("loose tolerance", func(t *testing.T) {
		code := `package test

const PI = 3.14
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		rule.SetTolerance(1.0)

		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Expected 0 violations with loose tolerance, got %d", len(violations))
		}
	})
}

func TestInaccurateConstantRuleExtractValue(t *testing.T) {
	t.Run("binary expression", func(t *testing.T) {
		code := `package test

const PI_DIV_2 = 3.14159 / 2
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Binary expression should not trigger for non-matching name, got %d violations", len(violations))
		}
	})

	t.Run("paren expression", func(t *testing.T) {
		code := `package test

const PI = (3.14159265358979323846)
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		rule := NewInaccurateConstantRule()
		violations := rule.Check(file, fset, "test.go")

		if len(violations) != 0 {
			t.Errorf("Expected 0 violations for accurate PI in parens, got %d", len(violations))
		}
	})
}

func TestInaccurateConstantRuleMetadata(t *testing.T) {
	rule := NewInaccurateConstantRule()

	if rule.Name() != "inaccurate-constant" {
		t.Errorf("Expected name 'inaccurate-constant', got %q", rule.Name())
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if !rule.Enabled() {
		t.Error("Rule should be enabled by default")
	}

	constants := rule.GetKnownConstants()
	if len(constants) == 0 {
		t.Error("Should have known constants")
	}

	if _, ok := constants["PI"]; !ok {
		t.Error("PI should be in known constants")
	}
}
