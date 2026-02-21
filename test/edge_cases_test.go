package test

var GlobalVar1 = 100
var globalVar2 = "test"

const (
	MultiConst1 = 1
	MultiConst2 = 2
)

type GlobalType struct {
	Field1 string
	Field2 int
}

func init() {
	GlobalVar1 = 200
}

var (
	PackageVar1 int
	PackageVar2 string
)

var ComplexGlobal = map[string]interface{}{
	"key": "value",
}

var SliceGlobal = []int{1, 2, 3}

var FuncGlobal = func() int {
	return 42
}

var (
	_ = "blank identifier"
)

var (
	ExportedGlobal   = "exported"
	unexportedGlobal = "unexported"
)
