package feedback

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.OpsFeedbackIndexReq) (resp *types.OpsFeedbackIndexResp, err error) {
	ret, err := l.svcCtx.Ops.OpsFeedbackIndex(l.ctx, utils.Copy[sys.OpsFeedbackIndexReq](req))
	if err != nil {
		return nil, err
	}
	var list = utils.CopySlice[types.OpsFeedback](ret.List)
	for _, v := range list {
		u, err := l.svcCtx.UserCache.GetData(l.ctx, v.RaiseUserID)
		if err != nil {
			continue
		}
		v.User = utils.Copy[types.UserCore](u)
	}
	return &types.OpsFeedbackIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     list,
	}, nil
}
