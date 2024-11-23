package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaLogic) Captcha(req *types.UserCaptchaReq) (resp *types.UserCaptchaResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	ret, err := l.svcCtx.UserRpc.UserCaptcha(l.ctx, &sys.UserCaptchaReq{Ip: ctxs.GetUserCtx(l.ctx).IP,
		Account: req.Account, Type: req.Type, Use: req.Use, Code: req.Code, CodeID: req.CodeID})
	if err != nil {
		l.Errorf("%s UserCaptcha err=%+v", utils.FuncName(), err)
		return nil, err
	}
	switch req.Type {
	case def.CaptchaTypeImage:
		url := l.svcCtx.Captcha.GetB64(ret.Code)
		return &types.UserCaptchaResp{
			CodeID: ret.CodeID,
			Expire: ret.Expire,
			Url:    url,
		}, nil
	default:
		return &types.UserCaptchaResp{
			CodeID: ret.CodeID,
			Expire: ret.Expire,
		}, nil
	}
}
