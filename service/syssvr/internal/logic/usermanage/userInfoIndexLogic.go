package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoIndexLogic {
	return &UserInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoIndexLogic) UserInfoIndex(in *sys.UserInfoIndexReq) (*sys.UserInfoIndexResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	f := relationDB.UserInfoFilter{
		UserName: in.UserName,
		Phone:    in.Phone,
		Email:    in.Email,
		UserIDs:  in.UserIDs,
	}
	ucs, err := l.UiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.UiDB.CountByFilter(l.ctx, f)
	info := make([]*sys.UserInfo, 0, len(ucs))
	for _, uc := range ucs {
		info = append(info, UserInfoToPb(l.ctx, uc, l.svcCtx))
	}
	if err != nil {
		return nil, err
	}
	return &sys.UserInfoIndexResp{
		List:  info,
		Total: total,
	}, nil

}
