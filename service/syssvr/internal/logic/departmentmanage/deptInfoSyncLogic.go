package departmentmanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	usermanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptInfoSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	cli svc.Clients
}

func NewDeptInfoSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptInfoSyncLogic {
	return &DeptInfoSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptInfoSyncLogic) DeptInfoSync(in *sys.DeptInfoSyncReq) (*sys.DeptInfoSyncResp, error) {
	cli, err := l.svcCtx.Cm.GetClients(l.ctx, in.AppCode)
	if err != nil || cli.DingMini == nil {
		return nil, errors.System.AddDetail(err)
	}
	l.cli = cli
	l.DeptInfoSyncDingTalk(&relationDB.SysDeptInfo{
		ID:         def.RootNode,
		Name:       "根节点",
		Status:     def.True,
		DingTalkID: def.RootNode,
	}, in)
	return &sys.DeptInfoSyncResp{}, nil
}

func (l *DeptInfoSyncLogic) DeptInfoSyncDingTalkUser(info *relationDB.SysDeptInfo, in *sys.DeptInfoSyncReq) error {
	if in.UserMode == 0 {
		return nil
	}
	c := 0
	page := 100
	for {
		req := request.NewDeptDetailUserInfo(int(info.DingTalkID), c, page)
		dings, err := l.cli.DingMini.GetDeptDetailUserInfo(req.Build())
		if err != nil {
			return errors.System.AddDetail(err)
		}
		old, err := relationDB.NewDeptUserRepo(l.ctx).FindByFilter(l.ctx,
			relationDB.DeptUserFilter{DeptID: info.ID, WithUser: true}, nil)
		if err != nil {
			return err
		}
		var (
			deptPhoneMap = map[string]*relationDB.SysUserInfo{}
			deptEmailMap = map[string]*relationDB.SysUserInfo{}
			deptIDMap    = map[string]*relationDB.SysUserInfo{}
		)
		for _, o := range old {
			if o.User == nil {
				continue
			}
			u := o.User
			if u.Phone.Valid {
				deptPhoneMap[u.Phone.String] = u
			}
			if u.Email.Valid {
				deptEmailMap[u.Email.String] = u
			}

			if u.DingTalkUserID.Valid {
				deptIDMap[u.DingTalkUserID.String] = u
			}
		}
		for _, ding := range dings.Page.List {
			po := deptIDMap[ding.UserId]
			if po == nil {
				po = deptPhoneMap[ding.Telephone]
			}
			if po == nil {
				po = deptEmailMap[ding.Email]
			}
			delete(deptIDMap, ding.UserId)
			delete(deptPhoneMap, ding.Telephone)
			delete(deptEmailMap, ding.Email)
			if po == nil {
				uc, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ding.UserId})
				if err != nil {
					if !errors.Cmp(err, errors.NotFind) {
						return err
					}
					uc, err = relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{ding.Email, ding.Telephone}})
					if !errors.Cmp(err, errors.NotFind) {
						return err
					}
				}
				if uc == nil {
					userID := l.svcCtx.UserID.GetSnowflakeId()
					uc = &relationDB.SysUserInfo{
						UserID:         userID,
						DingTalkUserID: sql.NullString{Valid: true, String: ding.UserId},
						NickName:       ding.Name,
					}
					if ding.OrgEmail != "" {
						uc.Email = sql.NullString{String: ding.OrgEmail, Valid: true}
						uc.UserName = uc.Email
					}
					if ding.Mobile != "" {
						uc.Phone = sql.NullString{String: ding.Mobile, Valid: true}
						uc.UserName = uc.Phone
					}
					if ding.Extension != "" {
						var tags = map[string]string{}
						err = json.Unmarshal([]byte(ding.Extension), &tags)
						if err == nil {
							uc.Tags = tags
						}
					}
					err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
						return usermanagelogic.Register(l.ctx, l.svcCtx, uc, tx)
					})
				}
				err = relationDB.NewDeptUserRepo(l.ctx).Insert(l.ctx, &relationDB.SysDeptUser{
					UserID: uc.UserID,
					DeptID: info.ID,
				})
				if err != nil {
					return err
				}
				continue
			}
			if in.UserMode == dept.UserModeAdd {
				continue
			}
			if po.NickName != ding.Name || po.DingTalkUserID.String != ding.UserId || po.Phone.String != ding.Telephone || po.Email.String != ding.Email {
				delete(deptPhoneMap, po.Phone.String)
				delete(deptEmailMap, po.Email.String)
				delete(deptIDMap, po.DingTalkUserID.String)
				po.NickName = ding.Name
				po.DingTalkUserID = sql.NullString{String: ding.UserId, Valid: true}
				if ding.OrgEmail != "" {
					po.Email = sql.NullString{String: ding.OrgEmail, Valid: true}
				}
				if ding.Extension != "" {
					var tags = map[string]string{}
					err = json.Unmarshal([]byte(ding.Extension), &tags)
					if err == nil {
						po.Tags = tags
					}
				}
				if ding.Mobile != "" {
					po.Phone = sql.NullString{String: ding.Mobile, Valid: true}
				}
				err = relationDB.NewUserInfoRepo(l.ctx).Update(l.ctx, po)
				if err != nil {
					return err
				}
			}
			if len(deptIDMap) > 0 && in.UserMode == dept.UserModeDelete { //如果存在删除的情况
				for _, one := range deptIDMap {
					err := relationDB.NewDeptUserRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DeptUserFilter{UserID: one.UserID})
					if err != nil {
						return err
					}
				}
			}
		}

		if !dings.Page.HasMore {
			break
		}
		c = dings.Page.NextCursor
	}
	return nil
}

