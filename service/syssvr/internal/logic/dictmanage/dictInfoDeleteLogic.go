package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoDeleteLogic {
	return &DictInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoDeleteLogic) DictInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	err := relationDB.NewDictInfoRepo(l.Info).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
