package log

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperIndexLogic {
	return &OperIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperIndexLogic) OperIndex(req *types.SysLogOperIndexReq) (resp *types.SysLogOperIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	info, err := l.svcCtx.LogRpc.OperLogIndex(l.ctx, &sys.OperLogIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		AppCode:      req.AppCode,
		OperName:     req.OperName,
		OperUserName: req.OperUserName,
		BusinessType: req.BusinessType,
		Code:         req.Code,
		UserID:       req.UserID,
	})
	return utils.Copy[types.SysLogOperIndexResp](info), err
}
