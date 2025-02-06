package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataOpenAccessCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataOpenAccessCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataOpenAccessCreateLogic {
	return &DataOpenAccessCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataOpenAccessCreateLogic) DataOpenAccessCreate(in *sys.OpenAccess) (*sys.WithID, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if in.TenantCode != "" && in.TenantCode != uc.TenantCode {
		if !(uc.TenantCode == def.TenantCodeDefault) {
			return nil, errors.Permissions
		}
		l.ctx = ctxs.BindTenantCode(l.ctx, in.TenantCode, 0)
	}
	if in.AccessSecret == "" {
		in.AccessSecret = utils.Random(32, 4)
	}
	if !uc.IsAdmin || in.UserID == 0 {
		in.UserID = uc.UserID
	}
	if !uc.IsAdmin && in.UserID != uc.UserID {
		return nil, errors.Permissions
	}
	po := &relationDB.SysDataOpenAccess{
		TenantCode:   dataType.TenantCode(in.TenantCode),
		UserID:       in.UserID,
		Code:         in.Code,
		AccessSecret: in.AccessSecret,
		Desc:         in.Desc,
		IpRange:      in.IpRange,
	}
	err := relationDB.NewDataOpenAccessRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, err
}
