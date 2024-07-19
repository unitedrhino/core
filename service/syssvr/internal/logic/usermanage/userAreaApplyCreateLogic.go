package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyCreateLogic {
	return &UserAreaApplyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyCreateLogic) UserAreaApplyCreate(in *sys.UserAreaApplyCreateReq) (*sys.Empty, error) {
	_, err := relationDB.NewAreaInfoRepo(l.ctx).FindOne(ctxs.WithAllArea(l.ctx), in.AreaID, nil)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("区域不存在")
		}
		return nil, err
	}
	err = relationDB.NewUserAreaApplyRepo(l.ctx).Insert(ctxs.WithRoot(l.ctx), &relationDB.SysUserAreaApply{
		UserID:   ctxs.GetUserCtx(l.ctx).UserID,
		AreaID:   stores.AreaID(in.AreaID),
		AuthType: in.AuthType,
	})
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return &sys.Empty{}, nil
		}

		return nil, err
	}
	return &sys.Empty{}, nil
}
