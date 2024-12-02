package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeptMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeptMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeptMultiCreateLogic {
	return &UserDeptMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeptMultiCreateLogic) UserDeptMultiCreate(in *sys.UserDeptMultiSaveReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if len(in.DeptIDs) != 0 {
		rs, err := relationDB.NewDeptInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeptInfoFilter{
			IDs: in.DeptIDs,
		})
		if err != nil {
			return nil, err
		}
		if int(rs) != len(in.DeptIDs) {
			return nil, errors.Parameter.WithMsg("有部门不存咋")
		}
	}
	var datas []*relationDB.SysDeptUser
	for _, v := range in.DeptIDs {
		datas = append(datas, &relationDB.SysDeptUser{
			DeptID: v,
			UserID: in.UserID,
		})
	}
	err := relationDB.NewDeptUserRepo(l.ctx).MultiInsert(l.ctx, datas)
	if err == nil {
		l.svcCtx.UserTokenInfo.SetData(l.ctx, in.UserID, nil)
	}
	return &sys.Empty{}, err
}
