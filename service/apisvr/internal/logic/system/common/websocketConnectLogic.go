package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	ws "gitee.com/i-Things/share/websocket"
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