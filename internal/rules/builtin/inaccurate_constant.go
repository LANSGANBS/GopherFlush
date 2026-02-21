package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"math"
	"strconv"
	"strings"
)

type InaccurateConstantRule struct {
	enabled       bool
	knownConstants map[string]float64
	tolerance     float64
}

type ConstantDefinition struct {
	Name  string
	Value float64
}

func NewInaccurateConstantRule() *InaccurateConstantRule {
	return &InaccurateConstantRule{
		enabled: true,
		tolerance: 1e-12,
		knownConstants: map[string]float64{
			"PI":    3.14159265358979323846,
			"E":     2.71828182845904523536,
			"PHI":   1.61803398874989484820,
			"SQRT2": 1.41421356237309504880,
			"LN2":   0.69314718055994530941,
			"LN10":  2.30258509299404568401,
		},
	}
}

func NewInaccurateConstantRuleWithConstants(constants []ConstantDefinition) *InaccurateConstantRule {
	rule := NewInaccurateConstantRule()
	for _, c := range constants {
		rule.knownConstants[strings.ToUpper(c.Name)] = c.Value
	}
	return rule
}

func (r *InaccurateConstantRule) Name() string {
	return "inaccurate-constant"
}

func (r *InaccurateConstantRule) Description() string {
	return "检测AI生成的不准确常量定义"
}

func (r *InaccurateConstantRule) Enabled() bool {
	return r.enabled
}

func (r *InaccurateConstantRule) SetTolerance(tolerance float64) {
	r.tolerance = tolerance
}

func (r *InaccurateConstantRule) AddConstant(name string, value float64) {
	if r.knownConstants == nil {
		r.knownConstants = make(map[string]float64)
	}
	r.knownConstants[strings.ToUpper(name)] = value
}

func (r *InaccurateConstantRule) GetKnownConstants() map[string]float64 {
	result := make(map[string]float64, len(r.knownConstants))
	for k, v := range r.knownConstants {
		result[k] = v
	}
	return result
}

func (r *InaccurateConstantRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					continue
				}

				violation := r.checkConstant(name.Name, valueSpec.Values[i], fset, filePath)
				if violation != nil {
					violations = append(violations, violation)
				}
			}
		}
	}

	return violations
}

func (r *InaccurateConstantRule) checkConstant(name string, value ast.Expr, fset *token.FileSet, filePath string) *types.Violation {
	expectedValue, exists := r.knownConstants[strings.ToUpper(name)]
	if !exists {
		return nil
	}

	actualValue := r.extractValue(value)
	if actualValue == 0 {
		return nil
	}

	diff := math.Abs(expectedValue - actualValue)

	if diff > r.tolerance {
		position := fset.Position(value.Pos())
		return &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityHigh,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "常量 '" + name + "' 的值不准确，当前值: " + strconv.FormatFloat(actualValue, 'f', -1, 64) + "，正确值应为: " + strconv.FormatFloat(expectedValue, 'f', -1, 64),
			Suggestion: "使用更准确的常量值，或使用标准库中的常量（如 math.Pi）",
		}
	}

	return nil
}

func (r *InaccurateConstantRule) extractValue(expr ast.Expr) float64 {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.FLOAT || v.Kind == token.INT {
			val, err := strconv.ParseFloat(v.Value, 64)
			if err == nil {
				return val
			}
		}
	case *ast.BinaryExpr:
		left := r.extractValue(v.X)
		right := r.extractValue(v.Y)
		if left == 0 || right == 0 {
			return 0
		}
		switch v.Op {
		case token.ADD:
			return left + right
		case token.SUB:
			return left - right
		case token.MUL:
			return left * right
		case token.QUO:
			if right != 0 {
				return left / right
			}
		}
	case *ast.ParenExpr:
		return r.extractValue(v.X)
	}
	return 0
}
