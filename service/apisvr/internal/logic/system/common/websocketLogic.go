package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/shared/utils"
	ws "gitee.com/i-Things/core/shared/websocket"
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketLogic {
	return &WebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketLogic) InitWebsocketConn(r *http.Request, conn *websocket.Conn) {

	//创建ws连接
	wsClient := ws.NewConn(l.ctx, l.svcCtx.Ws, r, conn)
	//开启读取进程
	utils.Go(l.ctx, wsClient.StartRead)
	//开启发送进程
	utils.Go(l.ctx, wsClient.StartWrite)
	return
}
