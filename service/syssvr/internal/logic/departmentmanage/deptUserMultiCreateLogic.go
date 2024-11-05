package departmentmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptUserMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptUserMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptUserMultiCreateLogic {
	return &DeptUserMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptUserMultiCreateLogic) DeptUserMultiCreate(in *sys.DeptUserMultiSaveReq) (*sys.Empty, error) {
	// todo: add your logic here and delete this line

	return &sys.Empty{}, nil
}
