package loglogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	LlDB *relationDB.LoginLogRepo
}

func NewLoginLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogCreateLogic {
	return &LoginLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		LlDB:   relationDB.NewLoginLogRepo(ctx),
	}
}

func (l *LoginLogCreateLogic) LoginLogCreate(in *sys.LoginLogCreateReq) (*sys.Empty, error) {
	err := l.LlDB.Insert(l.ctx, &relationDB.SysLoginLog{
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
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
