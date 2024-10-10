package common

import (
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/common"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

func WebsocketConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			// 读取存储空间大小
			ReadBufferSize: 1024,
			// 写入存储空间大小
			WriteBufferSize: 1024,
			// 允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		//ws连接失败
		if err != nil {
			result.Http(w, r, nil, err)
			logx.WithContext(r.Context()).Error("[ws]连接失败", "RemoteAddr:", r.RemoteAddr, "err", err)
			return
		}
		l := common.NewWebsocketConnectLogic(r.Context(), svcCtx)
		l.WebsocketConnect(r, conn)
	}
}
