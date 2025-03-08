package websocket

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

func subscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.WithContext(ctx).Errorf("websocket userSubscribe Decode userID:%v connectID:%v body:%v err:%v",
			c.userID, c.connectID, utils.Fmt(body), err)
		c.errorSend(body, err)
		return
	}
	md := utils.Md5Map(info.Params)
	subKey := info.Code + ":" + md

	handle := func(infos []map[string]any) {
		c.userSubscribeMutex.Lock()
		defer c.userSubscribeMutex.Unlock()
		for _, i := range infos {
			md := utils.Md5Map(i)
			key := info.Code + ":" + md
			if c.userSubscribe[subKey] == nil {
				c.userSubscribe[subKey] = make(map[string]struct{})
			}
			if subKey != key {
				c.userSubscribe[subKey][key] = struct{}{}
			}
			func() {
				dp.userSubscribeMutex.Lock()
				defer dp.userSubscribeMutex.Unlock()
				if _, ok := dp.userSubscribe[key]; !ok {
					dp.userSubscribe[key] = map[int64]*connection{}
				}
				dp.userSubscribe[key][c.connectID] = c
			}()
		}
		logx.WithContext(ctx).Infof("websocket userSubscribe userID:%v connectID:%v i:%v subKey:%v keys:%v params:%v subList:%v",
			c.userID, c.connectID, utils.Fmt(info), subKey, c.userSubscribe[subKey], infos, utils.Fmt(c.userSubscribe))
	}
	if checkSubscribe != nil {
		err = checkSubscribe(ctx, &info)
		if err != nil {
			logx.WithContext(ctx).Errorf("websocket userSubscribe checkSubscribe userID:%v connectID:%v body:%v err:%v",
				c.userID, c.connectID, utils.Fmt(body), err)
			c.errorSend(body, err)
			return
		}
		handle([]map[string]any{info.Params})
	}
	if checkSubscribe2 != nil {
		subs, err := checkSubscribe2(ctx, &info)
		if err != nil {
			logx.WithContext(ctx).Errorf("websocket userSubscribe checkSubscribe userID:%v connectID:%v body:%v err:%v",
				c.userID, c.connectID, utils.Fmt(body), err)
			c.errorSend(body, err)
			return
		}
		if len(subs) > 0 {
			handle(subs)
		} else {
			handle([]map[string]any{info.Params})
		}
	}
	var resp = WsResp{WsBody: WsBody{}}
	resp.WsBody.Type = SubRet
	c.sendMessage(resp)
}

func unSubscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(body, err)
		return
	}
	//err = NewUserSubscribe(store).Del(ctx, c.userID, &info)
	md := utils.Md5Map(info.Params)
	key := info.Code + ":" + md
	c.userSubscribeMutex.Lock()
	defer c.userSubscribeMutex.Unlock()
	keys := c.userSubscribe[key]
	delete(c.userSubscribe, key)
	func() {
		dp.userSubscribeMutex.Lock()
		defer dp.userSubscribeMutex.Unlock()
		if _, ok := dp.userSubscribe[key]; ok {
			delete(dp.userSubscribe[key], c.connectID)
		}
		if len(keys) == 0 {
			for k := range keys {
				if _, ok := dp.userSubscribe[k]; ok {
					delete(dp.userSubscribe[k], c.connectID)
				}
			}
		}
	}()
	var resp = WsResp{WsBody: WsBody{}}
	resp.WsBody.Type = UnSubRet
	c.sendMessage(resp)
}
