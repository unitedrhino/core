package project

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectMultiUpdateLogic {
	return &ProjectMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectMultiUpdateLogic) ProjectMultiUpdate(req *types.DataProjectMultiUpdateReq) error {
	dto := &sys.DataProjectMultiUpdateReq{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Projects:   ToProjectPbs(req.Projects),
	}
	_, err := l.svcCtx.DataM.DataProjectMultiUpdate(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.DataProjectMultiUpdate req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
