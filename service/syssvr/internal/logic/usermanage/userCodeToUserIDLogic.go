package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/users"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCodeToUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCodeToUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCodeToUserIDLogic {
	return &UserCodeToUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCodeToUserIDLogic) UserCodeToUserID(in *sys.UserCodeToUserIDReq) (*sys.UserCodeToUserIDResp, error) {
	switch in.LoginType {
	case users.RegWxMiniP:
		cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if cli.MiniProgram == nil {
			return nil, errors.System.AddMsg("未配置小程序")
		}
		auth := cli.MiniProgram.GetAuth()
		ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		return &sys.UserCodeToUserIDResp{
			OpenID:  ret.OpenID,
			UnionID: ret.UnionID,
		}, nil
	}
	return &sys.UserCodeToUserIDResp{}, errors.NotRealize
}
