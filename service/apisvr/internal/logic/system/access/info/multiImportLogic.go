package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiImportLogic {
	return &MultiImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiImportLogic) MultiImport(req *types.AccessMultiImportReq, f []byte) (resp *types.AccessMultiImportResp, err error) {
	ret, err := l.svcCtx.AccessRpc.AccessInfoMultiImport(l.ctx, &sys.AccessInfoMultiImportReq{
		Module: req.Module,
		Access: string(f),
	})
	if err != nil {
		return nil, err
	}
	return &types.AccessMultiImportResp{
		Total:       ret.Total,
		ErrCount:    ret.ErrCount,
		IgnoreCount: ret.IgnoreCount,
		SuccCount:   ret.SuccCount,
	}, nil
}
