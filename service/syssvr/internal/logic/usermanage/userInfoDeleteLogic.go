package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoDeleteLogic {
	return &UserInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoDeleteLogic) UserInfoDelete(in *sys.UserInfoDeleteReq) (*sys.Empty, error) {
	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{Code: ctxs.GetUserCtx(l.ctx).TenantCode})
	if err != nil {
		return nil, err
	}
	if ti.AdminUserID == in.UserID {
		return nil, errors.Permissions.AddMsg("超级管理员不允许删除")
	}
	pis, err := relationDB.NewProjectInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProjectInfoFilter{AdminUserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		err := uidb.Delete(l.ctx, cast.ToInt64(in.UserID))
		if err != nil {
			return err
		}
		err = relationDB.NewUserProfileRepo(tx).DeleteByFilter(l.ctx, relationDB.UserProfileFilter{UserID: in.UserID})
		if err != nil {
			return err
		}
		err = relationDB.NewUserRoleRepo(tx).DeleteByFilter(l.ctx, relationDB.UserRoleFilter{UserID: in.UserID})
		if err != nil {
			return err
		}
		cfg, err := relationDB.NewTenantConfigRepo(tx).FindOne(l.ctx)
		if err != nil {
			return err
		}
		if cfg.RegisterCreateProject != def.True { //如果是自动创建的,那还需要清除项目信息
			return nil
		}
		err = relationDB.NewProjectInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.ProjectInfoFilter{AdminUserID: in.UserID})
		return err
	})
	l.Infof("%s.delete uid=%v", utils.FuncName(), in.UserID)
	err = l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.CoreUserDelete, application.IDs{IDs: []int64{in.UserID}})
	if err != nil {
		l.Errorf("Publish userDelete %v err:%v", in, err)
	}
	if len(pis) != 0 {
		var ids []int64
		for _, v := range pis {
			ids = append(ids, int64(v.ProjectID))
		}
		err = l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.CoreProjectDelete, application.IDs{IDs: ids})
		if err != nil {
			l.Errorf("Publish projectDelete %v err:%v", in, err)
		}
	}

	return &sys.Empty{}, nil
}
