package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoReadLogic {
	return &ModuleInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoReadLogic) ModuleInfoRead(in *sys.WithIDCode) (*sys.ModuleInfo, error) {
	ret, err := relationDB.NewModuleInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ModuleInfoFilter{Codes: []string{in.Code}, ID: in.Id})
	return logic.ToModuleInfoPb(ret), err
}
