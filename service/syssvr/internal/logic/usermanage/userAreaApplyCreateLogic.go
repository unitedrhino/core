package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

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
	ai, err := relationDB.NewAreaInfoRepo(l.ctx).FindOne(ctxs.WithAdmin(l.ctx), in.AreaID, nil)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("区域不存在")
		}
		return nil, err
	}
	err = relationDB.NewUserAreaApplyRepo(l.ctx).Insert(ctxs.WithAdmin(l.ctx), &relationDB.SysUserAreaApply{
		ProjectID: ai.ProjectID,
		UserID:    ctxs.GetUserCtx(l.ctx).UserID,
		AreaID:    dataType.AreaID(in.AreaID),
		AuthType:  in.AuthType,
	})
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return &sys.Empty{}, nil
		}

		return nil, err
	}
	return &sys.Empty{}, nil
}
