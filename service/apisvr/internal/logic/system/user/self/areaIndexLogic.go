package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaIndexLogic {
	return &AreaIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaIndexLogic) AreaIndex() (resp *types.AreaInfoIndexResp, err error) {
	//var (
	//	areas []*sys.AreaInfo
	//)
	//uc := ctxs.GetUserCtx(l.ctx)
	//ret, err := l.svcCtx.UserRpc.UserAreaIndex(l.ctx, &sys.UserAreaIndexReq{
	//	UserID: uc.UserID,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//if len(ret.List) != 0 {
	//	var areaIDs []int64
	//	for _, v := range ret.List {
	//		areaIDs = append(areaIDs, v.AreaID)
	//	}
	//	ret2, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{AreaIDs: areaIDs})
	//	if err != nil {
	//		return nil, err
	//	}
	//	areas = ret2.List
	//}
	//return &types.AreaInfoIndexResp{
	//	List: info.ToAreaInfosTypes(areas),
	//}, nil
	return nil, err
}
