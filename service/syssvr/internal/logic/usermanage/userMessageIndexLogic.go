package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

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
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	return &sys.UserMessageIndexResp{Total: total, List: utils.CopySlice[sys.UserMessage](pos)}, nil
}
