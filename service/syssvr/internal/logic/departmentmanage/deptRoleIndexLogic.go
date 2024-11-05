package departmentmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptRoleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptRoleIndexLogic {
	return &DeptRoleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptRoleIndexLogic) DeptRoleIndex(in *sys.DeptRoleIndexReq) (*sys.DeptRoleIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.DeptRoleIndexResp{}, nil
}
