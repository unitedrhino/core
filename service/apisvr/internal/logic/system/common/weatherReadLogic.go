package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type WeatherReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWeatherReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeatherReadLogic {
	return &WeatherReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var key = "b7de434f83c146e480d13ba6a565ce30"

type respType[t any] struct {
	Code string `json:"code"`
	Now  t      `json:"now"`
}

func (l *WeatherReadLogic) WeatherRead(req *types.WeatherReadReq) (resp *types.WeatherReadResp, err error) {
	ret, err := l.svcCtx.Common.WeatherRead(l.ctx, utils.Copy[sys.WeatherReadReq](req))
	return utils.Copy[types.WeatherReadResp](ret), err
}
