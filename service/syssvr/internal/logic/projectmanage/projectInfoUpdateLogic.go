package projectmanagelogic

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"
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
	oldDevCount := po.DeviceCount
	err = l.setPoByPb(po, in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ProjectCache.SetData(l.ctx, in.ProjectID, nil)
	if err != nil {
		l.Error(err)
	}
	if oldDevCount != po.DeviceCount {
		err = relationDB.NewUserInfoRepo(l.ctx).UpdateDeviceCount(l.ctx, po.AdminUserID)
		if err != nil {
			l.Error(err)
		}
	}
	return &sys.Empty{}, nil
}
func (l *ProjectInfoUpdateLogic) setPoByPb(po *relationDB.SysProjectInfo, pb *sys.ProjectInfo) error {
	if pb.ProjectName != "" {
		po.ProjectName = pb.ProjectName
	}
	as, err := handAttachment(l.ctx, l.svcCtx, oss.BusinessProject, int64(po.ProjectID), po.Attachments, pb.Attachments)
	if err != nil {
		return err
	}
	po.Attachments = as
	//if pb.CompanyName != nil {
	//	po.CompanyName = pb.CompanyName.GetValue()
	//}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if uc.IsAdmin {
		if pb.DeviceCount != nil {
			po.DeviceCount = pb.DeviceCount.Value
		}
		if pb.DeviceOnlineCount != nil {
			po.DeviceOnlineCount = pb.DeviceOnlineCount.Value
		}
		if pb.AlarmStatus != 0 {
			po.AlarmStatus = pb.AlarmStatus
		}
		if pb.Status != 0 {
			po.Status = pb.Status
		}
		if pb.Type != "" {
			po.Type = pb.Type
		}
	}
	if pb.Sort != 0 {
		po.Sort = pb.Sort
	}
	if pb.Tags != nil {
		po.Tags = pb.Tags
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
	return err
}

func handAttachment(ctx context.Context, svcCtx *svc.ServiceContext, business string, id int64, old []*relationDB.Attachment, in []*sys.Attachment) (ret []*relationDB.Attachment, err error) {
	var oldA = map[int64]*relationDB.Attachment{}
	for _, v := range old {
		oldA[v.ID] = v
	}
	if len(in) != 0 {
		var up []*relationDB.Attachment

		for _, attachment := range in {
			if attachment.FilePath != "" {
				newPath := oss.GenFilePath(ctx, svcCtx.Config.Name, business, "attachment", fmt.Sprintf("%d/%d/%s", id, attachment.Id, oss.GetFileNameWithPath(attachment.FilePath)))
				path, err := svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(attachment.FilePath, newPath)
				if err != nil {
					return nil, errors.System.AddDetail(err)
				}
				up = append(up, &relationDB.Attachment{
					ID:       attachment.Id,
					FilePath: path,
					UseBy:    attachment.UseBy,
				})
			} else {
				o, ok := oldA[attachment.Id]
				if ok { //如果存在则直接把老的复制进去即可,如果没查到,则丢弃
					o.UseBy = attachment.UseBy
					up = append(up, o)
				}
				delete(oldA, attachment.Id)
			}
		}
		if len(up) != 0 {
			ret = up
		}
	} else {
		ret = []*relationDB.Attachment{}
	}
	if len(oldA) != 0 {
		defer func() {
			if err == nil {
				for _, v := range oldA {
					er := svcCtx.OssClient.Delete(ctx, v.FilePath, common.OptionKv{})
					if er != nil {
						logx.WithContext(ctx).Error(v, er)
					}
				}
			}
		}()
	}
	return
}
