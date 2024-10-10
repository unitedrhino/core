package area

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.DataAreaMultiUpdateReq) error {
	dto := &sys.DataAreaMultiUpdateReq{
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		ProjectID:  req.ProjectID,
		Areas:      ToAreaPbs(req.Areas),
	}
	_, err := l.svcCtx.DataM.DataAreaMultiUpdate(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserDataAuthManage req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
