package log

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/utils"

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
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total

	var logLoginInfo []*types.SysLogLoginInfo
	logLoginInfo = make([]*types.SysLogLoginInfo, 0, len(logLoginInfo))

	for _, i := range info.List {
		logLoginInfo = append(logLoginInfo, &types.SysLogLoginInfo{
			AppCode:       i.AppCode,
			UserID:        i.UserID,
			UserName:      i.UserName,
			IpAddr:        i.IpAddr,
			LoginLocation: i.LoginLocation,
			Browser:       i.Browser,
			Os:            i.Os,
			Code:          i.Code,
			Msg:           i.Msg,
			CreatedTime:   i.CreatedTime,
		})
	}

	return &types.SysLogLoginIndexResp{List: logLoginInfo, Total: total}, nil
}
