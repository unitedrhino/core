package slot

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/share/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cast"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const (
	AuthTypeCore = "core"
)

const (
	CodeAreaInfo      = "areaInfo"      //区域
	CodeUserSubscribe = "userSubscribe" //用户订阅
	CodeDataFilter    = "dataFilter"    //大数据过滤
	CodeDeviceSend    = "deviceSend"    //设备控制

	SubCodeCreate              = "create" //创建
	SubCodeDelete              = "delete" //创建
	SubCodePropertyControlSend = "propertyControlSend"
)

type Info struct {
	Code     string            `json:"code"` // 鉴权的编码
	SubCode  string            `json:"subCode"`
	SlotCode string            `json:"slotCode"` //slot的编码
	Method   string            `json:"method"`   // 请求方式 GET  POST
	Uri      string            `json:"uri"`      // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Hosts    []string          `json:"hosts"`    //访问的地址 host or host:port
	Body     string            `json:"body"`     // body 参数模板
	Handler  map[string]string `json:"handler"`  //http头
	AuthType string            `json:"authType"` //鉴权类型 core
}

type Infos []*Info

func (i Infos) Request(ctx context.Context, in any, retV any) error {
	uc := ctxs.GetUserCtx(ctx)
	for _, v := range i {
		greq := gorequest.New().Retry(3, time.Second*2)
		t, err := template.New(v.SlotCode + ":uri").Parse(v.Uri)
		if err != nil {
			return err
		}
		var out strings.Builder
		err = t.Execute(&out, in)
		if err != nil {
			return err
		}
		url := v.Hosts[0] + out.String()
		switch v.Method {
		case http.MethodGet:
			greq.Get(url)
		case http.MethodPost:
			var str string
			if v.Body == "" {
				str = utils.MarshalNoErr(in)
			} else {
				t, err = template.New(v.SlotCode + ":body").Parse(v.Body)
				if err != nil {
					return err
				}
				var out strings.Builder
				err = t.Execute(&out, in)
				if err != nil {
					return err
				}
				str = out.String()
			}
			greq.Post(url).Type("json").Send(str)
		}
		switch v.AuthType {
		case AuthTypeCore:
			greq.Set(ctxs.UserTokenKey, uc.Token).Set(ctxs.UserProjectID, cast.ToString(uc.ProjectID)).
				Set(ctxs.UserAppCodeKey, uc.AppCode)
		}
		var ret result.ResponseSuccessBean
		if retV != nil {
			ret.Data = retV
		}
		_, _, errs := greq.EndStruct(&ret)
		if errs != nil {
			return errors.System.AddDetail(errs)
		}
		if ret.Code != errors.OK.Code {
			return errors.NewCodeError(ret.Code, ret.Msg)
		}
	}
	return nil

}
