package usermanagelogic

import (
	"context"
	projectmanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/projectmanage"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
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
	if err != nil {
		return nil, err
	}
	tc, err := relationDB.NewTenantConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantConfigFilter{TenantCode: ctxs.GetUserCtx(l.ctx).TenantCode})
	if err != nil {
		return nil, err
	}
	pis, err := relationDB.NewProjectInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProjectInfoFilter{AdminUserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	if tc.CheckUserDelete == 1 {

		if len(pis) > 0 {
			return nil, errors.Permissions.AddMsg("名下还有项目,需要先转让项目给其他人才可以注销")
		}
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
		for _, v := range pis {
			if tc.CheckUserDelete != 1 { //如果是不检查项目下的设备,那么就直接全部删除
				err = projectmanagelogic.ProjectDelete(l.ctx, tx, int64(v.ProjectID))
				if err != nil {
					return err
				}
			}
		}

		return err
	})
	l.Infof("%s.delete uid=%v", utils.FuncName(), in.UserID)
	err = l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.CoreUserDelete, application.IDs{IDs: []int64{in.UserID}})
	if err != nil {
		l.Errorf("Publish userDelete %v err:%v", in, err)
	}

	return &sys.Empty{}, nil
}