func (l *DeptInfoSyncLogic) DeptInfoSyncDingTalk(info *relationDB.SysDeptInfo, in *sys.DeptInfoSyncReq) (*sys.DeptInfoSyncResp, error) {
	if err := l.DeptInfoSyncDingTalkUser(info, in); err != nil {
		return nil, err
	}
	req := request.NewDeptList()
	req.SetDeptId(int(info.DingTalkID))
	dings, err := l.cli.DingMini.GetDeptList(req.Build())
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	old, err := relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptInfoFilter{ParentID: info.ID}, nil)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	var (
		deptNameMap = map[string]*relationDB.SysDeptInfo{}
		deptIDMap   = map[int64]*relationDB.SysDeptInfo{}
	)
	for _, o := range old {
		deptNameMap[o.Name] = o
		if o.DingTalkID > def.RootNode {
			deptIDMap[o.DingTalkID] = o
		}
	}
	for _, ding := range dings.List {
		po := deptIDMap[int64(ding.Id)]
		if po == nil {
			po = deptNameMap[ding.Name]
		}
		delete(deptIDMap, int64(ding.Id))
		delete(deptNameMap, ding.Name)
		if po == nil {
			newOne := &relationDB.SysDeptInfo{
				ParentID:   info.ID,
				Name:       ding.Name,
				Status:     def.True,
				IDPath:     info.IDPath,
				DingTalkID: int64(ding.Id),
			}
			err := relationDB.NewDeptInfoRepo(l.ctx).Insert(l.ctx, newOne)
			if err != nil {
				return nil, err
			}
			newOne.IDPath = info.IDPath + fmt.Sprintf("%d-", newOne.ID)
			err = relationDB.NewDeptInfoRepo(l.ctx).Update(l.ctx, newOne)
			if err != nil {
				return nil, err
			}
			continue
		}
		if in.DeptMode == dept.DeptModeAdd {
			continue
		}
		if po.Name != ding.Name || po.DingTalkID != int64(ding.Id) {
			delete(deptNameMap, po.Name)
			po.Name = ding.Name
			po.DingTalkID = int64(ding.Id)
			err = relationDB.NewDeptInfoRepo(l.ctx).Update(l.ctx, po)
			if err != nil {
				return nil, err
			}

		}
	}
	if len(deptIDMap) > 0 && in.DeptMode == dept.DeptModeDelete { //如果存在删除的情况
		for _, one := range deptIDMap {
			err := relationDB.NewDeptInfoRepo(l.ctx).Delete(l.ctx, one.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	old, err = relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptInfoFilter{ParentID: info.ID}, nil)
	if err != nil {
		return nil, err
	}
	for _, o := range old {
		_, err := l.DeptInfoSyncDingTalk(o, in)
		if err != nil {
			return nil, err
		}
	}
	return &sys.DeptInfoSyncResp{}, nil
}
