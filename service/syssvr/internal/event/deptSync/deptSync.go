package deptSync

import (
	"context"
	"database/sql"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	departmentmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/departmentmanage"
	usermanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"
)

type DeptSync struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewDeptSync(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSync {
	return &DeptSync{
		ctx:    ctx,
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func DingSync(ctx context.Context, svcCtx *svc.ServiceContext, jobID int64, df *payload.DataFrame) error {
	return NewDeptSync(ctx, svcCtx).HandleDing(jobID, df)
}
func init() {
	departmentmanagelogic.DingSync = DingSync
}

type userModifyOrg struct {
	TimeStamp string `json:"timeStamp"`
	DiffInfo  []struct {
		Prev struct {
			ExtFields string `json:"extFields"`
			Name      string `json:"name"`
			Telephone string `json:"telephone"`
			Remark    string `json:"remark"`
			WorkPlace string `json:"workPlace"`
			JobNumber string `json:"jobNumber"`
			Email     string `json:"email"`
		} `json:"prev"`
		Curr struct {
			ExtFields string `json:"extFields"`
			Name      string `json:"name"`
			Telephone string `json:"telephone"`
			Remark    string `json:"remark"`
			WorkPlace string `json:"workPlace"`
			JobNumber string `json:"jobNumber"`
			Email     string `json:"email"`
		} `json:"curr"`
		Userid string `json:"userid"`
	} `json:"diffInfo"`
	EventId    string   `json:"eventId"`
	OptStaffId string   `json:"optStaffId"`
	UserId     []string `json:"userId"`
	DeptId     []int64  `json:"deptId"`
}

func (d DeptSync) HandleDing(jobID int64, df *payload.DataFrame) error {
	po, err := relationDB.NewDeptSyncJobRepo(d.ctx).FindOne(d.ctx, jobID)
	if err != nil {
		return err
	}
	if po.SyncMode != dept.SyncModeRealTime {
		return nil
	}
	var data userModifyOrg
	err = json.Unmarshal([]byte(df.Data), &data)
	if err != nil {
		return err
	}
	eventType := df.GetHeader("eventType")
	switch eventType {
	case "user_leave_org": //用户离职
		if len(data.UserId) == 0 {
			return nil
		}
		err = relationDB.NewUserInfoRepo(d.ctx).UpdateWithField(d.ctx, relationDB.UserInfoFilter{DingTalkUserIDs: data.UserId},
			map[string]any{
				"status": def.False,
			})
		return err

	case "user_add_org": //用户添加
		cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
			AppID:     po.ThirdConfig.AppID,
			AppKey:    po.ThirdConfig.AppKey,
			AppSecret: po.ThirdConfig.AppSecret,
		})
		if err != nil {
			return err
		}
		for _, v := range data.UserId {
			ui, er := cli.GetUserDetail(&request.UserDetail{
				UserId: v,
			})
			if er != nil {
				return er
			}

			userID := d.svcCtx.UserID.GetSnowflakeId()
			uc := &relationDB.SysUserInfo{
				UserID:         userID,
				DingTalkUserID: sql.NullString{Valid: true, String: ui.UserId},
				NickName:       ui.Name,
			}
			if ui.UnionId != "" {
				uc.DingTalkUnionID = sql.NullString{Valid: true, String: ui.UnionId}
			}
			var accounts []string
			if ui.OrgEmail != "" {
				accounts = append(accounts, ui.OrgEmail)
			}
			if ui.Mobile != "" {
				accounts = append(accounts, ui.Mobile)
			}
			uc, err = relationDB.NewUserInfoRepo(d.ctx).FindOneByFilter(d.ctx, relationDB.UserInfoFilter{Accounts: accounts})
			if err == nil {
				if ui.OrgEmail != "" {
					uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
				}
				if ui.Mobile != "" {
					uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
				}
				uc.DingTalkUserID = sql.NullString{Valid: true, String: ui.UserId}
				if ui.UnionId != "" {
					uc.DingTalkUnionID = sql.NullString{Valid: true, String: ui.UnionId}
				}
				err = relationDB.NewUserInfoRepo(d.ctx).Update(d.ctx, uc)
				return err
			}
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			uc.NickName = ui.Name
			if len(ui.Extension) != 0 {
				var tags = map[string]string{}
				err = json.Unmarshal([]byte(ui.Extension), &tags)
				if err == nil {
					uc.Tags = tags
				}
			}
			err = stores.GetTenantConn(d.ctx).Transaction(func(tx *gorm.DB) error {
				return usermanagelogic.Register(d.ctx, d.svcCtx, uc, tx)
			})
			return err
		}

	case "user_modify_org": //用户信息变更
		for _, v := range data.DiffInfo {
			uis, err := relationDB.NewUserInfoRepo(d.ctx).FindByFilter(d.ctx, relationDB.UserInfoFilter{DingTalkUserID: v.Userid}, nil)
			if err != nil {
				return err
			}
			for _, ui := range uis {
				if v.Prev.ExtFields != v.Curr.ExtFields {
					var tags = map[string]string{}
					err = json.Unmarshal([]byte(v.Curr.ExtFields), &tags)
					if err == nil {
						ui.Tags = tags
					}
				}
				if v.Prev.Name != v.Curr.Name {
					ui.NickName = v.Curr.Name
				}
				if v.Prev.Email != v.Curr.Email {
					ui.Email = sql.NullString{Valid: true, String: v.Curr.Email}
				}
				err = relationDB.NewUserInfoRepo(d.ctx).Update(d.ctx, ui)
				if err != nil {
					d.Error(err)
					continue
				}
			}
		}
	case "org_dept_remove": //部门删除
		if len(data.DeptId) == 0 {
			return nil
		}
		dis, err := relationDB.NewDeptInfoRepo(d.ctx).FindByFilter(d.ctx, relationDB.DeptInfoFilter{DingTalkIDs: data.DeptId}, nil)
		if err != nil {
			return err
		}
		for _, v := range dis {
			_, err = departmentmanagelogic.NewDeptInfoDeleteLogic(d.ctx, d.svcCtx).DeptInfoDelete(&sys.WithID{Id: v.ID})
			if err != nil {
				d.Error(err)
				continue
			}
		}
	case "org_dept_create", "org_dept_modify":
		if len(data.DeptId) == 0 {
			return nil
		}
		cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
			AppID:     po.ThirdConfig.AppID,
			AppKey:    po.ThirdConfig.AppKey,
			AppSecret: po.ThirdConfig.AppSecret,
		})
		if err != nil {
			return err
		}
		err = departmentmanagelogic.SyncDeptDing(d.ctx, cli, &relationDB.SysDeptInfo{
			ID:         def.RootNode,
			Name:       "根节点",
			Status:     def.True,
			DingTalkID: def.RootNode,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d DeptSync) Timing() error {
	d.Infof("DeptSync Timing")
	pos, err := relationDB.NewDeptSyncJobRepo(d.ctx).FindByFilter(d.ctx, relationDB.DeptSyncJobFilter{
		Direction: dept.SyncDirectionFrom, SyncModes: []int64{dept.SyncModeRealTime, dept.SyncModeTiming}}, nil)
	if err != nil {
		return err
	}
	for _, po := range pos {
		_, err := departmentmanagelogic.NewDeptSyncJobExecuteLogic(d.ctx, d.svcCtx).DeptSyncJobExecute(&sys.DeptSyncJobExecuteReq{JobID: po.ID})
		if err != nil {
			d.Error(err)
		}
	}
	return nil
}
