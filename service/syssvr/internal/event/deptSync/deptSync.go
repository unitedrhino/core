package deptSync

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	departmentmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/departmentmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"github.com/zeromicro/go-zero/core/logx"
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
}

func (d DeptSync) HandleDing(jobID int64, df *payload.DataFrame) error {
	fmt.Println(df)
	eventType := df.GetHeader("eventType")
	switch eventType {
	case "user_modify_org": //用户信息变更
		var data userModifyOrg
		err := json.Unmarshal([]byte(df.Data), &data)
		if err != nil {
			return err
		}
		for _, v := range data.DiffInfo {
			ui, err := relationDB.NewUserInfoRepo(d.ctx).FindByFilter(d.ctx, relationDB.UserInfoFilter{DingTalkUserID: v.Userid}, nil)
			if err != nil {
				return err
			}
			fmt.Println(ui)
		}

	}
	po, err := relationDB.NewDeptSyncJobRepo(d.ctx).FindOne(d.ctx, jobID)
	if err != nil {
		return err
	}
	if po.SyncMode != dept.SyncModeRealTime {
		return nil
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
