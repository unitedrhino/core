package usermanagelogic

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"

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

func (l *UserInfoUpdateLogic) UserInfoUpdate(in *sys.UserInfoUpdateReq) (*sys.Empty, error) {
	info := in.Info
	var (
		updateStatus bool
	)

	ui, err := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{UserIDs: []int64{info.UserID}, WithRoles: true})
	if err != nil {
		l.Errorf("%s.FindOne OperUserID=%d err=%v", utils.FuncName(), info.UserID, err)
		return nil, err
	}
	if in.WithRoot {
		if info.Phone != nil {
			ui.Phone = utils.AnyToNullString(info.Phone)
			ui.UserName = ui.Phone
		}
		if info.Email != nil {
			ui.Email = utils.AnyToNullString(info.Email)
			ui.UserName = ui.Email
		}
		if info.Status != 0 && ui.Status != info.Status {
			ui.Status = info.Status
			updateStatus = true
		}
		if info.UserName != "" {
			ui.UserName = sql.NullString{String: info.UserName, Valid: true}
		}
		if info.Tags != nil {
			ui.Tags = info.Tags
		}
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
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessUserManage, oss.SceneHeadIng, fmt.Sprintf("%d/%s", ui.UserID, oss.GetFileNameWithPath(info.HeadImg)))
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
	if updateStatus == true && ui.Status == def.False {
		err = relationDB.NewDataAreaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataAreaFilter{
			TargetID:   ui.UserID,
			TargetType: def.TargetUser,
		})
		if err != nil {
			l.Error(err)
		}
	}
	l.svcCtx.UserCache.SetData(l.ctx, ui.UserID, nil)
	err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.CoreUserUpdate, def.IDs{IDs: []int64{ui.UserID}})
	if err != nil {
		l.Errorf("Publish CoreUserUpdate %v err:%v", ui, err)
	}
	return &sys.Empty{}, nil
}
