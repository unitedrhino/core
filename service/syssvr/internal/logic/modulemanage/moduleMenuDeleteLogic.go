package modulemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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

func deleteMenu(ctx context.Context, tx *gorm.DB, id []int64) error {
	children, err := relationDB.NewMenuInfoRepo(tx).FindByFilter(ctx, relationDB.MenuInfoFilter{ParentIDs: id}, nil)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		var cids []int64
		for _, v := range children {
			cids = append(cids, v.ID)
		}
		err := deleteMenu(ctx, tx, cids)
		if err != nil {
			return err
		}
	}
	err = relationDB.NewMenuInfoRepo(tx).DeleteByFilter(ctx, relationDB.MenuInfoFilter{MenuIDs: id})
	if err != nil {
		return err
	}
	err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(ctx, relationDB.TenantAppMenuFilter{TempLateIDs: id})
	if err != nil {
		return err
	}
	return nil
}

func (l *ModuleMenuDeleteLogic) ModuleMenuDelete(in *sys.WithID) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		return deleteMenu(l.ctx, tx, []int64{in.Id})
	})
	return &sys.Empty{}, err
}
