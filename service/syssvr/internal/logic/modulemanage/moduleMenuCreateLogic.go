package modulemanagelogic

import (
	"context"
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

func createMenu(ctx context.Context, tx *gorm.DB, in *sys.MenuInfo) (int64, error) {
	in.Id = 0
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
	ams, err := relationDB.NewTenantAppModuleRepo(tx).FindByFilter(ctx, relationDB.TenantAppModuleFilter{
		ModuleCodes: []string{in.ModuleCode},
	}, nil)
	if err != nil {
		return 0, err
	}
	po := utils.Copy[relationDB.SysModuleMenu](in)
	err = relationDB.NewMenuInfoRepo(tx).Insert(ctx, po)
	if err != nil {
		return 0, err
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
	err = relationDB.NewTenantAppMenuRepo(tx).MultiInsert(ctx, data)
	if err != nil {
		return 0, err
	}
	return po.ID, nil

}

func (l *ModuleMenuCreateLogic) ModuleMenuCreate(in *sys.MenuInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if err := CheckModule(l.ctx, in.ModuleCode); err != nil {
		return nil, err
	}
	var id int64
	var err error
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		id, err = createMenu(l.ctx, tx, in)
		return err
	})
	return &sys.WithID{Id: id}, err
}
