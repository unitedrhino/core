package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

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
	ur, err := relationDB.NewDeptUserRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptUserFilter{UserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	if len(ur) == 0 {
		return &sys.UserDeptIndexResp{}, nil
	}
	var deptIDs []int64
	for _, v := range ur {
		deptIDs = append(deptIDs, v.DeptID)
	}
	rs, err := relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptInfoFilter{IDs: deptIDs}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.UserDeptIndexResp{List: utils.CopySlice[sys.DeptInfo](rs), Total: int64(len(deptIDs))}, nil
}
