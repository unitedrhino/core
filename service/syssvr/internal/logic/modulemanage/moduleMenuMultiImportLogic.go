package modulemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/domain/module"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
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
}

func NewModuleMenuMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuMultiImportLogic {
	return &ModuleMenuMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuMultiImportLogic) ModuleMenuMultiImport(in *sys.MenuMultiImportReq) (*sys.MenuMultiImportResp, error) {
	var pos []*relationDB.SysModuleMenu
	err := json.Unmarshal([]byte(in.Menu), &pos)
	if err != nil {
		return nil, errors.Parameter.AddMsg("导入的菜单格式不对").AddDetail(err)
	}
	switch in.Mode {
	case module.MenuImportModeUpdate:
		err = relationDB.NewMenuInfoRepo(l.ctx).MultiInsertOnly(l.ctx, pos)
	case module.MenuImportModeAll:
		err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			db := relationDB.NewMenuInfoRepo(tx)
			err := db.DeleteByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: in.ModuleCode})
			if err != nil {
				return err
			}
			err = db.MultiInsert(l.ctx, pos)
			return err
		})
	default:
		err = relationDB.NewMenuInfoRepo(l.ctx).MultiInsertOnly(l.ctx, pos)
	}
	if err != nil {
		return nil, err
	}
	ams, err := relationDB.NewTenantAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAppModuleFilter{
		ModuleCodes: []string{in.ModuleCode},
	}, nil)
	if err != nil {
		return nil, err
	}
	var data []*relationDB.SysTenantAppMenu
	for _, po := range pos {
		var template = *po
		template.ID = 0
		for _, am := range ams {
			tam := utils.Copy[relationDB.SysTenantAppMenu](template)
			tam.TempLateID = po.ID
			tam.TenantCode = am.TenantCode
			tam.AppCode = am.AppCode
			data = append(data, tam)
		}
	}
	if len(data) > 0 {
		err = relationDB.NewTenantAppMenuRepo(l.ctx).MultiInsert(l.ctx, data)
		if err != nil {
			return nil, err
		}
	}
	return &sys.MenuMultiImportResp{}, nil
}
