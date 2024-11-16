package modulemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuCreateLogic {
	return &ModuleMenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuCreateLogic) ModuleMenuCreate(in *sys.MenuInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if err := CheckModule(l.ctx, in.ModuleCode); err != nil {
		return nil, err
	}
	if in.Type == 0 {
		in.Type = 1
	}
	if in.ParentID == 0 {
		in.ParentID = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	ams, err := relationDB.NewTenantAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAppModuleFilter{
		ModuleCodes: []string{in.ModuleCode},
	}, nil)
	if err != nil {
		return nil, err
	}
	po := logic.ToMenuInfoPo(in)
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewMenuInfoRepo(l.ctx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		//导入租户的菜单中
		var data []*relationDB.SysTenantAppMenu
		var template = *po
		template.ID = 0
		for _, am := range ams {
			tam := utils.Copy[relationDB.SysTenantAppMenu](template)
			tam.TempLateID = po.ID
			tam.TenantCode = am.TenantCode
			tam.AppCode = am.AppCode
			data = append(data, tam)
		}
		err = relationDB.NewTenantAppMenuRepo(l.ctx).MultiInsert(l.ctx, data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
