package usermanagelogic

import (
	"context"
	projectmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/projectmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	uc := ctxs.GetUserCtx(l.ctx)
	ti, err := l.svcCtx.TenantCache.GetData(l.ctx, uc.TenantCode)
	if err != nil {
		return nil, err
	}
	if ti.AdminUserID == in.UserID {
		return nil, errors.Permissions.AddMsg("超级管理员不允许删除")
	}

	tc, err := l.svcCtx.TenantConfigCache.GetData(l.ctx, uc.TenantCode)
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
		err := uidb.Delete(ctxs.WithRoot(l.ctx), cast.ToInt64(in.UserID))
		if err != nil {
			return err
		}
		err = relationDB.NewUserProfileRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx), relationDB.UserProfileFilter{UserID: in.UserID})
		if err != nil {
			return err
		}
		err = relationDB.NewUserRoleRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx), relationDB.UserRoleFilter{UserID: in.UserID})
		for _, v := range pis {
			if tc.CheckUserDelete != 1 { //如果是不检查项目下的设备,那么就直接全部删除
				err = projectmanagelogic.ProjectDelete(ctxs.WithRoot(l.ctx), tx, int64(v.ProjectID))
				if err != nil {
					return err
				}
				err = l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreProjectInfoDelete, v.ProjectID)
				if err != nil {
					l.Error(err)
				}
			}
		}

		return err
	})
	l.Infof("%s.delete uid=%v", utils.FuncName(), in.UserID)
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserDelete, def.IDs{IDs: []int64{in.UserID}})
	if err != nil {
		l.Errorf("Publish userDelete %v err:%v", in, err)
	}
	l.svcCtx.UserCache.SetData(l.ctx, in.UserID, nil)

	return &sys.Empty{}, nil
}
