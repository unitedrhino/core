package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMessageIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMessageIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMessageIndexLogic {
	return &UserMessageIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserMessageIndexLogic) UserMessageIndex(in *sys.UserMessageIndexReq) (*sys.UserMessageIndexResp, error) {
	db := relationDB.NewUserMessageRepo(l.ctx)
	err := UpdateMsg(l.ctx, in.NotifyCode, in.Group)
	if err != nil {
		return nil, err
	}
	f := relationDB.UserMessageFilter{
		WithMessage: true,
		Group:       in.Group,
		NotifyCode:  in.NotifyCode,
		CreatedTime: utils.Copy[def.TimeRange](in.CreatedTime),
		IsRead:      in.IsRead,
		Str1:        in.Str1,
		Str2:        in.Str2,
		Str3:        in.Str3,
	}
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "isRead",
		Sort:  stores.OrderDesc,
	}, stores.OrderBy{
		Field: "createdTime",
		Sort:  stores.OrderDesc,
	}))
	return &sys.UserMessageIndexResp{Total: total, List: utils.CopySlice[sys.UserMessage](pos)}, nil
}
func getUserLastReadGlobal(ctx context.Context) (time.Time, error) {
	uc := ctxs.GetUserCtx(ctx)
	tStr, err := caches.GetStore().Hget("cache:sys:userLastReadGlobal", cast.ToString(uc.UserID))
	if err != nil {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, tStr)
}
func setUserLastReadGlobal(ctx context.Context, t time.Time) error {
	uc := ctxs.GetUserCtx(ctx)
	return caches.GetStore().Hset("cache:sys:userLastReadGlobal", cast.ToString(uc.UserID), t.Format(time.RFC3339))
}

func UpdateMsg(ctx context.Context, NotifyCode string, Group string) error {
	ll, err := getUserLastReadGlobal(ctx)
	if err != nil {
		ll = time.Time{}
	}
	mis, err := relationDB.NewMessageInfoRepo(ctx).FindByFilter(ctx, relationDB.MessageInfoFilter{
		Group:      Group,
		NotifyCode: NotifyCode,
		IsGlobal:   def.True,
		NotifyTime: stores.CmpGte(ll),
	}, nil)
	if err != nil {
		return err
	}
	if len(mis) != 0 {
		var users []*relationDB.SysUserMessage
		for _, v := range mis {
			users = append(users, &relationDB.SysUserMessage{
				UserID:    ctxs.GetUserCtx(ctx).UserID,
				Group:     v.Group,
				MessageID: v.ID,
				IsRead:    def.False,
			})
		}
		err = relationDB.NewUserMessageRepo(ctx).MultiInsert(ctx, users)
		if err != nil {
			return err
		}
		err = setUserLastReadGlobal(ctx, time.Now())
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	return nil
}
