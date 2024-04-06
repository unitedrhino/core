package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoUpdateLogic {
	return &ProjectInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 更新项目
func (l *ProjectInfoUpdateLogic) ProjectInfoUpdate(in *sys.ProjectInfo) (*sys.Empty, error) {
	ctxs.GetUserCtx(l.ctx).AllProject = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllProject = false
	}()
	if in.ProjectID == 0 {
		return nil, errors.Parameter
	}

	po, err := checkProject(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查项目出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
	}

	l.setPoByPb(po, in)

	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
func (l *ProjectInfoUpdateLogic) setPoByPb(po *relationDB.SysProjectInfo, pb *sys.ProjectInfo) {
	if pb.ProjectName != "" {
		po.ProjectName = pb.ProjectName
	}
	//if pb.CompanyName != nil {
	//	po.CompanyName = pb.CompanyName.GetValue()
	//}
	if pb.AdminUserID != 0 && pb.AdminUserID != po.AdminUserID {
		uc := ctxs.GetUserCtx(l.ctx)
		if uc.UserID == po.AdminUserID { //只有项目的管理员才有权限更换
			po.AdminUserID = pb.AdminUserID
		}
	}
	if pb.Position != nil {
		po.Position = logic.ToStorePoint(pb.Position)
	}
	//if pb.Region != nil {
	//	po.Region = pb.Region.GetValue()
	//}
	//if pb.Address != nil {
	//	po.Address = pb.Address.GetValue()
	//}
	if pb.Desc != nil {
		po.Desc = pb.Desc.GetValue()
	}
}
