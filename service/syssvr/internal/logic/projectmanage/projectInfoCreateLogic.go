package projectmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ProjectInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoCreateLogic {
	return &ProjectInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 新增项目
func (l *ProjectInfoCreateLogic) ProjectInfoCreate(in *sys.ProjectInfo) (*sys.ProjectWithID, error) {
	if in.ProjectName == "" {
		return nil, errors.Parameter
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllProject = true
	defer func() {
		uc.AllProject = false
	}()
	if in.AdminUserID == 0 {
		in.AdminUserID = uc.UserID
	}
	po := &relationDB.SysProjectInfo{
		ProjectID:   dataType.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName: in.ProjectName,
		//CompanyName: utils.ToEmptyString(in.CompanyName),
		AdminUserID:  in.AdminUserID,
		Ppsm:         in.Ppsm,
		Area:         in.Area.GetValue(),
		IsSysCreated: in.IsSysCreated,
		//Region:      utils.ToEmptyString(in.Region),
		Address:  utils.ToEmptyString(in.Address),
		Position: logic.ToStorePoint(in.Position),
		Desc:     utils.ToEmptyString(in.Desc),
	}
	po.Tags = in.Tags
	if po.Tags == nil {
		po.Tags = map[string]string{}
	}
	_, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, in.AdminUserID)
	if err != nil {
		return nil, err
	}
	if in.IsUpdateProjectImg && in.ProjectImg != "" {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProject, oss.SceneHeadIng,
			fmt.Sprintf("%d/%s", po.ProjectID, oss.GetFileNameWithPath(in.ProjectImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.ProjectImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.ProjectImg = path
	}
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		//tiDb := relationDB.NewTenantInfoRepo(tx)
		//ti, err := tiDb.FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{})
		//if err != nil {
		//	return err
		//}
		piDb := relationDB.NewProjectInfoRepo(tx)
		//total, err := piDb.CountByFilter(l.ctx, relationDB.ProjectInfoFilter{})
		//if err != nil {
		//	return err
		//}
		//if total >= ti.ProjectLimit {
		//	return errors.OutRange.WithMsgf("最多创建%v个项目", ti.ProjectLimit)
		//}
		err = piDb.Insert(l.ctx, po)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
		err = relationDB.NewDataProjectRepo(tx).Insert(l.ctx, &relationDB.SysDataProject{
			ProjectID:  int64(po.ProjectID),
			TargetType: def.TargetUser,
			TargetID:   po.AdminUserID,
			AuthType:   def.AuthAdmin,
		})
		return err
		//ti.ProjectLimit++
		//err = tiDb.Update(l.ctx, ti)
		//if err != nil {
		//	return err
		//}
	})
	if err != nil {
		l.Errorf("%s err=%+v", utils.FuncName(), err)
		return nil, err
	}
	cache.ClearProjectAuth(uc.UserID)
	return &sys.ProjectWithID{ProjectID: int64(po.ProjectID)}, nil
}
