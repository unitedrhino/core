package timedmanagelogic

import (
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/stores"
)

func ToTaskGroupPo(in *timedjob.TaskGroup) *relationDB.TimedTaskGroup {
	if in == nil {
		return nil
	}
	return &relationDB.TimedTaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}
func ToTaskGroupPb(in *relationDB.TimedTaskGroup) *timedjob.TaskGroup {
	if in == nil {
		return nil
	}
	return &timedjob.TaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}

func ToTaskGroupPbs(in []*relationDB.TimedTaskGroup) (ret []*timedjob.TaskGroup) {
	for _, v := range in {
		ret = append(ret, ToTaskGroupPb(v))
	}
	return
}

func ToTaskInfoPbs(in []*relationDB.TimedTaskInfo) (ret []*timedjob.TaskInfo) {
	for _, v := range in {
		ret = append(ret, ToTaskInfoPb(v))
	}
	return
}

func ToTaskInfoPb(in *relationDB.TimedTaskInfo) *timedjob.TaskInfo {
	if in == nil {
		return nil
	}
	return &timedjob.TaskInfo{
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

func ToTaskInfoPo(in *timedjob.TaskInfo) *relationDB.TimedTaskInfo {
	if in == nil {
		return nil
	}
	return &relationDB.TimedTaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Topics:    in.Topics,
		Priority:  in.Priority,
	}
}

func ToPageInfo(info *timedjob.PageInfo, defaultOrders ...stores.OrderBy) *stores.PageInfo {
	if info == nil {
		return nil
	}

	var orders = defaultOrders
	if infoOrders := info.GetOrders(); len(infoOrders) > 0 {
		orders = make([]stores.OrderBy, 0, len(infoOrders))
		for _, infoOd := range infoOrders {
			if infoOd.GetFiled() != "" {
				orders = append(orders, stores.OrderBy{Field: infoOd.GetFiled(), Sort: infoOd.GetSort()})
			}
		}
	}

	return &stores.PageInfo{
		Page:   info.GetPage(),
		Size:   info.GetSize(),
		Orders: orders,
	}
}

func ToPageInfoWithDefault(info *timedjob.PageInfo, defau *stores.PageInfo) *stores.PageInfo {
	if page := ToPageInfo(info); page == nil {
		return defau
	} else {
		if page.Page == 0 {
			page.Page = defau.Page
		}
		if page.Size == 0 {
			page.Size = defau.Size
		}
		if len(page.Orders) == 0 {
			page.Orders = defau.Orders
		}
		return page
	}
}
func ToTaskLog(po *relationDB.TimedTaskLog) *timedjob.TaskLog {
	var sql *timedjob.TaskLogSql
	if po.TimedTaskLogSql != nil {
		sql = &timedjob.TaskLogSql{
			SelectNum: po.SelectNum,
			ExecNum:   po.ExecNum,
		}
	}
	var script *timedjob.TaskLogScript
	if po.TimedTaskLogScript != nil {
		var execLog = []*timedjob.TaskExecLog{}
		for _, v := range po.ExecLog {
			execLog = append(execLog, &timedjob.TaskExecLog{
				Level:       v.Level,
				Content:     v.Content,
				CreatedTime: v.CreatedTime,
			})
		}
		script = &timedjob.TaskLogScript{ExecLog: execLog}
	}
	return &timedjob.TaskLog{
		Id:          po.ID,
		GroupCode:   po.GroupCode,
		TaskCode:    po.TaskCode,
		Params:      po.Params,
		ResultCode:  po.ResultCode,
		ResultMsg:   po.ResultMsg,
		CreatedTime: po.CreatedTime.Unix(),
		Sql:         sql,
		Script:      script,
	}
}
