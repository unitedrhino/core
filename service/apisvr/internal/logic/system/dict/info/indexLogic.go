package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/dict"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

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

func (l *IndexLogic) Index(req *types.DictInfoIndexReq) (resp *types.DictInfoIndexResp, err error) {
	ret, err := l.svcCtx.DictM.DictInfoIndex(l.ctx, utils.Copy[sys.DictInfoIndexReq](req))
	if err != nil {
		return nil, err
	}

	return &types.DictInfoIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     dict.ToInfosTypes(ret.List),
	}, nil
}
