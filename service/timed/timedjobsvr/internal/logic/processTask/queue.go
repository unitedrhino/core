package processTask

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func (t ProcessTask) Queue(ctx context.Context, task *domain.TaskInfo) error {
	if task.Queue == nil || !isValidPublishSubject(task.Queue.Topic) {
		err := errors.Parameter.AddMsgf("invalid nats subject: %q", task.Queue.Topic)
		logx.WithContext(ctx).Errorf("Queue.Publish invalid subject task:%+v", task)
		_ = relationDB.NewTaskInfoRepo(ctx).UpdateByFilter(ctx, &relationDB.TimedTaskInfo{Status: def.StatusWaitStop},
			relationDB.TaskFilter{IDs: []int64{task.ID}})
		return err
	}
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

func isValidPublishSubject(topic string) bool {
	topic = strings.TrimSpace(topic)
	if topic == "" || strings.ContainsAny(topic, " \t\r\n") {
		return false
	}
	parts := strings.Split(topic, ".")
	for _, part := range parts {
		if part == "" {
			return false
		}
		if part == "*" || part == ">" {
			return false
		}
	}
	return true
}
