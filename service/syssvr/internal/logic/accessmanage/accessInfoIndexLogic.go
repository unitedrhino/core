package accessmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoIndexLogic {
	return &AccessInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoIndexLogic) AccessInfoIndex(in *sys.AccessInfoIndexReq) (*sys.AccessInfoIndexResp, error) {
	f := relationDB.AccessFilter{
		Name:       in.Name,
		Code:       in.Code,
		Codes:      in.Codes,
		IsNeedAuth: in.IsNeedAuth,
		Group:      in.Group,
		WithApis:   in.WithApis,
	}
	total, err := relationDB.NewAccessRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewAccessRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var list []*sys.AccessInfo
	for _, v := range pos {
		list = append(list, ToAccessPb(v))
	}
	return &sys.AccessInfoIndexResp{List: list, Total: total}, nil
}
