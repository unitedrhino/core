package relationDB

import (
	"encoding/json"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
)

func ToTaskInfoDo(po *TimedTaskInfo) *domain.TaskInfo {
	var do domain.TaskInfo
	do.Code = po.Code
	do.GroupType = po.Group.Type
	do.GroupSubType = po.Group.SubType
	do.GroupCode = po.GroupCode
	do.Env = po.Group.Env
	switch po.Group.Type {
	case domain.TaskGroupTypeQueue:
		var param domain.ParamQueue
		json.Unmarshal([]byte(po.Params), &param)
		do.Queue = &param
	case domain.TaskGroupTypeScript:
		var sql domain.ParamScript
		json.Unmarshal([]byte(po.Params), &sql)
		do.Script = &domain.Script{Param: sql}
		json.Unmarshal([]byte(po.Group.Config), &do.Script.Config)
	case domain.TaskGroupTypeSql:
		var sql domain.ParamSql
		json.Unmarshal([]byte(po.Params), &sql)
		do.Sql = &domain.Sql{Param: sql}
		json.Unmarshal([]byte(po.Group.Config), &do.Sql.Config)
	}
	return &do
}
