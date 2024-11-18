package modulemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/domain/module"
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
	resp sys.MenuMultiImportResp
}

func NewModuleMenuMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuMultiImportLogic {
	return &ModuleMenuMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuMultiImportLogic) menuFmt(moduleCode string, in []*sys.MenuInfo) error {
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
			err := l.menuFmt(moduleCode, m.Children)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *ModuleMenuMultiImportLogic) Handle(tx *gorm.DB, in *sys.MenuMultiImportReq, parentID int64, menus []*sys.MenuInfo) error {
	l.resp.Total += int64(len(menus))
	db := relationDB.NewMenuInfoRepo(tx)
	paths := utils.ToSliceWithFunc(menus, func(in *sys.MenuInfo) string {
		return in.Path
	})
	olds, err := db.FindByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: in.ModuleCode, ParentID: parentID, Paths: paths}, nil)
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
			if len(menu.Children) > 0 {
				err = l.Handle(tx, in, id, menu.Children)
				if err != nil {
					return err
				}
			}
			continue
		}
		if in.Mode != module.MenuImportModeAdd { //如果只新增则不用处理这条
			err = updateMenu(l.ctx, tx, menu, old)
			if err != nil {
				return err
			}
		}
		if len(menu.Children) > 0 {
			err = l.Handle(tx, in, old.ID, menu.Children)
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
	err = l.menuFmt(in.ModuleCode, dos)
	if err != nil {
		return nil, err
	}
	if err := CheckModule(l.ctx, in.ModuleCode); err != nil {
		return nil, err
	}

	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if in.Mode == module.MenuImportModeAll {
			db := relationDB.NewMenuInfoRepo(tx)
			err = db.DeleteByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: in.ModuleCode})
			if err != nil {
				return err
			}
			err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantAppMenuFilter{ModuleCode: in.ModuleCode})
			if err != nil {
				return err
			}
		}
		err := l.Handle(tx, in, def.RootNode, dos)
		if err != nil {
			return err
		}
		return nil
	})

	return &l.resp, err
}
