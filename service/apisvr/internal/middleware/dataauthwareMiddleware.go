package middleware

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/config"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/def"
	"net/http"
)

type DataAuthWareMiddleware struct {
	cfg config.Config
}

func NewDataAuthWareMiddleware(cfg config.Config) *DataAuthWareMiddleware {
	caches.InitStore(cfg.CacheRedis)
	return &DataAuthWareMiddleware{cfg: cfg}
}

type DataAuthParam struct {
	ProjectID  string   `json:"projectID,string,optional"` //项目id
	ProjectIDs []string `json:"projectIDs,optional"`       //项目ids
	AreaID     string   `json:"areaID,string,optional"`    //项目区域id
	AreaIDs    []string `json:"areaIDs,optional"`          //项目区域ids
}

func (m *DataAuthWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		return
	}
}

func (m *DataAuthWareMiddleware) check(ctx context.Context, dataType def.AuthDataType, reqIDs []string) int {
	//if len(reqIDs) > 0 {
	//	authIDs, err := caches.GetUserDataAuth(ctx, dataType)
	//	diffIDs := utils.SliceLeftDiff(reqIDs, authIDs)
	//	if err == redis.Nil || (err == nil && len(diffIDs) > 0) { //没有数据权限
	//		logx.WithContext(ctx).Errorf("%s.没有数据权限 dataType=%#v, reqIDs=%#v, authIDs=%#v, diffIDs=%#v", utils.FuncName(), dataType, reqIDs, authIDs, diffIDs)
	//		return http.StatusUnauthorized
	//	} else if err != nil { //校验数据权限异常
	//		logx.WithContext(ctx).Errorf("%s.校验数据权限异常 dataType=%#v, reqIDs=%#v, authIDs=%#v, error=%#v", utils.FuncName(), dataType, reqIDs, authIDs, err)
	//		return http.StatusInternalServerError
	//	}
	//}
	return http.StatusOK
}
