package projectmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/oss/common"

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
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllProject = true
	defer func() {
		uc.AllProject = false
	}()
	if in.ProjectID == 0 {
		if uc.ProjectID <= def.NotClassified {
			return nil, errors.Parameter
		}
		in.ProjectID = uc.ProjectID
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
	err = l.svcCtx.ProjectCache.SetData(l.ctx, in.ProjectID, nil)
	if err != nil {
		l.Error(err)
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
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if uc.IsAdmin {
		if pb.DeviceCount != nil {
			po.DeviceCount = pb.DeviceCount.Value
		}
	}
	if pb.AdminUserID != 0 && pb.AdminUserID != po.AdminUserID {
		uc := ctxs.GetUserCtx(l.ctx)
		if uc.UserID == po.AdminUserID { //只有项目的管理员才有权限更换
			po.AdminUserID = pb.AdminUserID
		}
	}
	if pb.Position != nil {
		po.Position = logic.ToStorePoint(pb.Position)
	}
	if pb.Area != nil {
		po.Area = pb.Area.GetValue()
	}
	if pb.Ppsm != 0 {
		po.Ppsm = pb.Ppsm
	}
	//if pb.Region != nil {
	//	po.Region = pb.Region.GetValue()
	//}
	if pb.Address != nil {
		po.Address = pb.Address.GetValue()
	}
	if pb.IsUpdateProjectImg && pb.ProjectImg != "" {
		if po.ProjectImg != "" {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, po.ProjectImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", po.ProjectImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProject, oss.SceneHeadIng,
			fmt.Sprintf("%d/%s", pb.ProjectID, oss.GetFileNameWithPath(pb.ProjectImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(pb.ProjectImg, nwePath)
		if err != nil {
			l.Error(err)
		} else {
			po.ProjectImg = path
		}

	}
	if pb.Desc != nil {
		po.Desc = pb.Desc.GetValue()
	}
}
