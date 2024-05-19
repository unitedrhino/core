package coreExport

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	ws "gitee.com/i-Things/share/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"sync"
	"time"
)

const (
	asyncExecMax = 500
)

type publishStu struct {
	*ws.WsPublish
	ctx context.Context
}

type UserSubscribe struct {
	us          *ws.UserSubscribe
	publishChan chan publishStu //key是apisvr的节点id
	mutex       sync.RWMutex
	ServerMsg   *eventBus.FastEvent
}

func NewUserSubscribe(store kv.Store, ServerMsg *eventBus.FastEvent) *UserSubscribe {
	u := UserSubscribe{us: ws.NewUserSubscribe(store), publishChan: make(chan publishStu, asyncExecMax), ServerMsg: ServerMsg}
	utils.Go(context.Background(), func() {
		u.publish()
	})
	return &u
}

func (u *UserSubscribe) Publish(ctx context.Context, code string, data any, params ...map[string]any) error {
	pb := ws.WsPublish{
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
	logx.Infof("UserSubscribe.publish:%v", utils.Fmt(pb))

	return nil
}

func (u *UserSubscribe) publish() {
	execCache := make([]publishStu, 0, asyncExecMax)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		logx.Infof("UserSubscribe.publish publishs:%v", utils.Fmt(execCache))
		err := u.ServerMsg.Publish(context.Background(), fmt.Sprintf(eventBus.CoreApiUserPublish, 1), execCache)
		if err != nil {
			logx.Error(err)
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
