package areamanagelogic

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
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"strings"
)

type AreaInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoUpdateLogic {
	return &AreaInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 更新区域
func (l *AreaInfoUpdateLogic) AreaInfoUpdate(in *sys.AreaInfo) (*sys.Empty, error) {
	if in.AreaID == 0 || utils.SliceIn(in.AreaID, def.RootNode, def.NotClassified) {
		return nil, errors.Parameter
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		areaPo, err := checkArea(l.ctx, tx, in.AreaID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查区域出错")
		} else if areaPo == nil {
			return errors.Parameter.AddDetail(in.AreaID).WithMsg("检查区域不存在")
		}

		projPo, err := checkProject(l.ctx, tx, in.ProjectID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查项目出错")
		} else if projPo == nil {
			return errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
		}

		if in.AreaName != "" && in.AreaName != areaPo.AreaName { //如果修改了区域名称
			names := utils.GetNamePath(areaPo.AreaNamePath)
			names[len(names)-1] = in.AreaName
			newAreaNamePath := strings.Join(names, "-") + "-"
			aiDB := relationDB.NewAreaInfoRepo(tx)
			areas, err := aiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDPath: areaPo.AreaIDPath}, nil)
			if err != nil {
				return err
			}
			for _, v := range areas {
				v.AreaNamePath = strings.Replace(v.AreaNamePath, areaPo.AreaNamePath, newAreaNamePath, 1)
				err := aiDB.Update(l.ctx, v)
				if err != nil {
					return err
				}
			}
			areaPo.AreaNamePath = newAreaNamePath
		}

		l.setPoByPb(areaPo, in)

		err = relationDB.NewAreaInfoRepo(tx).Update(l.ctx, areaPo)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查出错")
		}
		return nil
	})
	l.svcCtx.AreaCache.SetData(l.ctx, in.AreaID, nil)
	return &sys.Empty{}, err
}
func (l *AreaInfoUpdateLogic) setPoByPb(po *relationDB.SysAreaInfo, pb *sys.AreaInfo) {
	//不支持更改 区域所属项目，因此不进行赋值

	//支持区域 改为 第一级区域（改字段前端必填）
	//po.ParentAreaID = pb.ParentAreaID
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if uc.IsAdmin {
		if pb.DeviceCount != nil {
			po.DeviceCount = pb.DeviceCount.Value
		}
		if pb.GroupCount != nil {
			po.GroupCount = pb.GroupCount.Value
		}
	}
	if pb.Tags != nil {
		po.Tags = pb.Tags
	}
	if pb.AreaName != "" {
		po.AreaName = pb.AreaName
	}
	if pb.Position != nil {
		po.Position = logic.ToStorePoint(pb.Position)
	}
	if pb.Desc != nil {
		po.Desc = pb.Desc.GetValue()
	}
	if pb.UseBy != "" {
		po.UseBy = pb.UseBy
	}
	if pb.IsUpdateAreaImg {
		if po.AreaImg != "" && !oss.IsCommonFile(l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneHeadIng, po.AreaImg) {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, po.AreaImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", po.AreaImg, err)
			}
		}
		if pb.AreaImg != "" {
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneHeadIng, fmt.Sprintf("%d/%s", pb.AreaID, oss.GetFileNameWithPath(pb.AreaImg)))
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(pb.AreaImg, nwePath)
			if err != nil {
				l.Error(err)
			} else {
				po.AreaImg = path
			}
		} else {
			po.AreaImg = ""
		}

	}
	if pb.IsUpdateConfigFile {
		if po.ConfigFile != "" && !oss.IsCommonFile(l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneConfigFile, po.ConfigFile) {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, po.ConfigFile, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", po.ConfigFile, err)
			}
		}
		if pb.ConfigFile != "" {
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneConfigFile, fmt.Sprintf("%d/%s", pb.AreaID, oss.GetFileNameWithPath(pb.ConfigFile)))
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(pb.ConfigFile, nwePath)
			if err != nil {
				l.Error(err)
			} else {
				po.ConfigFile = path
			}
		} else {
			po.ConfigFile = ""
		}

	}
}
