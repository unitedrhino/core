package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeptIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeptIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeptIndexLogic {
	return &UserDeptIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeptIndexLogic) UserDeptIndex(in *sys.UserDeptIndexReq) (*sys.UserDeptIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.UserDeptIndexResp{}, nil
}
