package logic

import (
	"encoding/json"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
)

func FillTaskInfoDo(do *domain.TaskInfo, po *relationDB.TimedTaskInfo) error {
	if po.Group == nil {
		return errors.Parameter.AddMsg("任务没有找到任务组")
	}
	do.Code = po.Code
	do.GroupType = po.Group.Type
	do.GroupSubType = po.Group.SubType
	do.GroupCode = po.GroupCode
	if do.Params == "" { //如果没有传,则用数据库里的
		do.Params = po.Params
	}
	do.Env = po.Group.Env
	switch po.Group.Type {
	case domain.TaskGroupTypeQueue:
		var param domain.ParamQueue
		err := json.Unmarshal([]byte(do.Params), &param)
		if err != nil {
			return err
		}
		do.Queue = &param
	case domain.TaskGroupTypeScript:
		var sql domain.ParamScript
		err := json.Unmarshal([]byte(do.Params), &sql)
		if err != nil {
			return err
		}
		do.Script = &domain.Script{Param: sql}
		err = json.Unmarshal([]byte(po.Group.Config), &do.Script.Config)
		if err != nil {
			return err
		}
	case domain.TaskGroupTypeSql:
		var sql domain.ParamSql
		err := json.Unmarshal([]byte(do.Params), &sql)
		if err != nil {
			return err
		}
		do.Sql = &domain.Sql{Param: sql}
		err = json.Unmarshal([]byte(po.Group.Config), &do.Sql.Config)
		if err != nil {
			return err
		}
	}
	return nil
}
