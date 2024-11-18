package modulemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
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
	var (
		pidMap = make(map[int64][]*sys.MenuInfo, len(pos))
		idMap  = make(map[int64]*sys.MenuInfo, len(pos))
		info   = make([]*sys.MenuInfo, 0, len(pos))
	)

	for _, v := range pos {
		i := logic.ToMenuInfoPb(v)
		idMap[i.Id] = i
		if i.ParentID == 1 { //根节点
			info = append(info, i)
			continue
		}
		pidMap[i.ParentID] = append(pidMap[i.ParentID], i)
	}
	fillChildren(info, pidMap)
	return &sys.MenuMultiExportResp{Menu: utils.MarshalNoErr(info)}, nil
}
