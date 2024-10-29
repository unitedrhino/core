package dept

type DeptMode = int64

var (
	DeptModeAdd    DeptMode = 1 //1 只新增,不修改(默认)
	DeptModeUpdate DeptMode = 2 //2 新增并修改
	DeptModeDelete DeptMode = 3 //3 新增修改及删除不存在的部门
)

type UserMode = int64

var (
	UserModeNone   UserMode = 0 //0 不同步
	UserModeAdd    UserMode = 1 //1 只新增,不修改(默认)
	UserModeUpdate UserMode = 2 //2 新增并修改
	UserModeDelete UserMode = 3 //3 新增修改及删除不存在的用户
)
