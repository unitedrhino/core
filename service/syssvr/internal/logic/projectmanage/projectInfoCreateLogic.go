package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
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
	ctxs.GetUserCtx(l.ctx).AllProject = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllProject = false
	}()

	po := &relationDB.SysProjectInfo{
		ProjectID:   stores.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName: in.ProjectName,
		//CompanyName: utils.ToEmptyString(in.CompanyName),
		AdminUserID: in.AdminUserID,
		//Region:      utils.ToEmptyString(in.Region),
		//Address:     utils.ToEmptyString(in.Address),
		Position: logic.ToStorePoint(in.Position),
		Desc:     utils.ToEmptyString(in.Desc),
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		tiDb := relationDB.NewTenantInfoRepo(tx)
		ti, err := tiDb.FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{})
		if err != nil {
			return err
		}
		piDb := relationDB.NewProjectInfoRepo(tx)
		total, err := piDb.CountByFilter(l.ctx, relationDB.ProjectInfoFilter{})
		if err != nil {
			return err
		}
		if total >= ti.ProjectLimit {
			return errors.OutRange.WithMsgf("最多创建%v个项目", ti.ProjectLimit)
		}
		err = piDb.Insert(l.ctx, po)
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return err
		}
		ti.ProjectLimit++
		err = tiDb.Update(l.ctx, ti)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorf("%s err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &sys.ProjectWithID{ProjectID: int64(po.ProjectID)}, nil
}
