// 模块菜单批量导入逻辑：支持按 path 增量/全量同步，并在全量导入后清理孤儿租户菜单，避免侧栏重复。
package modulemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/module"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	resp              sys.MenuMultiImportResp
	keptModuleMenuIDs map[int64]struct{} // 本次导入保留的模块菜单 ID（全量导入后用于删除多余节点）
}

func NewModuleMenuMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuMultiImportLogic {
	return &ModuleMenuMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func menuFmt(moduleCode string, in []*sys.MenuInfo) error {
	var path = map[string]struct{}{}
	for _, m := range in {
		m.Id = 0
		if m.ModuleCode != moduleCode {
			return errors.Parameter.AddMsg("导入和导出的模块编码需要一致")
		}
		if _, ok := path[m.Path]; ok {
			return errors.Parameter.AddMsgf("路由:%s 重复", m.Path)
		}
		path[m.Path] = struct{}{}
		if len(m.Children) > 0 {
			err := menuFmt(moduleCode, m.Children)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *ModuleMenuMultiImportLogic) Handle(tx *gorm.DB, ModuleCode string, Mode int64, parentID int64, menus []*sys.MenuInfo) error {
	l.resp.Total += int64(len(menus))
	db := relationDB.NewMenuInfoRepo(tx)
	paths := utils.ToSliceWithFunc(menus, func(in *sys.MenuInfo) string {
		return in.Path
	})
	olds, err := db.FindByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: ModuleCode, ParentID: parentID, Paths: paths}, nil)
	if err != nil {
		return err
	}
	var oldMap = map[string]*relationDB.SysModuleMenu{}
	for _, old := range olds {
		oldMap[old.Path] = old
	}
	for _, menu := range menus {
		menu.ParentID = parentID
		old, ok := oldMap[menu.Path]
		if !ok { //需要新建
			id, err := createMenu(l.ctx, tx, menu)
			if err != nil {
				return err
			}
			l.markKeptModuleMenu(id)
			if len(menu.Children) > 0 {
				err = l.Handle(tx, ModuleCode, Mode, id, menu.Children)
				if err != nil {
					return err
				}
			}
			continue
		}
		l.markKeptModuleMenu(old.ID)
		if Mode != module.MenuImportModeAdd { //如果只新增则不用处理这条
			err = updateMenu(l.ctx, tx, menu, old)
			if err != nil {
				return err
			}
		}
		if len(menu.Children) > 0 {
			err = l.Handle(tx, ModuleCode, Mode, old.ID, menu.Children)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *ModuleMenuMultiImportLogic) ModuleMenuMultiImport(in *sys.MenuMultiImportReq) (*sys.MenuMultiImportResp, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	var dos []*sys.MenuInfo
	err := json.Unmarshal([]byte(in.Menu), &dos)
	if err != nil {
		return nil, errors.Parameter.AddMsg("导入的菜单格式不对").AddDetail(err)
	}
	err = l.menuImport(in.ModuleCode, in.Mode, dos)

	return &l.resp, err
}
func (l *ModuleMenuMultiImportLogic) menuImport(ModuleCode string, Mode int64, dos []*sys.MenuInfo) error {
	err := menuFmt(ModuleCode, dos)
	if err != nil {
		return err
	}
	if err := CheckModule(l.ctx, ModuleCode); err != nil {
		return err
	}

	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if Mode == module.MenuImportModeAll {
			l.keptModuleMenuIDs = make(map[int64]struct{})
			// 先清理历史孤儿租户菜单，避免 template_id 失效导致的侧栏重复
			if err := cleanupOrphanTenantMenus(l.ctx, tx, ModuleCode); err != nil {
				return err
			}
		}
		err := l.Handle(tx, ModuleCode, Mode, def.RootNode, dos)
		if err != nil {
			return err
		}
		if Mode == module.MenuImportModeAll {
			if err := l.deleteModuleMenusExcept(l.ctx, tx, ModuleCode); err != nil {
				return err
			}
			if err := cleanupOrphanTenantMenus(l.ctx, tx, ModuleCode); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// markKeptModuleMenu 记录本次导入应保留的模块菜单 ID。
func (l *ModuleMenuMultiImportLogic) markKeptModuleMenu(id int64) {
	if l.keptModuleMenuIDs == nil {
		return
	}
	l.keptModuleMenuIDs[id] = struct{}{}
}

// deleteModuleMenusExcept 删除本次导入未覆盖的模块菜单及其租户副本。
func (l *ModuleMenuMultiImportLogic) deleteModuleMenusExcept(ctx context.Context, tx *gorm.DB, moduleCode string) error {
	menus, err := relationDB.NewMenuInfoRepo(tx).FindByFilter(ctx, relationDB.MenuInfoFilter{ModuleCode: moduleCode}, nil)
	if err != nil {
		return err
	}
	for _, m := range menus {
		if _, ok := l.keptModuleMenuIDs[m.ID]; ok {
			continue
		}
		if err := deleteMenu(ctx, tx, []int64{m.ID}); err != nil {
			return err
		}
	}
	return nil
}

// withAllTenant 在跨租户清理数据时临时放开租户隔离。
func withAllTenant(ctx context.Context) context.Context {
	uc := ctxs.GetUserCtx(ctx)
	if uc != nil {
		uc.AllTenant = true
	}
	return ctx
}

// cleanupOrphanTenantMenus 清理模块下 template_id 已失效或同 path 重复的租户菜单。
func cleanupOrphanTenantMenus(ctx context.Context, tx *gorm.DB, moduleCode string) error {
	ctx = withAllTenant(ctx)

	moduleMenus, err := relationDB.NewMenuInfoRepo(tx).FindByFilter(ctx, relationDB.MenuInfoFilter{ModuleCode: moduleCode}, nil)
	if err != nil {
		return err
	}
	validTemplateIDs := make(map[int64]struct{}, len(moduleMenus))
	for _, m := range moduleMenus {
		validTemplateIDs[m.ID] = struct{}{}
	}

	tenantMenus, err := relationDB.NewTenantAppMenuRepo(tx).FindByFilter(ctx, relationDB.TenantAppMenuFilter{ModuleCode: moduleCode}, nil)
	if err != nil {
		return err
	}

	type menuKey struct {
		tenantCode string
		appCode    string
		parentID   int64
		path       string
	}
	grouped := make(map[menuKey][]*relationDB.SysTenantAppMenu)
	var orphanIDs []int64

	for _, tm := range tenantMenus {
		if _, ok := validTemplateIDs[tm.TempLateID]; !ok {
			orphanIDs = append(orphanIDs, tm.ID)
			continue
		}
		k := menuKey{
			tenantCode: string(tm.TenantCode),
			appCode:    tm.AppCode,
			parentID:   tm.ParentID,
			path:       tm.Path,
		}
		grouped[k] = append(grouped[k], tm)
	}

	for _, list := range grouped {
		if len(list) <= 1 {
			continue
		}
		keepID := list[0].ID
		keepTemplateID := list[0].TempLateID
		for _, tm := range list[1:] {
			if tm.TempLateID > keepTemplateID {
				orphanIDs = append(orphanIDs, keepID)
				keepID = tm.ID
				keepTemplateID = tm.TempLateID
				continue
			}
			orphanIDs = append(orphanIDs, tm.ID)
		}
	}

	if len(orphanIDs) == 0 {
		return nil
	}
	if err := deleteTenantMenusAndRoleRefs(ctx, tx, moduleCode, orphanIDs); err != nil {
		return err
	}
	return nil
}

// deleteTenantMenusAndRoleRefs 删除租户菜单并清理角色菜单引用。
func deleteTenantMenusAndRoleRefs(ctx context.Context, tx *gorm.DB, moduleCode string, menuIDs []int64) error {
	if len(menuIDs) == 0 {
		return nil
	}
	err := tx.WithContext(ctx).
		Where("module_code = ? AND menu_id IN ?", moduleCode, menuIDs).
		Delete(&relationDB.SysRoleMenu{}).Error
	if err != nil {
		return stores.ErrFmt(err)
	}
	err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(ctx, relationDB.TenantAppMenuFilter{
		ModuleCode: moduleCode,
		MenuIDs:    menuIDs,
	})
	return err
}
