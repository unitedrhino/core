package modulemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoCreateLogic {
	return &ModuleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoCreateLogic) ModuleInfoCreate(in *sys.ModuleInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if in.Type == 0 {
		in.Type = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	po := logic.ToModuleInfoPo(in)
	po.ID = 0
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewModuleInfoRepo(tx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		//自动添加到全部模块中
		err = relationDB.NewAppModuleRepo(tx).Insert(l.ctx, &relationDB.SysAppModule{
			AppCode:    def.AppAll,
			ModuleCode: in.Code,
		})
		return err
	})

	return &sys.WithID{Id: po.ID}, err
}
