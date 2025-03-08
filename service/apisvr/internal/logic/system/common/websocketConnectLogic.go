package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	ws "gitee.com/unitedrhino/core/share/websocket"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type WebsocketConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketConnectLogic {
	return &WebsocketConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketConnectLogic) WebsocketConnect(r *http.Request, conn *websocket.Conn) error {
	userID := ctxs.GetUserCtx(l.ctx).UserID
	//创建ws连接
	wsClient := ws.NewConn(l.ctx, userID, l.svcCtx.Ws, r, conn)
	//开启读取进程
	utils.Go(l.ctx, wsClient.StartRead)
	//开启发送进程
	utils.Go(l.ctx, wsClient.StartWrite)
	return nil
}
