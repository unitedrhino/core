package task

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
)

func ToSendDelayReqPb(in *types.TimedTaskSendReq) *timedmanage.TaskSendReq {
	ret := timedmanage.TaskSendReq{GroupCode: in.GroupCode, Code: in.Code}
	if in.Option != nil {
		ret.Option = &timedmanage.TaskSendOption{
			Priority:  in.Option.Priority,
			ProcessIn: in.Option.ProcessIn,
			ProcessAt: in.Option.ProcessAt,
			Timeout:   in.Option.Timeout,
			Deadline:  in.Option.Deadline,
			TaskID:    in.Option.TaskID,
		}
	}
	if in.ParamSql != nil {
		ret.ParamSql = &timedmanage.TaskParamSql{Sql: in.ParamSql.Sql}
	}
	if in.ParamScript != nil {
		ret.ParamScript = &timedmanage.TaskParamScript{ExecContent: in.ParamScript.ExecContent, Param: in.ParamScript.Param}
	}
	if in.ParamQueue != nil {
		ret.ParamQueue = &timedmanage.TaskParamQueue{Topic: in.ParamQueue.Topic, Payload: in.ParamQueue.Payload}
	}
	return &ret
}

func ToGroupPb(in *types.TimedTaskGroup) *timedmanage.TaskGroup {
	if in == nil {
		return nil
	}
	return &timedmanage.TaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}

func ToGroupTypes(in *timedmanage.TaskGroup) *types.TimedTaskGroup {
	if in == nil {
		return nil
	}
	return &types.TimedTaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}
func ToTaskGroupsTypes(in []*timedmanage.TaskGroup) (ret []*types.TimedTaskGroup) {
	for _, v := range in {
		ret = append(ret, ToGroupTypes(v))
	}
	return
}

func ToTaskInfoPb(in *types.TimedTaskInfo) *timedmanage.TaskInfo {
	if in == nil {
		return nil
	}
	return &timedmanage.TaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		Topics:    in.Topics,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToTaskInfoTypes(in *timedmanage.TaskInfo) *types.TimedTaskInfo {
	if in == nil {
		return nil
	}
	return &types.TimedTaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Topics:    in.Topics,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToTaskInfosTypes(in []*timedmanage.TaskInfo) (ret []*types.TimedTaskInfo) {
	for _, v := range in {
		ret = append(ret, ToTaskInfoTypes(v))
	}
	return
}

func ToTaskLog(pb *timedjob.TaskLog) *types.TimedTaskLog {
	var sql *types.TimedTaskLogSql
	var script *types.TimedTaskLogScript
	if pb.Sql != nil {
		sql = &types.TimedTaskLogSql{
			SelectNum: pb.Sql.SelectNum,
			ExecNum:   pb.Sql.ExecNum,
		}
	}
	if pb.Script != nil {
		var execLog = []*types.TaskLogScript{}
		for _, v := range pb.Script.ExecLog {
			execLog = append(execLog, &types.TaskLogScript{
				Level:       v.Level,
				Content:     v.Content,
				CreatedTime: v.CreatedTime,
			})
		}
		script = &types.TimedTaskLogScript{
			ExecLog: execLog,
		}
	}
	return &types.TimedTaskLog{
		ID:          pb.Id,
		GroupCode:   pb.GroupCode,
		TaskCode:    pb.TaskCode,
		Params:      pb.Params,
		ResultCode:  pb.ResultCode,
		ResultMsg:   pb.ResultMsg,
		CreatedTime: pb.CreatedTime,
		Sql:         sql,
		Script:      script,
	}
}
