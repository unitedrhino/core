package modulemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuMultiExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuMultiExportLogic {
	return &ModuleMenuMultiExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuMultiExportLogic) ModuleMenuMultiExport(in *sys.MenuMultiExportReq) (*sys.MenuMultiExportResp, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	pos, err := relationDB.NewMenuInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: in.ModuleCode}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.MenuMultiExportResp{Menu: utils.MarshalNoErr(pos)}, nil
}
