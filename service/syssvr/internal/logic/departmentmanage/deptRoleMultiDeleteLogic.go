package departmentmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptRoleMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptRoleMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptRoleMultiDeleteLogic {
	return &DeptRoleMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptRoleMultiDeleteLogic) DeptRoleMultiDelete(in *sys.DeptRoleMultiSaveReq) (*sys.Empty, error) {
	// todo: add your logic here and delete this line

	return &sys.Empty{}, nil
}
