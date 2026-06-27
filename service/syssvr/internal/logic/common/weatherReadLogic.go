package commonlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/parnurzeal/gorequest"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
)

// singleflight: 防止缓存击穿时并发请求同时打到和风 API
var weatherSf syncx.SingleFlight

func init() {
	weatherSf = syncx.NewSingleFlight()
}

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

const (
	weatherGeoCacheTTL       = 30 * 24 * 60 * 60
	weatherDataCacheTTL      = 8 * 60 * 60
	weatherGeoCachePrecision = 1
)

// geoResp 和风 GeoAPI 城市搜索返回结构
type geoResp struct {
	Code     string `json:"code"`
	Location []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"location"`
}

func weatherGeoCacheKeyAndLookup(lon, lat float64) (string, string) {
	lonText := fmt.Sprintf("%.*f", weatherGeoCachePrecision, lon)
	latText := fmt.Sprintf("%.*f", weatherGeoCachePrecision, lat)
	return fmt.Sprintf("sys:common:geo:%s:%s", latText, lonText), fmt.Sprintf("%s,%s", lonText, latText)
}

func isWeatherPointValid(p *sys.Point) bool {
	if p == nil {
		return false
	}
	lon, lat := p.GetLongitude(), p.GetLatitude()
	if lon < -180 || lon > 180 || lat < -90 || lat > 90 {
		return false
	}
	return lon != 0 || lat != 0
}

func cloneWeatherPoint(p *sys.Point) *sys.Point {
	if p == nil {
		return nil
	}
	return &sys.Point{
		Longitude: p.GetLongitude(),
		Latitude:  p.GetLatitude(),
	}
}

func weatherProjectID(requestProjectID, userProjectID int64) (int64, bool) {
	if requestProjectID > def.NotClassified {
		return requestProjectID, true
	}
	if userProjectID > def.NotClassified {
		return userProjectID, false
	}
	return 0, false
}

func selectWeatherPosition(requestPosition *sys.Point, requestProjectID, userProjectID int64,
	project *sys.ProjectInfo, projectErr error) (*sys.Point, int64, error) {
	projectID, explicitProject := weatherProjectID(requestProjectID, userProjectID)
	if projectID != 0 {
		if projectErr == nil && project != nil && isWeatherPointValid(project.GetPosition()) {
			return cloneWeatherPoint(project.GetPosition()), projectID, nil
		}
		if explicitProject && projectErr != nil {
			return nil, projectID, projectErr
		}
		if isWeatherPointValid(requestPosition) {
			return cloneWeatherPoint(requestPosition), 0, nil
		}
		if projectErr != nil {
			return nil, projectID, projectErr
		}
	}
	if isWeatherPointValid(requestPosition) {
		return cloneWeatherPoint(requestPosition), 0, nil
	}
	return nil, 0, errors.Parameter.AddMsg("无法获取位置信息")
}

// resolveLocationID 经纬度 → 和风 LocationID（城市级），带 Redis 缓存（30 天）
func (l *WeatherReadLogic) resolveLocationID(lon, lat float64) (string, error) {
	geoKey, lookupLocation := weatherGeoCacheKeyAndLookup(lon, lat)
	cached, _ := caches.GetStore().GetCtx(l.ctx, geoKey)
	if cached != "" {
		return cached, nil
	}
	// 调和风 GeoAPI: /geo/v2/city/lookup
	var geo geoResp
	greq := gorequest.New().Retry(2, time.Second)
	resp, body, errs := greq.Get(fmt.Sprintf("https://%s/geo/v2/city/lookup?location=%s&key=%s&number=1",
		l.svcCtx.Config.Weather.ApiHost, lookupLocation, l.svcCtx.Config.Weather.ApiKey)).EndStruct(&geo)
	if errs != nil {
		return "", errors.System.AddDetail(string(body), resp, errs)
	}
	if resp.StatusCode != 200 || geo.Code != "200" || len(geo.Location) == 0 {
		return "", errors.System.AddMsgf("GeoAPI查询城市失败: code=%s, httpStatus=%d", geo.Code, resp.StatusCode)
	}
	locationID := geo.Location[0].ID
	// 城市映射极稳定，缓存 30 天
	caches.GetStore().SetexCtx(l.ctx, geoKey, locationID, weatherGeoCacheTTL)
	return locationID, nil
}

// fetchWeatherData 调和风 weather/now + air/now（air 失败降级），结果写入缓存（8h）
func (l *WeatherReadLogic) fetchWeatherData(locationID string) (*sys.WeatherReadResp, error) {
	var (
		weather respType[sys.WeatherReadResp]
		greq    = gorequest.New().Retry(3, time.Second*2)
	)
	// 实时天气（必须成功）
	resp, body, errs := greq.Get(fmt.Sprintf("https://%s/v7/weather/now?location=%s&key=%s",
		l.svcCtx.Config.Weather.ApiHost, locationID, l.svcCtx.Config.Weather.ApiKey)).EndStruct(&weather)
	if errs != nil {
		return nil, errors.System.AddDetail(string(body), resp, errs)
	}
	if resp.StatusCode != 200 {
		return nil, errors.System.AddDetail(string(body), resp, errs)
	}
	// 空气质量（降级：失败只记日志，不影响天气返回）
	var air respType[sys.WeatherAir]
	airResp, airBody, airErrs := greq.Get(fmt.Sprintf("https://%s/v7/air/now?location=%s&key=%s",
		l.svcCtx.Config.Weather.ApiHost, locationID, l.svcCtx.Config.Weather.ApiKey)).EndStruct(&air)
	if airErrs != nil || airResp.StatusCode != 200 {
		logx.Errorf("air/now 降级: locationID=%s, httpStatus=%d, errs=%v, body=%s",
			locationID, airResp.StatusCode, airErrs, string(airBody))
	} else {
		weather.Now.Air = &air.Now
	}
	// 以 LocationID 为 key 缓存 8 小时
	cacheKey := fmt.Sprintf("sys:common:weather:%s", locationID)
	caches.GetStore().SetexCtx(l.ctx, cacheKey, utils.MarshalNoErr(weather.Now), weatherDataCacheTTL)
	return &weather.Now, nil
}

func (l *WeatherReadLogic) WeatherRead(in *sys.WeatherReadReq) (*sys.WeatherReadResp, error) {
	if l.svcCtx.Config.Weather.ApiKey == "" || l.svcCtx.Config.Weather.ApiHost == "" {
		return &sys.WeatherReadResp{}, errors.Parameter.AddMsg("请联系管理员配置天气秘钥")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	var userProjectID int64
	if uc != nil {
		userProjectID = uc.ProjectID
	}
	requestProjectID := in.ProjectID
	projectID, _ := weatherProjectID(requestProjectID, userProjectID)
	var (
		project    *sys.ProjectInfo
		projectErr error
	)
	if projectID != 0 {
		project, projectErr = l.svcCtx.ProjectCache.GetData(l.ctx, projectID)
	}
	position, selectedProjectID, err := selectWeatherPosition(in.Position, requestProjectID, userProjectID, project, projectErr)
	if err != nil {
		return nil, err
	}
	if projectErr != nil && requestProjectID == 0 && in.Position != nil {
		logx.Errorf("天气项目定位读取失败，回退请求定位: projectID=%d, err=%v", projectID, projectErr)
	}
	in.ProjectID = selectedProjectID
	in.Position = position

	// 第一层缓存：经纬度 → 城市 LocationID（30 天）
	locationID, err := l.resolveLocationID(in.Position.Longitude, in.Position.Latitude)
	if err != nil {
		return nil, err
	}

	// 第二层缓存：LocationID → 天气数据（8 小时）
	cacheKey := fmt.Sprintf("sys:common:weather:%s", locationID)
	ret, _ := caches.GetStore().GetCtx(l.ctx, cacheKey)
	if ret != "" {
		var rett sys.WeatherReadResp
		if err := json.Unmarshal([]byte(ret), &rett); err == nil {
			return &rett, nil
		}
		logx.Errorf("天气缓存反序列化失败, key=%s, err=%v", cacheKey, err)
	}

	// singleflight: 同一 LocationID 并发请求只发一次和风调用
	val, err := weatherSf.Do(locationID, func() (any, error) {
		return l.fetchWeatherData(locationID)
	})
	if err != nil {
		return nil, err
	}
	return val.(*sys.WeatherReadResp), nil
}
