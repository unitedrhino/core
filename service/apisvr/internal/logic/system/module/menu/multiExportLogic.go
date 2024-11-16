package menu

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量导出菜单
func NewMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiExportLogic {
	return &MultiExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiExportLogic) MultiExport(req *types.MenuMultiExportReq) (resp *types.MenuMultiExportResp, err error) {
	ret, err := l.svcCtx.ModuleRpc.ModuleMenuMultiExport(l.ctx, utils.Copy[sys.MenuMultiExportReq](req))

	return utils.Copy[types.MenuMultiExportResp](ret), err
}
