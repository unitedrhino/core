package departmentmanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	usermanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobExecuteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobExecuteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobExecuteLogic {
	return &DeptSyncJobExecuteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var runMutex sync.Mutex
var runMap = make(map[int64]struct{})

func (l *DeptSyncJobExecuteLogic) DeptSyncJobExecute(in *sys.DeptSyncJobExecuteReq) (*sys.DeptSyncJobExecuteResp, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	po, err := relationDB.NewDeptSyncJobRepo(l.ctx).FindOne(l.ctx, in.JobID)
	if err != nil {
		return nil, err
	}
	canRun := func() bool {
		runMutex.Lock()
		defer runMutex.Unlock()
		if _, ok := runMap[po.ID]; ok {
			return false
		}
		runMap[po.ID] = struct{}{}
		return true
	}()
	if !canRun {
		return &sys.DeptSyncJobExecuteResp{}, nil
	}
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		defer func() {
			runMutex.Lock()
			defer runMutex.Unlock()
			delete(runMap, po.ID)
		}()
		start := time.Now()
		err := func() error {
			cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
				AppID:     po.ThirdConfig.AppID,
				AppKey:    po.ThirdConfig.AppKey,
				AppSecret: po.ThirdConfig.AppSecret,
			})
			if err != nil {
				return err
			}
			err = SyncDeptDing(ctx, cli, &relationDB.SysDeptInfo{
				ID:         def.RootNode,
				Name:       "根节点",
				Status:     def.True,
				DingTalkID: def.RootNode,
			})
			if err != nil {
				return err
			}
			var needSyncDeptMap = map[int64]*relationDB.SysDeptInfo{}
			var needSyncDepts = []*relationDB.SysDeptInfo{}
			var deptIDPaths = []string{}
			if len(po.SyncDeptIDs) == 0 { //指定只同步这几个部门的用户
				needSyncDepts, err = relationDB.NewDeptInfoRepo(ctx).FindByFilter(ctx, relationDB.DeptInfoFilter{}, &stores.PageInfo{Size: 2000})
				if err != nil {
					return err
				}
			} else {
				needSyncDepts, err = relationDB.NewDeptInfoRepo(ctx).FindByFilter(ctx, relationDB.DeptInfoFilter{IDs: po.SyncDeptIDs}, nil)
				if err != nil {
					return err
				}
				var idPaths []string
				for _, d := range needSyncDepts {
					idPaths = append(idPaths, d.IDPath)
				}
				depts, err := relationDB.NewDeptInfoRepo(ctx).FindByFilter(ctx, relationDB.DeptInfoFilter{IDPaths: idPaths}, nil)
				if err != nil {
					return err
				}
				needSyncDepts = append(needSyncDepts, depts...)
			}
			for _, d := range needSyncDepts {
				needSyncDeptMap[d.ID] = d
			}
			for _, d := range needSyncDeptMap {
				dd := d
				deptIDPaths = append(deptIDPaths, d.IDPath)
				err := SyncDeptUserDing(ctx, l.svcCtx, cli, dd)
				if err != nil {
					l.Error(dd, err)
				}
			}
			if len(deptIDPaths) > 0 {
				err = usermanagelogic.FillDeptUserCount(ctx, l.svcCtx, deptIDPaths...)
				if err != nil {
					l.Error(err)
				}
			}
			return nil
		}()
		logx.WithContext(ctx).Infof("DeptSyncJobExecute jobID:%v use:%v err:%v", po.ID, time.Now().Sub(start), err)
	})

	return &sys.DeptSyncJobExecuteResp{}, nil
}

