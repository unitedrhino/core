package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuDeleteLogic {
	return &ModuleMenuDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuDeleteLogic) ModuleMenuDelete(in *sys.WithID) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewMenuInfoRepo(l.ctx).Delete(l.ctx, in.Id)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.TenantAppMenuFilter{TempLateID: in.Id})
		if err != nil {
			return err
		}
		return nil
	})
	return &sys.Response{}, err
}
