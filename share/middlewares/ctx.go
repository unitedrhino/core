package middlewares

import (
	"bytes"
	"encoding/json"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
	"io"
	"net/http"
	"strings"
)

const bodySize = 256
const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "ur_requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "code", "tenantCode"},
		Buckets:   []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 750, 1000, 2000, 5000, 10000, 20000, 50000, 100000},
	})
)

func InitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = ctxs.InitCtxWithReq(r)
		reqBody, _ := io.ReadAll(r.Body)                //读取 reqBody
		r.Body = io.NopCloser(bytes.NewReader(reqBody)) //重建 reqBody
		if len(reqBody) > bodySize {
			reqBody = reqBody[0:bodySize]
		}
		reqBody = bytes.ReplaceAll(reqBody, []byte("\n"), []byte{})
		reqBody = bytes.ReplaceAll(reqBody, []byte("\r"), []byte{})
		uc := ctxs.GetUserCtxNoNil(r.Context())
		r = ctxs.NeedResp(r)
		startTime := timex.Now()
		defer func() {
			resp := ctxs.GetResp(r)
			var respBody []byte
			if resp == nil {
				resp = &http.Response{}
			} else if resp.Body != nil {
				respBody, _ = io.ReadAll(resp.Body) //读取 respBody
				if len(respBody) > bodySize {
					respBody = respBody[0:bodySize]
				}
			}
			useTime := timex.Since(startTime)
			metricServerReqDur.Observe(useTime.Milliseconds(),
				r.URL.Path, cast.ToString(resp.StatusCode), uc.TenantCode)
			if strings.Contains(r.Header.Get(def.ContentType), def.ApplicationJson) {
				logx.WithContext(r.Context()).Infof("[HTTP %v %v] %s use:%v uc:[%v]  reqBody:[%v] respBody:[%v]",
					resp.StatusCode, resp.Status, r.RequestURI, useTime, utils.Fmt(uc), string(reqBody), string(respBody))
			} else {
				logx.WithContext(r.Context()).Infof("[HTTP %v %v] %s use:%v uc:[%v]  respBody:[%v]",
					resp.StatusCode, resp.Status, r.RequestURI, useTime, utils.Fmt(uc), string(respBody))
			}

		}()
		defer utils.Recoverf(r.Context(), "uri:%s uc:%v  req:%v",
			r.RequestURI, utils.Fmt(uc), string(reqBody))
		defer func() {
			if p := recover(); p != nil {
				result.Http(w, r, nil, errors.Panic)
				utils.HandleThrow(r.Context(), "uri:%s uc:%v  req:%v err:%v",
					r.RequestURI, utils.Fmt(uc), string(reqBody), cast.ToString(p))
				ret := ctxs.GetResp(r)
				err := errors.Panic.AddDetail(p)
				if ret != nil {
					//将接口的应答结果写入r.Response，为操作日志记录接口提供应答信息
					var temp http.Response
					temp.StatusCode = int(err.GetCode())
					temp.Status = err.GetMsg()
					if ret != nil {
						bs, _ := json.Marshal(ret)
						temp.Body = io.NopCloser(bytes.NewReader(bs))
					}
					*ret = temp
				}
			}
		}()
		next(w, r)

	}
}
