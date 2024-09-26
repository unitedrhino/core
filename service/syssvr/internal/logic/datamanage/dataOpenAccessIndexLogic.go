package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataOpenAccessIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataOpenAccessIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataOpenAccessIndexLogic {
	return &DataOpenAccessIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataOpenAccessIndexLogic) DataOpenAccessIndex(in *sys.OpenAccessIndexReq) (*sys.OpenAccessIndexResp, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if in.TenantCode != "" && ctxs.IsRoot(l.ctx) != nil {
		return nil, errors.Permissions
	}
	if !uc.IsAdmin {
		in.UserID = uc.UserID
	}
	if in.TenantCode != "" {
		l.ctx = ctxs.BindTenantCode(l.ctx, in.TenantCode, 0)
	}
	f := relationDB.DataOpenAccessFilter{
		TenantCode: in.TenantCode,
		UserID:     in.UserID,
		Code:       in.Code,
	}
	total, err := relationDB.NewDataOpenAccessRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewDataOpenAccessRepo(l.ctx).FindByFilter(l.ctx, f, utils.Copy[stores.PageInfo](in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.OpenAccessIndexResp{
		Total: total,
		List:  utils.CopySlice[sys.OpenAccess](pos),
	}, nil
}
