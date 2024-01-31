package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuIndexLogic {
	return &TenantAppMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuIndexLogic) TenantAppMenuIndex(in *sys.TenantAppMenuIndexReq) (*sys.TenantAppMenuIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil && in.Code != "" {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	f := relationDB.TenantAppMenuFilter{
		ModuleCode: in.ModuleCode,
		TenantCode: in.Code,
		AppCode:    in.AppCode,
		MenuIDs:    in.MenuIDs,
	}
	resp, err := relationDB.NewTenantAppMenuRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	//total, err := relationDB.NewMenuInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	//if err != nil {
	//	return nil, err
	//}
	info := make([]*sys.TenantAppMenu, 0, len(resp))
	if !in.IsRetTree {
		for _, v := range resp {
			i := logic.ToTenantAppMenuInfoPb(v)
			info = append(info, i)
		}
		return &sys.TenantAppMenuIndexResp{List: info}, nil
	}

	var (
		pidMap = make(map[int64][]*sys.TenantAppMenu, len(resp))
		idMap  = make(map[int64]*sys.TenantAppMenu, len(resp))
	)
	for _, v := range resp {
		i := logic.ToTenantAppMenuInfoPb(v)
		idMap[i.Info.Id] = i
		if i.Info.ParentID == 1 { //根节点
			info = append(info, i)
			continue
		}
		pidMap[i.Info.ParentID] = append(pidMap[i.Info.ParentID], i)
	}
	fillChildren(info, pidMap)
	return &sys.TenantAppMenuIndexResp{List: info}, nil
}

func fillChildren(in []*sys.TenantAppMenu, pidMap map[int64][]*sys.TenantAppMenu) {
	for _, v := range in {
		id := v.TemplateID
		if id == 0 { //如果id为0,则是自定义节点,自定义节点的父节点是id
			id = v.Info.Id
		}
		cs := pidMap[id]
		if cs != nil {
			v.Children = cs
			fillChildren(cs, pidMap)
		}
	}
}
