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

func (l *IndexLogic) Index(req *types.AccessIndexReq) (resp *types.AccessIndexResp, err error) {
	ret, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Group:      req.Group,
		Code:       req.Code,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
		WithApis:   req.WithApis,
		AuthTypes:  req.AuthTypes,
	})
	if err != nil {
		return nil, err
	}
	return &types.AccessIndexResp{
		List:     ToAccessInfosTypes(ret.List),
		PageResp: logic.ToPageResp(req.Page, ret.Total),
	}, nil
}
