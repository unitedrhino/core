package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppCreateLogic {
	return &TenantAppCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppCreateLogic) TenantAppCreate(in *sys.TenantAppInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		//todo 需要检查租户是否存在
		err := relationDB.NewTenantAppRepo(tx).Insert(l.ctx, &relationDB.SysTenantApp{
			TenantCode:     dataType.TenantCode(in.Code),
			AppCode:        in.AppCode,
			WxMini:         utils.Copy[relationDB.SysTenantThird](in.WxMini),
			WxOpen:         utils.Copy[relationDB.SysTenantThird](in.WxOpen),
			DingMini:       utils.Copy[relationDB.SysTenantThird](in.DingMini),
			Android:        utils.Copy[relationDB.SysThirdApp](in.Android),
			IsAutoRegister: in.IsAutoRegister,
			Config:         in.Config,
		})
		if err != nil {
			return err
		}
		for _, module := range in.Modules {
			err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, module.Code, module.MenuIDs)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &sys.Empty{}, err
}
func ModuleCreate(ctx context.Context, tx *gorm.DB, tenantCode, appCode string, moduleCode string, menuIDs []int64) error {
	if _, err := relationDB.NewTenantAppModuleRepo(tx).FindOneByFilter(ctx,
		relationDB.TenantAppModuleFilter{TenantCode: tenantCode, AppCode: appCode, ModuleCodes: []string{moduleCode}}); err == nil || !errors.Cmp(err, errors.NotFind) { //如果报错或者已经有了则跳过
		return err
	}
	mi, err := relationDB.NewModuleInfoRepo(tx).FindOneByFilter(ctx,
		relationDB.ModuleInfoFilter{Codes: []string{moduleCode}, WithMenus: true})
	if err != nil {
		return err
	}
	var (
		menuMap = make(map[int64]*relationDB.SysModuleMenu)
		//menuTree = make(map[int64]*relationDB.SysModuleMenu)
		allMenu = false
	)
	if len(menuIDs) == 0 {
		allMenu = true
	}
	for _, m := range mi.Menus {
		if allMenu {
			if m.IsAllTenant == def.False { //如果菜单不是给所有租户用的则跳过
				continue
			}
			menuIDs = append(menuIDs, m.ID)
		}
		menuMap[m.ID] = m
	}
	var (
		insertMenus []*relationDB.SysTenantAppMenu
	)

	for _, id := range menuIDs {
		m := menuMap[id]
		if m == nil { //模板里不存在无法添加
			continue
		}
		v := utils.Copy[relationDB.SysTenantAppMenu](m)
		v.TempLateID = m.ID
		v.TenantCode = dataType.TenantCode(tenantCode)
		v.AppCode = appCode
		v.ID = 0
		insertMenus = append(insertMenus, v)
	}
	err = relationDB.NewTenantAppMenuRepo(tx).MultiInsert(ctx, insertMenus)
	if err != nil {
		return err
	}
	err = relationDB.NewTenantAppModuleRepo(tx).Insert(ctx, &relationDB.SysTenantAppModule{
		TenantCode: dataType.TenantCode(tenantCode), SysAppModule: relationDB.SysAppModule{AppCode: appCode, ModuleCode: moduleCode}})
	return err
}
