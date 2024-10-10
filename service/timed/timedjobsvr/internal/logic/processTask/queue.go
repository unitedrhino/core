package processTask

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func (t ProcessTask) Queue(ctx context.Context, task *domain.TaskInfo) error {
	err := t.svcCtx.PubJob.Publish(ctx, task.GroupSubType, task.Queue.Topic, []byte(task.Queue.Payload))
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		GroupCode:  task.GroupCode,
		TaskCode:   task.Code,
		Params:     task.Params,
		ResultCode: e.GetCode(),
		ResultMsg:  e.GetMsg(),
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("Queue.Publish.Insert err:%v", er)
	}
	return err
}
