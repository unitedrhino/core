package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppPolicyReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppPolicyReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppPolicyReadLogic {
	return &AppPolicyReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppPolicyReadLogic) AppPolicyRead(in *sys.AppPolicyReadReq) (*sys.AppPolicy, error) {
	po, err := relationDB.NewAppPolicyRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AppPolicyFilter{
		AppCode: in.AppCode,
		Code:    in.Code,
	})
	return utils.Copy[sys.AppPolicy](po), err
}
