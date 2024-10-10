package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.UserInfoCreateReq) (resp *types.UserCreateResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), req)
	info := req.Info
	//性别参数如果未传或者无效，则默认指定为男性
	if info.Sex != 1 && info.Sex != 2 {
		info.Sex = 1
	}
	resp1, err1 := l.svcCtx.UserRpc.UserInfoCreate(l.ctx, &sys.UserInfoCreateReq{
		Info:    user.UserInfoToRpc(info),
		RoleIDs: req.RoleIDs,
	})
	if err1 != nil {
		er := errors.Fmt(err1)
		l.Errorf("%s.rpc.Register req=%v err=%v rpc_err=%v", utils.FuncName(), req, er, err)
		return &types.UserCreateResp{}, er
	}
	if resp1 == nil {
		l.Errorf("%s.rpc.Register return nil req=%v", utils.FuncName(), req)
		return &types.UserCreateResp{}, errors.System.AddDetail("register core rpc return nil")
	}

	return &types.UserCreateResp{UserID: resp1.UserID}, nil
}
