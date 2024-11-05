package departmentmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptUserIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptUserIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptUserIndexLogic {
	return &DeptUserIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptUserIndexLogic) DeptUserIndex(in *sys.DeptUserIndexReq) (*sys.DeptUserIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.DeptUserIndexResp{}, nil
}
