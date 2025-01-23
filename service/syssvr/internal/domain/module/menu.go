package module

type MenuImportMode = int64

var (
	MenuImportModeAdd    MenuImportMode = 1 //只新增
	MenuImportModeUpdate MenuImportMode = 2 //新增并修改
	MenuImportModeAll    MenuImportMode = 3 //新增修改并删除没有的

)
