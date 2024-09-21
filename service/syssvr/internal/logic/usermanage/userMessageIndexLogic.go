package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"time"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	ll, err := l.GetUserLastReadGlobal()
	if err != nil {
		ll = time.Time{}
	}
	mis, err := relationDB.NewMessageInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.MessageInfoFilter{
		Group:      in.Group,
		NotifyCode: in.NotifyCode,
		IsGlobal:   def.True,
		NotifyTime: stores.CmpGte(ll),
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(mis) != 0 {
		var users []*relationDB.SysUserMessage
		for _, v := range mis {
			users = append(users, &relationDB.SysUserMessage{
				UserID:    ctxs.GetUserCtx(l.ctx).UserID,
				Group:     v.Group,
				MessageID: v.ID,
				IsRead:    def.False,
			})
		}
		err = relationDB.NewUserMessageRepo(l.ctx).MultiInsert(l.ctx, users)
		if err != nil {
			return nil, err
		}
		err = l.SetUserLastReadGlobal(time.Now())
		if err != nil {
			l.Error(err)
		}
	}
	f := relationDB.UserMessageFilter{
		WithMessage: true,
		Group:       in.Group,
		NotifyCode:  in.NotifyCode,
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
func (l *UserMessageIndexLogic) GetUserLastReadGlobal() (time.Time, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	tStr, err := caches.GetStore().Hget("cache:sys:userLastReadGlobal", cast.ToString(uc.UserID))
	if err != nil {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, tStr)
}
func (l *UserMessageIndexLogic) SetUserLastReadGlobal(t time.Time) error {
	uc := ctxs.GetUserCtx(l.ctx)
	return caches.GetStore().Hset("cache:sys:userLastReadGlobal", cast.ToString(uc.UserID), t.Format(time.RFC3339))
}
