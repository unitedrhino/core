package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserTaRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	ur     *UserRegisterLogic
	logx.Logger
}

func NewUserTaRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTaRegisterLogic {
	return &UserTaRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ur:     NewUserRegisterLogic(ctx, svcCtx),
	}
}

func (l *UserTaRegisterLogic) UserTaRegister(in *sys.UserTaRegisterReq) (*sys.UserRegisterResp, error) {
	if in.Password == "" {
		return nil, errors.Parameter.AddMsg("密码必填")
	}
	// 验证密码
	if err := l.ur.validatePassword(in.Password); err != nil {
		return nil, err
	}
	return &sys.UserRegisterResp{}, nil
}
