package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptInfoDeleteLogic {
	return &DeptInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptInfoDeleteLogic) DeptInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewDeptInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewDeptInfoRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DeptInfoFilter{IDPath: old.IDPath})

	return &sys.Empty{}, nil
}
