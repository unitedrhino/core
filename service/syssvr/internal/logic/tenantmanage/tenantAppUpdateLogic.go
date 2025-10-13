package tenantmanagelogic

import (
	"context"
	"fmt"

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
	"gorm.io/gorm"

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
	if in.Code == "" {
		in.Code = uc.TenantCode
	}
	old, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
	if err != nil {
		return nil, err
	}
	//app, err := relationDB.NewAppInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AppInfoFilter{Codes: []string{in.AppCode}})
	//if err != nil {
	//	return nil, err
	//}
	var thirdFields = map[string]any{}
	if in.WxMini != nil {
		old.WxMini = utils.Copy[relationDB.SysTenantThird](in.WxMini)
		thirdFields["wx_mini_app_id"] = in.WxMini.AppID
	}
	if in.DingMini != nil {
		old.DingMini = utils.Copy[relationDB.SysTenantThird](in.DingMini)
		thirdFields["ding_mini_app_id"] = in.WxMini.AppID
	}
	if in.WxOpen != nil {
		old.WxOpen = utils.Copy[relationDB.SysTenantThird](in.WxOpen)
		thirdFields["wx_open_app_id"] = in.WxMini.AppID
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
		thirdFields["android_version"] = in.Android.Version
		thirdFields["android_file_path"] = in.Android.FilePath
		thirdFields["android_version_desc"] = in.Android.VersionDesc
	}

	if in.LoginTypes != nil {
		old.LoginTypes = in.LoginTypes
	}
	if in.IsAutoRegister != 0 {
		old.IsAutoRegister = in.IsAutoRegister
	}
	if in.Config != "" {
		old.Config = in.Config
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		db := relationDB.NewTenantAppRepo(tx)
		err = db.Update(ctxs.WithRoot(l.ctx), old)
		if err != nil {
			return err
		}
		err = db.UpdateWithField(ctxs.WithRoot(l.ctx), relationDB.TenantAppFilter{AppCodes: []string{in.AppCode}}, thirdFields)
		return err
	})
	if err == nil {
		ctx := l.ctx
		if in.Code != "" {
			ctx = ctxs.BindTenantCode(l.ctx, in.Code, def.RootNode)
		}
		po, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
		if err != nil {
			l.Errorf("NewTenantAppRepo.FindOneByFilter err:%v", err)
		} else {
			err = l.svcCtx.ThirdClientsManage.SetOne(ctx, po)
			if err != nil {
				l.Errorf("ThirdClientsManage.SetOne err:%v", err)
			}
		}
	}
	return &sys.Empty{}, err
}