func SyncDeptUserDing(ctx context.Context, svcCtx *svc.ServiceContext, cli *dingClient.DingTalk, info *relationDB.SysDeptInfo) error {
	c := 0
	page := 100
	var hasMore = true
	for hasMore {
		req := request.NewDeptDetailUserInfo(int(info.DingTalkID), c, page)
		limit.Wait(ctx)
		dings, err := cli.GetDeptDetailUserInfo(req.Build())
		if err != nil {
			return errors.System.AddDetail(err)
		}
		hasMore = dings.Page.HasMore
		c = dings.Page.NextCursor
		old, err := relationDB.NewDeptUserRepo(ctx).FindByFilter(ctx,
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
				po = deptPhoneMap[ding.Mobile]
			}
			if po == nil {
				po = deptEmailMap[ding.Email]
			}
			delete(deptIDMap, ding.UserId)
			delete(deptPhoneMap, ding.Mobile)
			delete(deptEmailMap, ding.Email)
			if po == nil {
				uc, err := relationDB.NewUserInfoRepo(ctx).FindOneByFilter(ctx, relationDB.UserInfoFilter{DingTalkUserID: ding.UserId, DingTalkUnionID: ding.UnionId})
				if err != nil {
					if !errors.Cmp(err, errors.NotFind) {
						return err
					}
					uc, err = relationDB.NewUserInfoRepo(ctx).FindOneByFilter(ctx, relationDB.UserInfoFilter{Accounts: []string{ding.Email, ding.Mobile}})
					if !errors.Cmp(err, errors.NotFind) {
						return err
					}
				}
				if uc == nil {
					userID := svcCtx.UserID.GetSnowflakeId()
					uc = &relationDB.SysUserInfo{
						UserID:         userID,
						DingTalkUserID: sql.NullString{Valid: true, String: ding.UserId},
						NickName:       ding.Name,
					}
					if ding.UnionId != "" {
						uc.DingTalkUnionID = sql.NullString{Valid: true, String: ding.UnionId}
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
					err = stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
						return usermanagelogic.Register(ctx, svcCtx, uc, tx)
					})
				}
				err = relationDB.NewDeptUserRepo(ctx).Insert(ctx, &relationDB.SysDeptUser{
					UserID:     uc.UserID,
					DeptID:     info.ID,
					DeptIDPath: info.IDPath,
				})
				if err != nil {
					return err
				}
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
				err = relationDB.NewUserInfoRepo(ctx).Update(ctx, po)
				if err != nil {
					return err
				}
			}
			if len(deptIDMap) > 0 { //如果存在删除的情况
				for _, one := range deptIDMap {
					err := relationDB.NewDeptUserRepo(ctx).DeleteByFilter(ctx, relationDB.DeptUserFilter{UserID: one.UserID})
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

var limit = rate.NewLimiter(38, 38)

func SyncDeptDing(ctx context.Context, cli *dingClient.DingTalk, info *relationDB.SysDeptInfo) error {
	req := request.NewDeptList()
	req.SetDeptId(int(info.DingTalkID))
	limit.Wait(ctx)
	dings, err := cli.GetDeptList(req.Build())
	if err != nil {
		return errors.System.AddDetail(err)
	}

	old, err := relationDB.NewDeptInfoRepo(ctx).FindByFilter(ctx, relationDB.DeptInfoFilter{ParentID: info.ID}, nil)
	if err != nil {
		return errors.System.AddDetail(err)
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
			err := relationDB.NewDeptInfoRepo(ctx).Insert(ctx, newOne)
			if err != nil {
				return err
			}
			newOne.IDPath = info.IDPath + fmt.Sprintf("%d-", newOne.ID)
			err = relationDB.NewDeptInfoRepo(ctx).Update(ctx, newOne)
			if err != nil {
				return err
			}
			continue
		}
		if po.Name != ding.Name || po.DingTalkID != int64(ding.Id) {
			delete(deptNameMap, po.Name)
			po.Name = ding.Name
			po.DingTalkID = int64(ding.Id)
			err = relationDB.NewDeptInfoRepo(ctx).Update(ctx, po)
			if err != nil {
				return err
			}

		}
	}
	if len(deptIDMap) > 0 { //如果存在删除的情况
		for _, one := range deptIDMap {
			err := relationDB.NewDeptInfoRepo(ctx).Delete(ctx, one.ID)
			if err != nil {
				return err
			}
		}
	}
	old, err = relationDB.NewDeptInfoRepo(ctx).FindByFilter(ctx, relationDB.DeptInfoFilter{ParentID: info.ID}, nil)
	if err != nil {
		return err
	}
	for _, o := range old {
		oo := o
		err := SyncDeptDing(ctx, cli, oo)
		if err != nil {
			logx.WithContext(ctx).Error(oo, err)
		}
	}
	return nil
}

var DingSync func(ctx context.Context, svcCtx *svc.ServiceContext, jobID int64, df *payload.DataFrame) error

func DeptSyncAddDing(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysDeptSyncJob) error {
	svcCtx.DingStreamMapMutex.Lock()
	defer svcCtx.DingStreamMapMutex.Unlock()
	cli1 := svcCtx.DingStreamMap[string(po.TenantCode)]
	if cli1 != nil {
		defer cli1.Close()
	}
	cli := dingClient.NewDingStream(po.ThirdConfig.AppKey, po.ThirdConfig.AppSecret)
	cli.RegisterAllEventRouter(func(c context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
		ctx, span := ctxs.StartSpan(c, "dingStreamEvent", "")
		defer span.End()
		startTime := timex.Now()
		duration := timex.Since(startTime)
		ctx = ctxs.WithRoot(ctx)
		ctx = ctxs.BindTenantCode(ctx, string(po.TenantCode), 0)
		err := DingSync(ctx, svcCtx, po.ID, df)
		logx.WithContext(ctx).WithDuration(duration).Infof(
			"subscribeDingStream df:%v err:%v",
			utils.Fmt(df), err)
		return event.NewSuccessResponse()
	})
	err := cli.Start(context.Background())
	if err != nil {
		logx.Error(err)
		return err
	}
	svcCtx.DingStreamMap[string(po.TenantCode)] = cli
	return nil
}

func DeptSyncDelDing(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysDeptSyncJob) error {
	svcCtx.DingStreamMapMutex.Lock()
	defer svcCtx.DingStreamMapMutex.Unlock()
	cli1 := svcCtx.DingStreamMap[string(po.TenantCode)]
	if cli1 != nil {
		cli1.Close()
		delete(svcCtx.DingStreamMap, string(po.TenantCode))
	}
	return nil
}
