package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
func (l *ProjectInfoDeleteLogic) ProjectInfoDelete(in *sys.ProjectWithID) (*sys.Response, error) {
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

	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewProjectInfoRepo(tx).Delete(l.ctx, in.ProjectID)
		if err != nil {
			return err
		}
		err = relationDB.NewAreaInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.AreaInfoFilter{ProjectID: in.ProjectID})
		if err != nil {
			return errors.Fmt(err).WithMsg("删除区域及子区域出错")
		}
		err = relationDB.NewDataAreaRepo(tx).DeleteByFilter(l.ctx, relationDB.DataAreaFilter{ProjectID: in.ProjectID})
		if err != nil {
			return err
		}
		err = relationDB.NewUserAreaApplyRepo(tx).DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{ProjectID: in.ProjectID})
		if err != nil {
			return err
		}
		return nil
	})

	return &sys.Response{}, nil
}
