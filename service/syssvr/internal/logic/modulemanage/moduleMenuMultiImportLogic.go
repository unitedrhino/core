package modulemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/domain/module"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	var pos []*relationDB.SysModuleMenu
	err := json.Unmarshal([]byte(in.Menu), &pos)
	if err != nil {
		return nil, errors.Parameter.AddMsg("导入的菜单格式不对").AddDetail(err)
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		switch in.Mode {
		case module.MenuImportModeUpdate:
			err = relationDB.NewMenuInfoRepo(tx).MultiInsert(l.ctx, pos)
		case module.MenuImportModeAll:
			db := relationDB.NewMenuInfoRepo(tx)
			err = db.DeleteByFilter(l.ctx, relationDB.MenuInfoFilter{ModuleCode: in.ModuleCode})
			if err != nil {
				return err
			}
			err = db.MultiInsert(l.ctx, pos)
			if err != nil {
				return err
			}
			err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantAppMenuFilter{ModuleCode: in.ModuleCode})
			if err != nil {
				return err
			}
		default:
			err = relationDB.NewMenuInfoRepo(tx).MultiInsertOnly(l.ctx, pos)
		}
		if err != nil {
			return err
		}
		ams, err := relationDB.NewTenantAppModuleRepo(tx).FindByFilter(l.ctx, relationDB.TenantAppModuleFilter{
			ModuleCodes: []string{in.ModuleCode},
		}, nil)
		if err != nil {
			return err
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
			err = relationDB.NewTenantAppMenuRepo(tx).MultiInsert(l.ctx, data)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &sys.MenuMultiImportResp{Total: int64(len(pos))}, err
}
