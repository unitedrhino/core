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

type LoginIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginIndexLogic {
	return &LoginIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginIndexLogic) LoginIndex(req *types.SysLogLoginIndexReq) (resp *types.SysLogLoginIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	info, err := l.svcCtx.LogRpc.LoginLogIndex(l.ctx, &sys.LoginLogIndexReq{
		AppCode:       req.AppCode,
		Page:          logic.ToSysPageRpc(req.Page),
		IpAddr:        req.IpAddr,
		LoginLocation: req.LoginLocation,
		Date:          &sys.DateRange{Start: req.DateRange.Start, End: req.DateRange.End},
		UserID:        req.UserID,
		UserName:      req.UserName,
		Code:          req.Code,
	})
	if err != nil {
		return nil, err
	}
	return utils.Copy[types.SysLogLoginIndexResp](info), err
}
