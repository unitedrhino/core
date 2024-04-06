package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoDeleteLogic {
	return &UserInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoDeleteLogic) UserInfoDelete(in *sys.UserInfoDeleteReq) (*sys.Empty, error) {
	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{Code: ctxs.GetUserCtx(l.ctx).TenantCode})
	if err != nil {
		return nil, err
	}
	if ti.AdminUserID == in.UserID {
		return nil, errors.Permissions.AddMsg("超级管理员不允许删除")
	}
	err = l.UiDB.Delete(l.ctx, cast.ToInt64(in.UserID))
	if err != nil {
		l.Errorf("%s.Delete uid=%d err=%+v", utils.FuncName(), in.UserID, err)
		return nil, err
	}

	l.Infof("%s.delete uid=%v", utils.FuncName(), in.UserID)

	return &sys.Empty{}, nil
}
