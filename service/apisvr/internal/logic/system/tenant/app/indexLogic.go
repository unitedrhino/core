package app

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.TenantAppIndexReq) (resp *types.TenantAppIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	if len(ret.List) == 0 {
		return &types.TenantAppIndexResp{}, nil
	}
	appCodes := make([]string, 0)
	codeIDMap := make(map[string]*sys.TenantAppInfo)
	for _, v := range ret.List {
		appCodes = append(appCodes, v.AppCode)
		codeIDMap[v.AppCode] = v
	}
	apps, err := l.svcCtx.AppRpc.AppInfoIndex(l.ctx, &sys.AppInfoIndexReq{
		Codes: appCodes,
	})
	if err != nil {
		return nil, err
	}
	var retList []*types.TenantApp
	for _, v := range apps.List {
		ta := codeIDMap[v.Code]
		v.Id = ta.Id //修正为关联的id
		val := utils.Copy[types.TenantApp](v)
		if ta.WxMini != nil && ta.WxMini.AppID != "" {
			val.WxMini = utils.Copy[types.ThirdAppConfig](ta.WxMini)
		}
		if ta.DingMini != nil && ta.DingMini.AppID != "" {
			val.DingMini = utils.Copy[types.ThirdAppConfig](ta.DingMini)
		}
		if ta.WxOpen != nil && ta.WxOpen.AppID != "" {
			val.WxOpen = utils.Copy[types.ThirdAppConfig](ta.WxOpen)
		}
		if ta.Android != nil {
			val.Android = utils.Copy[types.ThirdApp](ta.Android)
		}
		val.LoginTypes = ta.LoginTypes
		retList = append(retList, val)
	}
	return &types.TenantAppIndexResp{
		List: retList,
	}, nil
}
