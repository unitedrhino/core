package access

import (
	"context"
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

// 获取开放认证列表
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DataOpenAccessIndexReq) (resp *types.DataOpenAccessIndexResp, err error) {
	ret, err := l.svcCtx.DataM.DataOpenAccessIndex(l.ctx, utils.Copy[sys.OpenAccessIndexReq](req))
	if err != nil {
		return nil, err
	}
	resp = utils.Copy[types.DataOpenAccessIndexResp](ret)
	for _, v := range resp.List {
		u, err := l.svcCtx.UserCache.GetData(l.ctx, v.UserID)
		if err != nil {
			continue
		}
		v.User = utils.Copy[types.UserCore](u)
	}
	return utils.Copy[types.DataOpenAccessIndexResp](ret), nil
}
