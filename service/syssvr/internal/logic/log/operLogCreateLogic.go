package loglogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"sync"

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

var asyncOperInsert *stores.AsyncInsert[relationDB.SysOperLog]
var operOnce sync.Once

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	operOnce.Do(func() {
		asyncOperInsert = stores.NewAsyncInsert[relationDB.SysOperLog]()
	})
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
	//OperName，BusinessType 用Route查接口管理表获得
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	asyncOperInsert.AsyncInsert(&relationDB.SysOperLog{
		TenantCode:   stores.TenantCode(uc.TenantCode),
		AppCode:      uc.AppCode,
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
	return &sys.Empty{}, nil
}
