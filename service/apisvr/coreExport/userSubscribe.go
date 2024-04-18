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
	asyncExecMax = 100
)

type publishStu struct {
	*ws.WsPublish
	ctx context.Context
}

type UserSubscribe struct {
	us          *ws.UserSubscribe
	publishChan map[int64]chan publishStu //key是apisvr的节点id
	mutex       sync.RWMutex
	ServerMsg   *eventBus.FastEvent
}

func NewUserSubscribe(store kv.Store, ServerMsg *eventBus.FastEvent) *UserSubscribe {
	return &UserSubscribe{us: ws.NewUserSubscribe(store), publishChan: map[int64]chan publishStu{}, ServerMsg: ServerMsg}
}

func (u *UserSubscribe) Publish(ctx context.Context, info *ws.SubscribeInfo, data any) error {
	ret, err := u.us.IndexInfo(ctx, info)
	if err != nil {
		return err
	}
	if len(ret) == 0 {
		return nil
	}
	func() { //初始化channel
		u.mutex.Lock()
		defer u.mutex.Unlock()
		for k := range ret {
			if u.publishChan[k] == nil {
				c := make(chan publishStu, asyncExecMax)
				u.publishChan[k] = c
				kk := k
				utils.Go(ctx, func() {
					u.publish(kk, c)
				})
			}
		}
	}()
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	for k, vs := range ret {
		for _, v := range vs {
			u.publishChan[k] <- publishStu{
				WsPublish: &ws.WsPublish{
					UserID: v,
					Code:   info.Code,
					Data:   data,
				},
				ctx: ctxs.CopyCtx(ctx),
			}
		}
	}
	return nil
}

func (u *UserSubscribe) publish(nodeID int64, infos chan publishStu) {
	execCache := make([]publishStu, 0, asyncExecMax)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		err := u.ServerMsg.Publish(context.Background(), fmt.Sprintf(eventBus.CoreApiUserPublish, nodeID), execCache)
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
		case e := <-infos:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				exec()
			}
		}
	}
}
