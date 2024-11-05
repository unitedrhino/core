package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeptMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeptMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeptMultiDeleteLogic {
	return &UserDeptMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeptMultiDeleteLogic) UserDeptMultiDelete(in *sys.UserDeptMultiSaveReq) (*sys.Empty, error) {
	// todo: add your logic here and delete this line

	return &sys.Empty{}, nil
}
