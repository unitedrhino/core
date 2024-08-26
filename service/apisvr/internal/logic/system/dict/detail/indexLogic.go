package detail

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/dict"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.DictDetailIndexReq) (resp *types.DictDetailIndexResp, err error) {
	ret, err := l.svcCtx.DictM.DictDetailIndex(l.ctx, utils.Copy[sys.DictDetailIndexReq](req))
	if err != nil {
		return nil, err
	}

	return &types.DictDetailIndexResp{
		Total: ret.Total,
		List:  dict.ToDetailsTypes(ret.List),
	}, nil
}
