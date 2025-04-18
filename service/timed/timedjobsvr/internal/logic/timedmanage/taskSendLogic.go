package timedmanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/processTask"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"github.com/hibiken/asynq"
	"time"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskSendLogic {
	return &TaskSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送延时请求,如果任务不存在,则会自动创建,但是自动创建的需要填写param
func (l *TaskSendLogic) TaskSend(in *timedjob.TaskSendReq) (*timedjob.TaskWithTaskID, error) {
	tg := relationDB.NewTaskGroupRepo(l.ctx)
	group, err := tg.FindOneByFilter(l.ctx, relationDB.TaskGroupFilter{Codes: []string{in.GroupCode}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("任务组未找到,请填写正确的任务组")
		}
		return nil, err
	}
	ti := relationDB.NewTaskInfoRepo(l.ctx)
	task, err := ti.FindOneByFilter(l.ctx, relationDB.TaskFilter{Codes: []string{in.Code}, WithGroup: true})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}

	if task == nil { //如果数据库中没有这个任务,需要检查并动态创建
		var param string
		switch group.Type {
		case domain.TaskGroupTypeQueue:
			if in.ParamQueue == nil {
				return nil, errors.Parameter.AddMsg("任务组为消息发送类型,请填写消息发送参数")
			}
			p, _ := json.Marshal(domain.ParamQueue{Topic: in.ParamQueue.Topic, Payload: in.ParamQueue.Payload})
			param = string(p)
		case domain.TaskGroupTypeSql:
			if in.ParamSql == nil {
				return nil, errors.Parameter.AddMsg("任务组为sql执行类型,请填写sql执行参数")
			}
			p, _ := json.Marshal(domain.ParamSql{Sql: in.ParamSql.Sql})
			param = string(p)
		case domain.TaskGroupTypeScript:
			if in.ParamScript == nil {
				return nil, errors.Parameter.AddMsg("任务组为script执行类型,请填写script执行参数")
			}
			p, _ := json.Marshal(domain.ParamScript{ExecContent: in.ParamScript.ExecContent, Param: in.ParamScript.Param})
			param = string(p)
		}
		property := int64(3)
		if in.GetOption() != nil && in.GetOption().Priority != 0 {
			property = in.Option.Priority
		}
		task = &relationDB.TimedTaskInfo{
			GroupCode: in.GroupCode,
			Type:      domain.TaskTypeDelay,
			Code:      in.Code,
			Params:    param,
			Status:    def.StatusRunning,
			Priority:  property,
		}
		err := ti.Insert(l.ctx, task)
		if err != nil {
			return nil, err
		}
	} else { //如果传了参数需要更新参数内容
		oldDo := relationDB.ToTaskInfoDo(task)
		switch group.Type {
		case domain.TaskGroupTypeQueue:
			if in.ParamQueue != nil {
				p, _ := json.Marshal(domain.ParamQueue{Topic: in.ParamQueue.Topic, Payload: in.ParamQueue.Payload})
				task.Params = string(p)
			}
		case domain.TaskGroupTypeSql:
			if in.ParamSql != nil {
				ps := domain.ParamSql{Sql: oldDo.Sql.Param.Sql}
				if in.ParamSql.Sql != "" {
					ps.Sql = in.ParamSql.Sql
				}
				p, _ := json.Marshal(ps)
				task.Params = string(p)
			}
		case domain.TaskGroupTypeScript:
			if in.ParamScript != nil {
				ps := domain.ParamScript{ExecContent: oldDo.Script.Param.ExecContent, Param: oldDo.Script.Param.Param}
				if in.ParamScript.ExecContent != "" {
					ps.ExecContent = in.ParamScript.ExecContent
				}
				if in.ParamScript.Param != nil {
					ps.Param = in.ParamScript.Param
				}
				p, _ := json.Marshal(ps)
				task.Params = string(p)
			}
		}
	}
	taskInfo := domain.TaskInfo{
		ID:     task.ID,
		Params: task.Params,
	}
	if in.Option == nil { //立即执行
		err := logic.FillTaskInfoDo(&taskInfo, task)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		return &timedjob.TaskWithTaskID{}, processTask.NewProcessTask(l.ctx, l.svcCtx, func(ctx context.Context, req *timedjob.TaskSendReq) error {
			_, err := NewTaskSendLogic(ctx, l.svcCtx).TaskSend(req)
			return err
		}).Process(l.ctx, taskInfo)
	}
	payload, _ := json.Marshal(taskInfo)
	aTask := asynq.NewTask(getTaskCode(task), payload, asynq.Queue(domain.ToPriority(task.Priority)))
	var opts []asynq.Option
	if in.Option != nil {
		var opt = asynq.ProcessAt(time.Unix(in.Option.ProcessAt, 0))
		if in.Option.ProcessIn != 0 {
			opt = asynq.ProcessIn(time.Duration(in.Option.ProcessIn) * time.Second)
		}
		opts = append(opts, opt)
		if in.Option.Timeout != 0 {
			opts = append(opts, asynq.Timeout(time.Duration(in.Option.Timeout)*time.Second))
		}
		if in.Option.Deadline != 0 {
			opts = append(opts, asynq.Deadline(time.Unix(in.Option.Deadline, 0)))
		}
		if in.Option.TaskID != "" {
			opts = append(opts, asynq.TaskID(in.Option.TaskID))
		}
	}
	t, err := l.svcCtx.AsynqClient.Enqueue(aTask, opts...)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	return &timedjob.TaskWithTaskID{TaskID: t.ID}, nil
}

func getTaskCode(j *relationDB.TimedTaskInfo) string {
	return fmt.Sprintf("delay:%v:%s", j.GroupCode, j.Code)
}
