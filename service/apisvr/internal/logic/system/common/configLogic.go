package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config() (resp *types.ConfigResp, err error) {
	rsp, err := l.svcCtx.Common.Config(l.ctx, &sys.Empty{})
	if err != nil {
		err = errors.Fmt(err)
		l.Errorf("%s.rpc.SysConfig err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &types.ConfigResp{Map: types.Map{Mode: rsp.Map.Mode, AccessKey: rsp.Map.AccessKey, AccessSecret: rsp.Map.AccessSecret},
		Oss: types.Oss{Host: l.svcCtx.Config.OssConf.CustomHost}}, nil
}
