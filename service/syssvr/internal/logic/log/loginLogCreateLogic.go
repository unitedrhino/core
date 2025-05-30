package loglogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"sync"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	LlDB *relationDB.LoginLogRepo
}

var asyncLoginInsert *stores.AsyncInsert[relationDB.SysLoginLog]
var loginOnce sync.Once

func NewLoginLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogCreateLogic {
	loginOnce.Do(func() {
		asyncLoginInsert = stores.NewAsyncInsert[relationDB.SysLoginLog](stores.GetTenantConn(ctx), "")
	})
	return &LoginLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		LlDB:   relationDB.NewLoginLogRepo(ctx),
	}
}

func (l *LoginLogCreateLogic) LoginLogCreate(in *sys.LoginLogCreateReq) (*sys.Empty, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	asyncLoginInsert.AsyncInsert(&relationDB.SysLoginLog{
		TenantCode:    dataType.TenantCode(uc.TenantCode),
		AppCode:       in.AppCode,
		UserID:        in.UserID,
		UserName:      in.UserName,
		IpAddr:        in.IpAddr,
		LoginLocation: in.LoginLocation,
		Browser:       in.Browser,
		Os:            in.Os,
		Code:          in.Code,
		Msg:           in.Msg,
	})
	return &sys.Empty{}, nil
}
