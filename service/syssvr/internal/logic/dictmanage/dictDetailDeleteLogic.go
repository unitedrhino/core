package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailDeleteLogic {
	return &DictDetailDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailDeleteLogic) DictDetailDelete(in *sys.WithID) (*sys.Response, error) {
	err := relationDB.NewDictDetailRepo(l.ctx).Delete(l.ctx, in.Id)

	return &sys.Response{}, err
}
