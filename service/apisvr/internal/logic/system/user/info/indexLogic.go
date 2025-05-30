package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"

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
		Page:           logic.ToSysPageRpc(req.Page),
		UserName:       req.UserName,
		Phone:          req.Phone,
		Email:          req.Email,
		Account:        req.Account,
		UserIDs:        req.UserIDs,
		NickName:       req.NickName,
		HasAccessAreas: req.HasAccessAreas,
		RoleCode:       req.RoleCode,
		DeptID:         req.DeptID,
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
		var (
			roles []*sys.RoleInfo
			depts []*sys.DeptInfo
		)
		if req.WithRoles == true {
			ret, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
				UserID: i.UserID,
			})
			if err != nil {
				return nil, err
			}
			roles = ret.List
		}
		if req.WithDepts {
			ret, err := l.svcCtx.UserRpc.UserDeptIndex(l.ctx, &sys.UserDeptIndexReq{UserID: i.UserID})
			if err != nil {
				return nil, err
			}
			depts = ret.List
		}

		userInfo = append(userInfo, user.UserInfoToApi(i, user.UserOpt{Roles: roles, Depts: depts}))
	}

	return &types.UserInfoIndexResp{List: userInfo, PageResp: logic.ToPageResp(req.Page, total)}, nil
}
func Cover(in *wrapperspb.StringValue) *wrapperspb.StringValue {
	if in == nil {
		return nil
	}

	return &wrapperspb.StringValue{Value: maskPhoneNumberOrEmail(in.Value)}
}

// maskPhoneNumberOrEmail 接受一个手机号或邮箱作为参数，并在适当位置替换为星号
func maskPhoneNumberOrEmail(input string) string {
	if strings.Count(input, "@") == 0 && (len(input) == 10 || len(input) == 11) {
		// 假设输入是一个手机号
		return input[:3] + "****" + input[len(input)-4:]
	} else if strings.Count(input, "@") == 1 {
		// 假设输入是一个邮箱
		localPart := strings.Split(input, "@")[0]
		domainPart := strings.Split(input, "@")[1]
		return localPart[len(localPart)-4:] + "@" + "****" + domainPart[len(domainPart)-4:]
	}
	return input
}
