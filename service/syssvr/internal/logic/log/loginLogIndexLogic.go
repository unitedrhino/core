package loglogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/stores"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	LlDB *relationDB.LoginLogRepo
}

func NewLoginLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogIndexLogic {
	return &LoginLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		LlDB:   relationDB.NewLoginLogRepo(ctx),
	}
}

func (l *LoginLogIndexLogic) LoginLogIndex(in *sys.LoginLogIndexReq) (*sys.LoginLogIndexResp, error) {
	f := relationDB.LoginLogFilter{
		IpAddr:        in.IpAddr,
		LoginLocation: in.LoginLocation,
		Data: &relationDB.DateRange{
			Start: in.Date.Start,
			End:   in.Date.End,
		},
		UserID:   in.UserID,
		UserName: in.UserName,
		Code:     in.Code,
	}
	resp, err := l.LlDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  2,
	}))
	if err != nil {
		return nil, err
	}
	total, err := l.LlDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.LoginLogInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.LoginLogInfo{
			AppCode:       v.AppCode,
			UserID:        v.UserID,
			UserName:      v.UserName,
			IpAddr:        v.IpAddr,
			LoginLocation: v.LoginLocation,
			Browser:       v.Browser,
			Os:            v.Os,
			Code:          v.Code,
			Msg:           v.Msg,
			CreatedTime:   v.CreatedTime.Unix(),
		})
	}

	return &sys.LoginLogIndexResp{List: info, Total: total}, nil
}
