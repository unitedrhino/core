package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/users"
	"gitee.com/i-Things/share/utils"
	"github.com/golang-jwt/jwt/v5"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantOpenCheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantOpenCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantOpenCheckTokenLogic {
	return &TenantOpenCheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantOpenCheckTokenLogic) TenantOpenCheckToken(in *sys.TenantOpenCheckTokenReq) (*sys.TenantOpenCheckTokenResp, error) {
	var claim users.OpenClaims
	err := users.ParseTokenWithFunc(&claim, in.Token, func(token *jwt.Token) (interface{}, error) {
		if claim.TenantCode == "" || claim.UserID == 0 || claim.Code == "" {
			return nil, errors.TokenInvalid
		}
		po, err := relationDB.NewDataOpenAccessRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DataOpenAccessFilter{
			TenantCode: claim.TenantCode,
			UserID:     claim.UserID,
			Code:       claim.Code,
		})
		if err != nil {
			return nil, err
		}
		return []byte(po.AccessSecret), nil
	})
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	return &sys.TenantOpenCheckTokenResp{UserID: claim.UserID, IsAdmin: def.True, UserName: "open", TenantCode: claim.TenantCode}, nil
}
