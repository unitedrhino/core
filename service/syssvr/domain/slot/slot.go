package slot

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
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

type Info struct {
	Code     string            `json:"code"`     // 鉴权的编码
	SlotCode string            `json:"slotCode"` //slot的编码
	Method   string            `json:"method"`   // 请求方式 GET  POST
	Uri      string            `json:"uri"`      // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Proto    string            `json:"proto"`    //http https
	Hosts    []string          `json:"hosts"`    //访问的地址 host or host:port
	Body     string            `json:"body"`     // body 参数模板
	Handler  map[string]string `json:"handler"`  //http头
	AuthType string            `json:"authType"` //鉴权类型 core
}

type Infos []*Info

func (i Infos) Request(ctx context.Context, in any) error {
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
			t, err = template.New(v.SlotCode + ":body").Parse(v.Body)
			if err != nil {
				return err
			}
			var out strings.Builder
			err = t.Execute(&out, in)
			if err != nil {
				return err
			}
			greq.Post(url).Type(gorequest.TypeJSON).Send(out.String())
		}
		switch v.AuthType {
		case AuthTypeCore:
			greq.Set(ctxs.UserTokenKey, uc.Token).Set(ctxs.UserProjectID, cast.ToString(uc.ProjectID)).
				Set(ctxs.UserAppCodeKey, uc.AppCode).Set(ctxs.UserRoleKey, cast.ToString(uc.RoleID))
		}
		var ret result.ResponseSuccessBean
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
