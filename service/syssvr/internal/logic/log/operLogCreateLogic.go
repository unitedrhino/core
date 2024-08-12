package loglogic

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OlDB *relationDB.OperLogRepo
	UiDB *relationDB.UserInfoRepo
	AiDB *relationDB.ApiInfoRepo
}

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	return &OperLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OlDB:   relationDB.NewOperLogRepo(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *OperLogCreateLogic) OperLogCreate(in *sys.OperLogCreateReq) (*sys.Empty, error) {
	//OperUserName 用uid查用户表获得
	uc := ctxs.GetUserCtx(l.ctx)
	//OperName，BusinessType 用Route查接口管理表获得

	err := l.OlDB.Insert(l.ctx, &relationDB.SysOperLog{
		AppCode:      in.AppCode,
		OperUserID:   uc.UserID,
		OperUserName: uc.Account,
		OperName:     in.OperName,
		BusinessType: in.BusinessType,
		Uri:          in.Uri,
		OperIpAddr:   in.OperIpAddr,
		OperLocation: in.OperLocation,
		Req:          sql.NullString{String: in.Req, Valid: true},
		Resp:         sql.NullString{String: in.Resp, Valid: true},
		Code:         in.Code,
		Msg:          in.Msg,
	})
	if err != nil {
		return nil, err
	}

	return &sys.Empty{}, nil
}
