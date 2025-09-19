package commonlogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

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
	return &sys.ConfigResp{Map: utils.Copy[sys.Map](l.svcCtx.Config.Map), IsEnableTaRegister: l.svcCtx.Config.TenantConfig.IsEnableTaRegister}, nil
}
