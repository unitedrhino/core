package tenantmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppUpdateLogic {
	return &TenantAppUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppUpdateLogic) TenantAppUpdate(in *sys.TenantAppInfo) (*sys.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.Code != "" && uc.TenantCode != def.TenantCodeDefault {
		return nil, errors.Permissions
	}
	old, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
	if err != nil {
		return nil, err
	}
	if in.WxMini != nil {
		old.WxMini = utils.Copy[relationDB.SysTenantThird](in.WxMini)
	}
	if in.DingMini != nil {
		old.DingMini = utils.Copy[relationDB.SysTenantThird](in.DingMini)
	}
	if in.WxOpen != nil {
		old.WxOpen = utils.Copy[relationDB.SysTenantThird](in.WxOpen)
	}
	if in.Android != nil {
		if in.Android.IsUpdateFilePath {
			if old.Android != nil && old.Android.FilePath != "" {
				err := l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, old.Android.FilePath, common.OptionKv{})
				if err != nil {
					l.Errorf("Delete file err path:%v,err:%v", old.Android.FilePath, err)
				}
			}
			if in.Android.FilePath != "" {
				nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessApp, oss.SceneFirmware, fmt.Sprintf("%s/%s", old.AppCode, oss.GetFileNameWithPath(in.Android.FilePath)))
				path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(in.Android.FilePath, nwePath)
				if err != nil {
					l.Error(err)
				} else {
					in.Android.FilePath = path
				}
			} else {
				in.Android.FilePath = ""
			}
		}
		old.Android = utils.Copy[relationDB.SysThirdApp](in.Android)
	}
	if in.LoginTypes != nil {
		old.LoginTypes = in.LoginTypes
	}
	err = relationDB.NewTenantAppRepo(l.ctx).Update(ctxs.WithRoot(l.ctx), old)
	if err == nil {
		ctx := l.ctx
		if in.Code != "" {
			ctx = ctxs.BindTenantCode(l.ctx, in.Code, def.RootNode)
		}
		err := l.svcCtx.Cm.ClearClients(ctx, in.AppCode)
		if err != nil {
			l.Error(err)
		}
	}
	return &sys.Empty{}, err
}
