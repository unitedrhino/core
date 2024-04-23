package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"

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

func (l *IndexLogic) Index(req *types.UserInfoIndexReq) (resp *types.UserInfoIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	info, err := l.svcCtx.UserRpc.UserInfoIndex(l.ctx, &sys.UserInfoIndexReq{
		Page:     logic.ToSysPageRpc(req.Page),
		UserName: req.UserName,
		Phone:    req.Phone,
		Email:    req.Email,
		Account:  req.Account,
	})
	if err != nil {
		return nil, err
	}

	var userInfo []*types.UserInfo
	var total int64
	var needCover bool
	total = info.Total
	uc := ctxs.GetUserCtx(l.ctx)
	if !uc.IsAdmin || uc.TenantCode != def.TenantCodeDefault || uc.IsAllData != true {
		needCover = true
	}
	userInfo = make([]*types.UserInfo, 0, len(userInfo))
	for _, i := range info.List {
		if needCover {
			i.Password = ""
			i.WechatUnionID = ""
			i.LastIP = ""
			i.RegIP = ""
			i.CreatedTime = 0
			i.IsAllData = 0
			i.Phone = Cover(i.Phone)
			i.Email = Cover(i.Email)
		}
		userInfo = append(userInfo, user.UserInfoToApi(i, nil, nil))
	}

	return &types.UserInfoIndexResp{userInfo, total}, nil
}
func Cover(in *wrapperspb.StringValue) *wrapperspb.StringValue {
	if in == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: "xxx"}
}
