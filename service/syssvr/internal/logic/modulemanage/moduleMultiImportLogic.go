package modulemanagelogic

import (
	"context"
	"encoding/json"
	tenantmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMultiImportLogic {
	return &ModuleMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMultiImportLogic) ModuleMultiImport(in *sys.ModuleMultiImportReq) (*sys.ModuleMultiImportResp, error) {
	var modules []*relationDB.SysModuleInfo
	err := json.Unmarshal([]byte(in.Modules), &modules)
	if err != nil {
		return nil, err
	}
	var resp sys.ModuleMultiImportResp
	for _, m := range modules {
		_, err := NewModuleInfoCreateLogic(l.ctx, l.svcCtx).ModuleInfoCreate(utils.Copy[sys.ModuleInfo](m))
		if err != nil && !errors.Cmp(err, errors.Duplicate) {
			l.Error(m, err)
			resp.ErrCount++
			continue
		}
		err = relationDB.NewAppModuleRepo(l.ctx).Insert(l.ctx, &relationDB.SysAppModule{
			AppCode:    "core",
			ModuleCode: m.Code,
		})
		if err != nil && !errors.Cmp(err, errors.Duplicate) {
			l.Error(m, err)
		}
		_, err = tenantmanagelogic.NewTenantAppModuleCreateLogic(l.ctx, l.svcCtx).TenantAppModuleCreate(&sys.TenantModuleCreateReq{
			Code:       def.TenantCodeDefault,
			AppCode:    def.AppCore,
			ModuleCode: m.Code,
		})
		if err != nil {
			l.Error(m, err)
		}
		resp.SuccCount++
		if len(m.Menus) > 0 {
			info := genMenuTree(m.Menus)
			err := NewModuleMenuMultiImportLogic(l.ctx, l.svcCtx).menuImport(m.Code, 3, info)
			if err != nil {
				l.Error(m, err)
			}
		}
	}
	return &resp, nil
}
