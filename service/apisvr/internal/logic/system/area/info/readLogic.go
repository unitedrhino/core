package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.AreaInfoReadReq) (resp *types.AreaInfo, err error) {
	dmResp, err := l.svcCtx.AreaM.AreaInfoRead(l.ctx, &sys.AreaInfoReadReq{AreaID: req.AreaID, ProjectID: req.ProjectID, WithChildren: req.WithChildren})
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return ToAreaInfoTypes(dmResp), nil
}
