package commonlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/parnurzeal/gorequest"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type WeatherReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWeatherReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeatherReadLogic {
	return &WeatherReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type respType[t any] struct {
	Code string `json:"code"`
	Now  t      `json:"now"`
}

//var key = "b7de434f83c146e480d13ba6a565ce30"

func (l *WeatherReadLogic) WeatherRead(in *sys.WeatherReadReq) (*sys.WeatherReadResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.Position == nil && in.ProjectID == 0 {
		in.ProjectID = uc.ProjectID
	}
	if in.ProjectID != 0 {
		pi, err := l.svcCtx.ProjectCache.GetData(l.ctx, in.ProjectID)
		if err != nil {
			return nil, err
		}
		in.Position = &sys.Point{
			Longitude: pi.Position.GetLongitude(),
			Latitude:  pi.Position.GetLatitude(),
		}
	}
	cacheKey := fmt.Sprintf("sys:common:weather:%.2f:%.2f", in.Position.Latitude, in.Position.Longitude)
	ret, err := caches.GetStore().GetCtx(l.ctx, cacheKey)
	if ret != "" {
		var rett sys.WeatherReadResp
		json.Unmarshal([]byte(ret), &rett)
		return &rett, nil
	}
	tc, err := l.svcCtx.TenantConfigCache.GetData(l.ctx, uc.TenantCode)
	if err != nil {
		return nil, err
	}
	key := tc.WeatherKey
	var (
		weather respType[sys.WeatherReadResp]
		air     respType[sys.WeatherAir]
		greq    = gorequest.New().Retry(3, time.Second*2)
	)
	//参考: https://dev.qweather.com/
	resp, body, errs := greq.Get(fmt.Sprintf("https://devapi.qweather.com/v7/weather/now?location=%v,%v&key=%s",
		in.Position.Longitude, in.Position.Latitude, key)).EndStruct(&weather)
	if errs != nil {
		return nil, errors.System.AddDetail(string(body), resp, errs)
	}
	resp, body, errs = greq.Get(fmt.Sprintf("https://devapi.qweather.com/v7/air/now?location=%v,%v&key=%s",
		in.Position.Longitude, in.Position.Latitude, key)).EndStruct(&air)
	if errs != nil {
		return nil, errors.System.AddDetail(string(body), resp, errs)
	}
	weather.Now.Air = &air.Now
	caches.GetStore().SetexCtx(l.ctx, cacheKey, utils.MarshalNoErr(weather.Now), 60*60*1) //1个小时的有效期
	return &weather.Now, nil
}
