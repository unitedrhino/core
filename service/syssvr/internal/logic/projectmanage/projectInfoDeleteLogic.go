package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoDeleteLogic {
	return &ProjectInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 删除项目
func (l *ProjectInfoDeleteLogic) ProjectInfoDelete(in *sys.ProjectWithID) (*sys.Empty, error) {
	if in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("项目ID参数必填")
	}
	ctxs.GetUserCtx(l.ctx).AllProject = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllProject = false
	}()
	po, err := checkProject(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查项目出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if uc != nil && (!uc.IsAdmin && uc.UserID != po.AdminUserID) { //只有管理员或项目负责人可以删除
		return nil, errors.Permissions.WithMsg("只有管理员可以删除")
	}

	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{})
	if err != nil {
		return nil, err
	}
	if ti.DefaultProjectID == in.ProjectID || po.IsSysCreated == def.True {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("默认项目禁止删除")
	}

	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		err = ProjectDelete(l.ctx, tx, in.ProjectID)
		return err
	})
	if err == nil {
		err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.CoreProjectInfoDelete, in.ProjectID)
		if err != nil {
			l.Error(err)
		}
	}
	return &sys.Empty{}, err
}

func ProjectDelete(ctx context.Context, tx *gorm.DB, id int64) error {
	err := relationDB.NewProjectInfoRepo(tx).Delete(ctx, id)
	if err != nil {
		return err
	}
	err = relationDB.NewProjectProfileRepo(tx).DeleteByFilter(ctx, relationDB.ProjectProfileFilter{ProjectID: id})
	if err != nil {
		return err
	}
	err = relationDB.NewAreaInfoRepo(tx).DeleteByFilter(ctx, relationDB.AreaInfoFilter{ProjectID: id})
	if err != nil {
		return errors.Fmt(err).WithMsg("删除区域及子区域出错")
	}
	err = relationDB.NewDataProjectRepo(tx).DeleteByFilter(ctx, relationDB.DataProjectFilter{ProjectID: id})
	if err != nil {
		return err
	}
	err = relationDB.NewDataAreaRepo(tx).DeleteByFilter(ctx, relationDB.DataAreaFilter{ProjectID: id})
	if err != nil {
		return err
	}
	err = relationDB.NewUserAreaApplyRepo(tx).DeleteByFilter(ctx, relationDB.UserAreaApplyFilter{ProjectID: id})
	if err != nil {
		return err
	}
	return nil
}
