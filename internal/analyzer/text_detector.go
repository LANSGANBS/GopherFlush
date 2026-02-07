package analyzer

// detectFunctions 检测函数定义
func (ta *TextAnalyzer) detectFunctions(lines []string) []*FunctionInfo {
	functions := []*FunctionInfo{}

	switch ta.language {
	case LanguagePython:
		functions = ta.detectPythonFunctions(lines)
	case LanguageJava:
		functions = ta.detectJavaFunctions(lines)
	case LanguageC, LanguageCPP:
		functions = ta.detectCFunctions(lines)
	}

	return functions
}

// detectGlobalVars 检测全局变量
func (ta *TextAnalyzer) detectGlobalVars(lines []string) []*GlobalVarInfo {
	globals := []*GlobalVarInfo{}

	switch ta.language {
	case LanguagePython:
		globals = ta.detectPythonGlobals(lines)
	case LanguageJava:
		globals = ta.detectJavaGlobals(lines)
	case LanguageC, LanguageCPP:
		globals = ta.detectCGlobals(lines)
	}

	return globals
}
