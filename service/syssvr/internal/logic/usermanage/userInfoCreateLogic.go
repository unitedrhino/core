package usermanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoCreateLogic {
	return &UserInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoCreateLogic) UserInfoInsert(in *sys.UserInfoCreateReq) (int64, error) {
	info := in.Info
	var userID int64
	//首先校验账号格式使用正则表达式，对用户账号做格式校验：只能是大小写字母，数字和下划线，减号
	if info.UserName != "" {
		err := CheckUserName(info.UserName)
		if err != nil && !utils.SliceIn(info.UserName, info.Email.GetValue(), info.Phone.GetValue()) {
			return 0, err
		}
	}
	if info.Tags == nil {
		info.Tags = map[string]string{}
	}
	if info.Phone.GetValue() != "" {
		if !utils.IsPhone(info.Phone.GetValue()) {
			return 0, errors.Parameter.AddMsgf("手机号格式错误")
		}
		if info.UserName == "" {
			info.UserName = info.Phone.GetValue()
		}
	}
	if info.Email.GetValue() != "" {
		if !utils.IsEmail(info.Email.GetValue()) {
			return 0, errors.Parameter.AddMsgf("邮箱格式错误")
		}
		if info.UserName == "" {
			info.UserName = info.Email.GetValue()
		}
	}
	if info.UserName == "" {
		return 0, errors.Parameter.AddMsgf("用户名,手机号和邮箱至少要填一个")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if uc == nil {
		return 0, errors.Permissions.WithMsg("无租户号")
	}
	if len(in.RoleIDs) == 0 { //填充默认角色
		t, err := relationDB.NewTenantConfigRepo(l.ctx).FindOne(l.ctx)
		if err != nil {
			return 0, err
		}
		in.RoleIDs = []int64{t.RegisterRoleID}
	}
	//校验密码强度
	err := CheckPwd(l.svcCtx, info.Password)
	if err != nil {
		return 0, err
	}
	if info.Role == 0 {
		info.Role = in.RoleIDs[0]
	} else if !utils.SliceIn(info.Role, in.RoleIDs...) {
		return 0, errors.Parameter.AddMsgf("用户默认角色不存在")
	}
	count, err := relationDB.NewRoleInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.RoleInfoFilter{IDs: in.RoleIDs})
	if err != nil {
		return 0, err
	}
	if int(count) != len(in.RoleIDs) {
		return 0, errors.Parameter.AddMsgf("角色有不存在的")
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		var account = []string{info.UserName}
		if info.Phone.GetValue() != "" {
			account = append(account, info.Phone.GetValue())
		}
		if info.Email.GetValue() != "" {
			account = append(account, info.Email.GetValue())
		}
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: account})
		if err == nil { //已注册
			//提示重复注册
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return err
		}
		//1.生成uid
		userID = l.svcCtx.UserID.GetSnowflakeId()

		//2.对密码进行md5加密
		password := utils.MakePwd(info.Password, userID, false)
		ui := relationDB.SysUserInfo{
			UserID:    userID,
			UserName:  sql.NullString{String: info.UserName, Valid: true},
			Password:  password,
			NickName:  info.NickName,
			Role:      info.Role,
			IsAllData: info.IsAllData,
		}
		if info.Email.GetValue() != "" {
			ui.Email = sql.NullString{String: info.Email.GetValue(), Valid: true}
		}
		if info.Phone.GetValue() != "" {
			ui.Phone = sql.NullString{String: info.Phone.GetValue(), Valid: true}
		}
		err = uidb.Insert(l.ctx, &ui)
		if err != nil { //并发情况下有可能重复所以需要再次判断一次
			if errors.Cmp(err, errors.NotFind) {
				return errors.DuplicateUsername.AddDetail(info.UserName)
			}
			l.Errorf("%s.Inserts err=%#v", utils.FuncName(), err)
			return err
		}
		err := relationDB.NewUserRoleRepo(tx).MultiUpdate(l.ctx, ui.UserID, in.RoleIDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
func (l *UserInfoCreateLogic) UserInfoCreate(in *sys.UserInfoCreateReq) (*sys.UserCreateResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	userID, err := l.UserInfoInsert(in)
	if err != nil {
		return nil, err
	}
	e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{userID}})
	if e != nil {
		l.Errorf("Publish CoreUserCreate %v err:%v", userID, e)
	}
	return &sys.UserCreateResp{UserID: userID}, nil
}
