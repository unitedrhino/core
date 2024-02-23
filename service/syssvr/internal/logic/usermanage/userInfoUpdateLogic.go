package usermanagelogic

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(in *sys.UserInfoUpdateReq) (*sys.Response, error) {
	info := in.Info
	ui, err := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{UserIDs: []int64{info.UserID}, WithRoles: true})
	if err != nil {
		l.Errorf("%s.FindOne UserID=%d err=%v", utils.FuncName(), info.UserID, err)
		return nil, err
	}
	if in.WithRoot {
		if info.Phone != nil {
			ui.Phone = utils.AnyToNullString(info.Phone)
		}
		if info.Email != nil {
			ui.Email = utils.AnyToNullString(info.Email)
		}
	}
	if info.UserName != "" {
		ui.UserName = sql.NullString{String: info.UserName, Valid: true}
	}
	if info.NickName != "" {
		ui.NickName = info.NickName
	}

	//性別有效才賦值，否則使用旧值
	if ui.Sex == def.Unknown {
		ui.Sex = def.Male
	} else {
		ui.Sex = info.Sex
	}

	//设置数据超管
	if info.IsAllData == 1 || info.IsAllData == 2 {
		ui.IsAllData = info.IsAllData
	}
	if info.Role != 0 && info.Role != ui.Role {
		ui.Role = info.Role
	}
	if info.IsUpdateHeadImg && info.HeadImg != "" {
		if ui.HeadImg != "" {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, ui.HeadImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", ui.HeadImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessUserManage, oss.SceneUserInfo, fmt.Sprintf("%d/%s", ui.UserID, oss.GetFileNameWithPath(info.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(info.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		ui.HeadImg = path
	}

	if info.Role != 0 { //默认角色只能修改为授权的角色
		for _, r := range ui.Roles {
			if r.RoleID == info.Role {
				ui.Role = info.Role
			}
		}
	}

	err = l.UiDB.Update(l.ctx, ui)
	if err != nil {
		l.Errorf("%s.Update ui=%v err=%v", utils.FuncName(), ui, err)
		return nil, err
	}
	l.Infof("%s.modified usersvr info = %+v", utils.FuncName(), ui)

	return &sys.Response{}, nil
}
