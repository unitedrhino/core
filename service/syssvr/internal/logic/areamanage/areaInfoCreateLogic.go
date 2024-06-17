package areamanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AreaInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoCreateLogic {
	return &AreaInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增区域
func (l *AreaInfoCreateLogic) AreaInfoCreate(in *sys.AreaInfo) (*sys.AreaWithID, error) {
	if in.AreaName == "" || ////root节点不为0
		in.ParentAreaID == def.NotClassified { //未分类不能有下属的区域
		return nil, errors.Parameter
	}
	uc := ctxs.GetUserCtx(l.ctx)
	list := l.svcCtx.Slot.Get(l.ctx, "areaInfo", "create")
	err := list.Request(l.ctx, in, nil)
	if err != nil {
		return nil, err
	}
	if in.ProjectID == 0 {
		in.ProjectID = uc.ProjectID
	}
	if !uc.IsAdmin {
		if uc.ProjectAuth == nil || uc.ProjectAuth[in.ProjectID] == nil || uc.ProjectAuth[in.ProjectID].AuthType != def.AuthAdmin {
			return nil, errors.Permissions.AddMsg("只有项目管理员才能创建区域")
		}
	}
	if in.ParentAreaID == 0 {
		in.ParentAreaID = def.RootNode
	}

	var areaID = l.svcCtx.AreaID.GetSnowflakeId()
	var areaIDPath string = cast.ToString(areaID) + "-"
	var areaNamePath = in.AreaName + "-"
	areaPo := &relationDB.SysAreaInfo{
		AreaID:       stores.AreaID(areaID),
		ParentAreaID: in.ParentAreaID,                //创建时必填
		ProjectID:    stores.ProjectID(in.ProjectID), //创建时必填
		AreaIDPath:   areaIDPath,
		AreaNamePath: areaNamePath,
		AreaName:     in.AreaName,
		Position:     logic.ToStorePoint(in.Position),
		Desc:         utils.ToEmptyString(in.Desc),
		IsLeaf:       def.True,
		UseBy:        in.UseBy,
	}
	if in.IsUpdateAreaImg && in.AreaImg != "" {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BUsinessArea, oss.SceneHeadIng, fmt.Sprintf("%d/%s", areaID, oss.GetFileNameWithPath(in.AreaImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.AreaImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		areaPo.AreaImg = path
	}
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		projPo, err := checkProject(l.ctx, tx, in.ProjectID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查项目出错")
		} else if projPo == nil {
			return errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
		}
		aiRepo := relationDB.NewAreaInfoRepo(tx)
		if in.ParentAreaID != def.RootNode { //有选了父级项目区域
			pa, err := checkParentArea(l.ctx, tx, in.ParentAreaID)
			if err != nil {
				return err
			}
			areaPo.AreaIDPath = pa.AreaIDPath + cast.ToString(areaID) + "-"
			areaPo.AreaNamePath = pa.AreaNamePath + in.AreaName + "-"
			pa.LowerLevelCount++
			err = addSubAreaIDs(l.ctx, tx, pa, int64(areaPo.AreaID))
			if err != nil {
				return err
			}
			pa.IsLeaf = def.False
			err = aiRepo.Update(l.ctx, pa)
			if err != nil {
				return err
			}
		}
		err = aiRepo.Insert(l.ctx, areaPo)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
		return nil
	})
	if err == nil {
		FillProjectAreaCount(l.ctx, l.svcCtx, int64(areaPo.ProjectID))
	}
	cache.ClearProjectAuth(uc.UserID)

	return &sys.AreaWithID{AreaID: int64(areaPo.AreaID)}, err
}

func FillProjectAreaCount(ctx context.Context, svcCtx *svc.ServiceContext, projectID int64) {
	ctxs.GoNewCtx(ctx, func(ctx context.Context) {
		count, err := relationDB.NewAreaInfoRepo(ctx).CountByFilter(ctx, relationDB.AreaInfoFilter{ProjectID: projectID})
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
		err = relationDB.NewProjectInfoRepo(ctx).Update(ctx,
			&relationDB.SysProjectInfo{ProjectID: stores.ProjectID(projectID), AreaCount: count},
			"areaCount")
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
		err = svcCtx.ProjectCache.SetData(ctx, projectID, nil)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	})

}
