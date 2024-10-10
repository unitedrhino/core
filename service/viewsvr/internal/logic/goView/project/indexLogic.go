package project

import (
	"context"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/logic"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/viewsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/types"

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

func (l *IndexLogic) Index(req *types.ProjectInfoIndexReq) (resp *types.ProjectInfoIndexResp, err error) {
	size, err := relationDB.NewProjectInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.ProjectInfoFilter{})
	if err != nil {
		return nil, err
	}
	pi, err := relationDB.NewProjectInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProjectInfoFilter{}, logic.ToPageInfo(req.Page))
	if err != nil {
		return nil, err
	}
	var (
		list []*types.ProjectInfo
	)
	for _, v := range pi {
		list = append(list, ToProjectInfoTypes(v))
	}
	return &types.ProjectInfoIndexResp{
		List:  list,
		Total: size,
		Num:   int64(len(list)),
	}, nil
}
