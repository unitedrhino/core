package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.AppInfoIndexReq) (resp *types.AppInfoIndexResp, err error) {
	ret, err := l.svcCtx.AppRpc.AppInfoIndex(l.ctx, &sys.AppInfoIndexReq{
		Name: req.Name,
		Page: logic.ToSysPageRpc(req.Page),
		Code: req.Code,
	})
	if err != nil {
		return nil, err
	}
	return &types.AppInfoIndexResp{
		Total: ret.Total,
		List:  ToAppInfosTypes(ret.List),
	}, nil
}
