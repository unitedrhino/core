package websocket

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"sync"
	"time"
)

const (
	asyncExecMax = 500
)

type publishStu struct {
	*WsPublish
	ctx context.Context
}

type UserSubscribe struct {
	publishChan chan publishStu //key是apisvr的节点id
	mutex       sync.RWMutex
	ServerMsg   *eventBus.FastEvent
}

func NewUserSubscribe(store kv.Store, ServerMsg *eventBus.FastEvent) *UserSubscribe {
	u := UserSubscribe{publishChan: make(chan publishStu, asyncExecMax), ServerMsg: ServerMsg}
	utils.Go(context.Background(), func() {
		u.publish()
	})
	return &u
}

func (u *UserSubscribe) Publish(ctx context.Context, code string, data any, params ...map[string]any) error {
	pb := WsPublish{
		Code: code,
		Data: data,
	}
	for _, param := range params {
		pb.Params = append(pb.Params, utils.Md5Map(param))
	}
	u.publishChan <- publishStu{
		WsPublish: &pb,
		ctx:       ctxs.CopyCtx(ctx),
	}
	logx.WithContext(ctx).Debugf("websocket UserSubscribe.publish pb:%v params:%v", utils.Fmt(pb), utils.Fmt(params))
	return nil
}

func (u *UserSubscribe) publish() {
	execCache := make([]publishStu, 0, asyncExecMax)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		logx.WithContext(execCache[0].ctx).Debugf("websocket UserSubscribe.publish publishs:%v", utils.Fmt(execCache))
		err := u.ServerMsg.Publish(execCache[0].ctx, fmt.Sprintf(eventBus.CoreApiUserPublish, 1), execCache)
		if err != nil {
			logx.WithContext(execCache[0].ctx).Error(err)
		}
		execCache = execCache[0:0] //清空切片
	}
	tick := time.Tick(time.Second)
	for {
		select {
		case _ = <-tick:
			exec()
		case e := <-u.publishChan:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				exec()
			}
		}
	}
}
