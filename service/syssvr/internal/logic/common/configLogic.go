package commonlogic

import (
	"context"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigLogic) Config(in *sys.Empty) (*sys.ConfigResp, error) {
	return &sys.ConfigResp{Map: &sys.Map{
		Mode:         l.svcCtx.Config.Map.Mode,
		AccessKey:    l.svcCtx.Config.Map.AccessKey,
		AccessSecret: l.svcCtx.Config.Map.AccessSecret,
	}}, nil
}
