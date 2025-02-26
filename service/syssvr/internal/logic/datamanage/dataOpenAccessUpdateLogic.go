package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataOpenAccessUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataOpenAccessUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataOpenAccessUpdateLogic {
	return &DataOpenAccessUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataOpenAccessUpdateLogic) DataOpenAccessUpdate(in *sys.OpenAccess) (*sys.Empty, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	old, err := relationDB.NewDataOpenAccessRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if string(old.TenantCode) != uc.TenantCode && !(uc.TenantCode == def.TenantCodeDefault && uc.IsAdmin) {
		return nil, errors.Permissions
	}
	if !uc.IsAdmin && old.UserID != uc.UserID {
		return nil, errors.Permissions
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.IpRange != nil {
		old.IpRange = in.IpRange
	}
	err = relationDB.NewDataOpenAccessRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
