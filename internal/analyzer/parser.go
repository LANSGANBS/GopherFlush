package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser Go代码解析器
type Parser struct {
	fset *token.FileSet
}

// NewParser 创建新的解析器
func NewParser() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseFile 解析单个Go文件
func (p *Parser) ParseFile(filePath string) (*ast.File, error) {
	// TODO: 实现文件解析
	return parser.ParseFile(p.fset, filePath, nil, parser.ParseComments)
}

// GetFileSet 获取文件集
func (p *Parser) GetFileSet() *token.FileSet {
	return p.fset
}
