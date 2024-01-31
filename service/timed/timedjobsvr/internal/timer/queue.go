package timer

import (
	"context"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func (t Timed) Queue(ctx context.Context, task *domain.TaskInfo) error {
	err := t.SvcCtx.PubJob.Publish(ctx, task.GroupSubType, task.Queue.Topic, []byte(task.Queue.Payload))
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
