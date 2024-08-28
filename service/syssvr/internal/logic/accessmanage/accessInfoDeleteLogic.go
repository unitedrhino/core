package accessmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoDeleteLogic {
	return &AccessInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoDeleteLogic) AccessInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	po, err := relationDB.NewAccessRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewAccessRepo(l.ctx).Delete(l.ctx, in.Id)
		if err != nil {
			return err
		}
		err = relationDB.NewApiInfoRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.ApiInfoFilter{
			AccessCode: po.Code,
		})
		return err
	})
	return &sys.Empty{}, err
}
