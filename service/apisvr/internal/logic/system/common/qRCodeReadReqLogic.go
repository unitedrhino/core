package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QRCodeReadReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQRCodeReadReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QRCodeReadReqLogic {
	return &QRCodeReadReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QRCodeReadReqLogic) QRCodeReadReq(req *types.QRCodeReadReq) (resp *types.QRCodeReadResp, err error) {
	ret, err := l.svcCtx.Common.QRCodeRead(l.ctx, utils.Copy[sys.QRCodeReadReq](req))
	if err != nil {
		return nil, err
	}
	return &types.QRCodeReadResp{Buffer: ret.Buffer}, nil
}
